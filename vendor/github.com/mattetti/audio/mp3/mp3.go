// mp3 is a package used to access mp3 information
// See: http://sea-mist.se/fou/cuppsats.nsf/all/857e49b9bfa2d753c125722700157b97/$file/Thesis%20report-%20MP3%20Decoder.pdf
// Uses some code from https://github.com/tcolgate/mp3 under MIT license Tristan Colgate-McFarlane and badgerodon
package mp3

import (
	"errors"
	"io"
)

//go:generate stringer -type=FrameVersion
const (
	MPEG25 FrameVersion = iota
	MPEGReserved
	MPEG2
	MPEG1
)

//go:generate stringer -type=FrameLayer
const (
	LayerReserved FrameLayer = iota
	Layer3
	Layer2
	Layer1
)

//go:generate stringer -type=FrameEmphasis
const (
	EmphNone FrameEmphasis = iota
	Emph5015
	EmphReserved
	EmphCCITJ17
)

//go:generate stringer -type=FrameChannelMode
const (
	Stereo FrameChannelMode = iota
	JointStereo
	DualChannel
	SingleChannel
)

var (
	// ID31HBytes are the 2 bytes starting the ID3 v1 tag
	ID31HBytes = []byte{0xFF, 0xFB}
	// XingTAGID Xing vbr tag
	XingTAGID = []byte{0x58, 0x69, 0x6E, 0x67}
	// InfoTAGID Info cbr tag
	InfoTAGID = []byte{0x49, 0x6E, 0x66, 0x6F}
)

var (
	// ErrNoSyncBits implies we could not find a valid frame header sync bit before EOF
	ErrNoSyncBits = errors.New("EOF before sync bits found")

	// ErrPrematureEOF indicates that the filed ended before a complete frame could be read
	ErrPrematureEOF = errors.New("EOF mid stream")

	ErrInvalidHeader = errors.New("invalid header")

	// ErrInvalidBitrate indicates that the header information did not contain a recognized bitrate
	ErrInvalidBitrate FrameBitRate = -1

	// ErrInvalidSampleRate indicates that no samplerate could be found for the frame header provided
	ErrInvalidSampleRate = FrameSampleRate(-1)

	bitrates = map[FrameVersion]map[FrameLayer][15]int{
		MPEG1: { // MPEG 1
			Layer1: {0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448}, // Layer1
			Layer2: {0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384},    // Layer2
			Layer3: {0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320},     // Layer3
		},
		MPEG2: { // MPEG 2, 2.5
			Layer1: {0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256}, // Layer1
			Layer2: {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},      // Layer2
			Layer3: {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},      // Layer3
		},
	}

	sampleRates = map[FrameVersion][3]int{
		MPEG1:        {44100, 48000, 32000},
		MPEG2:        {22050, 24000, 16000},
		MPEG25:       {11025, 12000, 8000},
		MPEGReserved: {0, 0, 0},
	}

	slotSize = map[FrameLayer]int{
		LayerReserved: 0,
		Layer3:        1,
		Layer2:        1,
		Layer1:        4,
	}

	samplesPerFrame = map[FrameVersion]map[FrameLayer]int{
		MPEG1: {
			Layer1: 384,
			Layer2: 1152,
			Layer3: 1152,
		},
		MPEG2: {
			Layer1: 384,
			Layer2: 1152,
			Layer3: 576,
		},
	}
)

func New(r io.Reader) *Decoder {
	return &Decoder{r: r}
}
