package caf

import (
	"errors"
	"io"
)

var (
	fileHeaderID = [4]byte{'c', 'a', 'f', 'f'}

	// Chunk IDs
	StreamDescriptionChunkID = [4]byte{'d', 'e', 's', 'c'}
	AudioDataChunkID         = [4]byte{'d', 'a', 't', 'a'}
	ChannelLayoutChunkID     = [4]byte{'c', 'h', 'a', 'n'}
	FillerChunkID            = [4]byte{'f', 'r', 'e', 'e'}
	MarkerChunkID            = [4]byte{'m', 'a', 'r', 'k'}
	RegionChunkID            = [4]byte{'r', 'e', 'g', 'n'}
	InstrumentChunkID        = [4]byte{'i', 'n', 's', 't'}
	MagicCookieID            = [4]byte{'k', 'u', 'k', 'i'}
	InfoStringsChunkID       = [4]byte{'i', 'n', 'f', 'o'}
	EditCommentsChunkID      = [4]byte{'e', 'd', 'c', 't'}
	PacketTableChunkID       = [4]byte{'p', 'a', 'k', 't'}
	StringsChunkID           = [4]byte{'s', 't', 'r', 'g'}
	UUIDChunkID              = [4]byte{'u', 'u', 'i', 'd'}
	PeakChunkID              = [4]byte{'p', 'e', 'a', 'k'}
	OverviewChunkID          = [4]byte{'o', 'v', 'v', 'w'}
	MIDIChunkID              = [4]byte{'m', 'i', 'd', 'i'}
	UMIDChunkID              = [4]byte{'u', 'm', 'i', 'd'}
	FormatListID             = [4]byte{'l', 'd', 's', 'c'}
	IXMLChunkID              = [4]byte{'i', 'X', 'M', 'L'}

	// Format IDs
	// Linear PCM
	AudioFormatLinearPCM = [4]byte{'l', 'p', 'c', 'm'}
	// Apple’s implementation of IMA 4:1 ADPCM. Has no format flags.
	AudioFormatAppleIMA4 = [4]byte{'i', 'm', 'a', '4'}
	// MPEG-4 AAC. The mFormatFlags field must contain the MPEG-4 audio object type constant indicating the specific kind of data.
	AudioFormatMPEG4AAC = [4]byte{'a', 'a', 'c', ' '}
	// MACE 3:1; has no format flags.
	AudioFormatMACE3 = [4]byte{'M', 'A', 'C', '3'}
	// MACE 6:1; has no format flags.
	AudioFormatMACE6 = [4]byte{'M', 'A', 'C', '6'}
	// μLaw 2:1; has no format flags.
	AudioFormatULaw = [4]byte{'u', 'l', 'a', 'w'}
	// aLaw 2:1; has no format flags.
	AudioFormatALaw = [4]byte{'a', 'l', 'a', 'w'}
	// MPEG-1 or 2, Layer 1 audio. Has no format flags.
	AudioFormatMPEGLayer1 = [4]byte{'.', 'm', 'p', '1'}
	// MPEG-1 or 2, Layer 2 audio. Has no format flags.
	AudioFormatMPEGLayer2 = [4]byte{'.', 'm', 'p', '2'}
	// MPEG-1 or 2, Layer 3 audio (that is, MP3). Has no format flags.
	AudioFormatMPEGLayer3 = [4]byte{'.', 'm', 'p', '3'}
	// Apple Lossless; has no format flags.
	AudioFormatAppleLossless = [4]byte{'a', 'l', 'a', 'c'}

	// ErrFmtNotSupported is a generic error reporting an unknown format.
	ErrFmtNotSupported = errors.New("format not supported")
	// ErrUnexpectedData is a generic error reporting that the parser encountered unexpected data.
	ErrUnexpectedData = errors.New("unexpected data content")
)

func New(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

//func NewParser(r io.Reader, ch chan *TBD) *Decoder {
//return &Decoder{r: r, Ch: ch}
//}
