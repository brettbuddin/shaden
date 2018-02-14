package aiff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/mattetti/audio"
)

// Decoder is the wrapper structure for the AIFF container
type Decoder struct {
	r io.ReadSeeker

	// ID is always 'FORM'. This indicates that this is a FORM chunk
	ID [4]byte
	// Size contains the size of data portion of the 'FORM' chunk.
	// Note that the data portion has been
	// broken into two parts, formType and chunks
	Size uint32
	// Form describes what's in the 'FORM' chunk. For Audio IFF files,
	// formType (aka Format) is always 'AIFF'.
	// This indicates that the chunks within the FORM pertain to sampled sound.
	Form [4]byte

	// Data coming from the COMM chunk
	commSize        uint32
	NumChans        uint16
	NumSampleFrames uint32
	BitDepth        uint16
	SampleRate      int
	//
	PCMSize  uint32
	PCMChunk *Chunk

	// AIFC data
	Encoding     [4]byte
	EncodingName string

	err             error
	pcmDataAccessed bool

	Debug bool
}

// NewDecoder creates a new reader reading the given reader and pushing audio data to the given channel.
// It is the caller's responsibility to call Close on the reader when done.
func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{r: r}
}

// SampleBitDepth returns the bit depth encoding of each sample.
func (d *Decoder) SampleBitDepth() int32 {
	if d == nil {
		return 0
	}
	return int32(d.BitDepth)
}

// PCMLen returns the total number of bytes in the PCM data chunk
func (d *Decoder) PCMLen() int64 {
	if d == nil {
		return 0
	}
	return int64(d.PCMSize)
}

// Err returns the first non-EOF error that was encountered by the Decoder.
func (d *Decoder) Err() error {
	if d.err == io.EOF {
		return nil
	}
	return d.err
}

// EOF returns positively if the underlying reader reached the end of file.
func (d *Decoder) EOF() bool {
	if d == nil || d.err == io.EOF {
		return true
	}
	return false
}

// WasPCMAccessed returns positively if the PCM data was previously accessed.
func (d *Decoder) WasPCMAccessed() bool {
	if d == nil {
		return false
	}
	return d.pcmDataAccessed
}

// NextChunk returns the next available chunk
func (d *Decoder) NextChunk() (*Chunk, error) {
	if d.err = d.readHeaders(); d.err != nil {
		d.err = fmt.Errorf("failed to read header - %v", d.err)
		return nil, d.err
	}

	var (
		id   [4]byte
		size uint32
	)

	id, size, d.err = d.iDnSize()
	if d.err != nil {
		d.err = fmt.Errorf("error reading chunk header - %v", d.err)
		return nil, d.err
	}

	c := &Chunk{
		ID:   id,
		Size: int(size),
		R:    io.LimitReader(d.r, int64(size)),
	}
	return c, d.err
}

// IsValidFile verifies that the file is valid/readable.
func (d *Decoder) IsValidFile() bool {
	d.ReadInfo()
	if d.err != nil && d.err != io.EOF {
		return false
	}
	if d.NumChans < 1 {
		return false
	}
	if d.BitDepth < 8 {
		return false
	}
	if d, err := d.Duration(); err != nil || d <= 0 {
		return false
	}

	return true
}

// Duration returns the time duration for the current AIFF container
func (d *Decoder) Duration() (time.Duration, error) {
	if d == nil {
		return 0, errors.New("can't calculate the duration of a nil pointer")
	}
	d.ReadInfo()
	if err := d.Err(); err != nil {
		return 0, err
	}
	duration := time.Duration(float64(d.NumSampleFrames) / float64(d.SampleRate) * float64(time.Second))
	return duration, nil
}

// FwdToPCM forwards the underlying reader until the start of the PCM chunk.
// If the PCM chunk was already read, no data will be found (you need to rewind).
func (d *Decoder) FwdToPCM() error {
	if d.err = d.readHeaders(); d.err != nil {
		d.err = fmt.Errorf("failed to read header - %v", d.err)
		return nil
	}

	// read the file information to setup the audio clip
	// find the beginning of the SSND chunk and set the clip reader to it.
	var rewindBytes int64

	var chunk *Chunk
	for d.err == nil {
		chunk, d.err = d.NextChunk()
		if d.err != nil {
			return d.err
		}
		switch chunk.ID {
		case COMMID:
			if err := d.parseCommChunk(uint32(chunk.Size)); err != nil {
				return err
			}
			// if we found the sound data before the COMM,
			// we need to rewind the reader so we can properly
			// set the clip reader.
			if rewindBytes > 0 {
				d.r.Seek(-rewindBytes, 1)
			}
		case SSNDID:
			//            SSND chunk: Must be defined
			//   0      4 bytes  "SSND"
			//   4      4 bytes  <Chunk size(x)>
			//   8      4 bytes  <Offset(n)>
			//  12      4 bytes  <block size>
			//  16     (n)bytes  Comment
			//  16+(n) (s)bytes  <Sample data>

			var offset uint32
			if d.err = chunk.ReadBE(&offset); d.err != nil {
				d.err = fmt.Errorf("PCM offset failed to parse - %s", d.err)
				return d.err
			}

			if d.err = chunk.ReadBE(&d.PCMSize); d.err != nil {
				d.err = fmt.Errorf("PCMSize failed to parse - %s", d.err)
				return d.err
			}
			if offset > 0 {
				d.PCMSize -= offset
				// skip pcm comment
				buf := make([]byte, offset)
				if err := chunk.ReadBE(&buf); err != nil {
					return err
				}
			}
			d.PCMChunk = chunk
			d.pcmDataAccessed = true
			return nil

		default:
			// if we read SSN but didn't read the COMM, we need to track location
			if d.SampleRate == 0 {
				rewindBytes += int64(chunk.Size)
			}
			chunk.Done()
		}
	}
	return nil
}

// Reset resets the decoder (and rewind the underlying reader)
func (d *Decoder) Reset() {
	d.ID = [4]byte{}
	d.Size = 0
	d.Form = [4]byte{}
	d.commSize = 0
	d.NumChans = 0
	d.NumSampleFrames = 0
	d.BitDepth = 0
	d.SampleRate = 0
	d.Encoding = [4]byte{}
	d.EncodingName = ""
	d.err = nil
	d.pcmDataAccessed = false
	d.r.Seek(0, 0)
}

// Format returns the audio format of the decoded content.
func (d *Decoder) Format() *audio.Format {
	if d == nil {
		return nil
	}
	return &audio.Format{
		NumChannels: int(d.NumChans),
		SampleRate:  int(d.SampleRate),
		BitDepth:    int(d.BitDepth),
		Endianness:  binary.LittleEndian,
	}
}

// FullPCMBuffer is an inneficient way to access all the PCM data contained in the
// audio container. The entire PCM data is held in memory.
// Consider using Buffer() instead.
func (d *Decoder) FullPCMBuffer() (*audio.PCMBuffer, error) {
	if !d.WasPCMAccessed() {
		err := d.FwdToPCM()
		if err != nil {
			return nil, d.err
		}
	}
	format := d.Format()

	buf := audio.NewPCMIntBuffer(make([]int, 4096), format)
	decodeF, err := sampleDecodeFunc(int(d.BitDepth))
	if err != nil {
		return nil, fmt.Errorf("could not get sample decode func %v", err)
	}

	i := 0
	for err == nil {
		buf.Ints[i], err = decodeF(d.PCMChunk)
		if err != nil {
			break
		}
		i++
		// grow the underlying slice if needed
		if i == len(buf.Ints) {
			buf.Ints = append(buf.Ints, make([]int, 4096)...)
		}
	}
	buf.Ints = buf.Ints[:i]

	if err == io.EOF {
		err = nil
	}

	return buf, err
}

// PCMBuffer populates the passed PCM buffer
func (d *Decoder) PCMBuffer(buf *audio.PCMBuffer) error {
	if buf == nil {
		return nil
	}

	if !d.pcmDataAccessed {
		if d.Debug {
			fmt.Println("forwarding to PCM Data")
		}
		err := d.FwdToPCM()
		if err != nil {
			return d.err
		}
	}

	// TODO: avoid a potentially unecessary allocation
	format := &audio.Format{
		NumChannels: int(d.NumChans),
		SampleRate:  int(d.SampleRate),
		BitDepth:    int(d.BitDepth),
		Endianness:  binary.LittleEndian,
	}

	decodeF, err := sampleDecodeFunc(int(d.BitDepth))
	if err != nil {
		return fmt.Errorf("could not get sample decode func %v", err)
	}

	// Note that we populate the buffer even if the
	// size of the buffer doesn't fit an even number of frames.
	if d.Debug {
		fmt.Println("populating %d samples", len(buf.Ints))
	}
	for i := 0; i < len(buf.Ints); i++ {
		buf.Ints[i], err = decodeF(d.r)
		if err != nil {
			break
		}
	}
	if err == io.EOF {
		err = nil
	}
	buf.Format = format
	if buf.DataType != audio.Integer {
		buf.DataType = audio.Integer
	}

	return err
}

// String implements the Stringer interface.
func (d *Decoder) String() string {
	out := fmt.Sprintf("Format: %s - ", d.Form)
	if d.Form == aifcID {
		out += fmt.Sprintf("%s - ", d.EncodingName)
	}
	if d.SampleRate != 0 {
		out += fmt.Sprintf("%d channels @ %d / %d bits - ", d.NumChans, d.SampleRate, d.BitDepth)
		dur, _ := d.Duration()
		out += fmt.Sprintf("Duration: %f seconds\n", dur.Seconds())
	}
	return out
}

// iDnSize returns the next ID + block size
func (d *Decoder) iDnSize() ([4]byte, uint32, error) {
	var ID [4]byte
	var blockSize uint32
	if d.err = binary.Read(d.r, binary.BigEndian, &ID); d.err != nil {
		return ID, blockSize, d.err
	}
	if d.err = binary.Read(d.r, binary.BigEndian, &blockSize); d.err != nil {
		return ID, blockSize, d.err
	}
	return ID, blockSize, nil
}

// readHeaders is safe to call multiple times
// byte size of the header: 12
func (d *Decoder) readHeaders() error {
	// prevent the headers to be re-read
	if d.Size > 0 {
		return nil
	}
	if d.err = binary.Read(d.r, binary.BigEndian, &d.ID); d.err != nil {
		return d.err
	}
	// Must start by a FORM header/ID
	if d.ID != formID {
		d.err = fmt.Errorf("%s - %s", ErrFmtNotSupported, d.ID)
		return d.err
	}

	if d.err = binary.Read(d.r, binary.BigEndian, &d.Size); d.err != nil {
		return d.err
	}
	if d.err = binary.Read(d.r, binary.BigEndian, &d.Form); d.err != nil {
		return d.err
	}

	// Must be a AIFF or AIFC form type
	if d.Form != aiffID && d.Form != aifcID {
		d.err = fmt.Errorf("%s - %s", ErrFmtNotSupported, d.Form)
		return d.err
	}

	return nil
}

// ReadInfo reads the underlying reader until the comm header is parsed.
// This method is safe to call multiple times.
func (d *Decoder) ReadInfo() {
	if d == nil || d.SampleRate > 0 {
		return
	}
	if d.err = d.readHeaders(); d.err != nil {
		d.err = fmt.Errorf("failed to read header - %v", d.err)
		return
	}

	var (
		id          [4]byte
		size        uint32
		rewindBytes int64
	)
	for d.err != io.EOF {
		id, size, d.err = d.iDnSize()
		if d.err != nil {
			d.err = fmt.Errorf("error reading chunk header - %v", d.err)
			break
		}
		switch id {
		case COMMID:
			d.parseCommChunk(size)
			// if we found other chunks before the COMM,
			// we need to rewind the reader so we can properly
			// read the rest later.
			if rewindBytes > 0 {
				d.r.Seek(-(rewindBytes + int64(size)), 1)
				break
			}
			return
		default:
			// we haven't read the COMM chunk yet, we need to track location to rewind
			if d.SampleRate == 0 {
				rewindBytes += int64(size)
			}
			if d.err = d.jumpTo(int(size)); d.err != nil {
				return
			}
		}
	}
}

func (d *Decoder) parseCommChunk(size uint32) error {
	d.commSize = size
	// don't re-parse the comm chunk
	if d.NumChans > 0 {
		return nil
	}

	if d.err = binary.Read(d.r, binary.BigEndian, &d.NumChans); d.err != nil {
		d.err = fmt.Errorf("num of channels failed to parse - %s", d.err)
		return d.err
	}
	if d.err = binary.Read(d.r, binary.BigEndian, &d.NumSampleFrames); d.err != nil {
		d.err = fmt.Errorf("num of sample frames failed to parse - %s", d.err)
		return d.err
	}
	if d.err = binary.Read(d.r, binary.BigEndian, &d.BitDepth); d.err != nil {
		d.err = fmt.Errorf("sample size failed to parse - %s", d.err)
		return d.err
	}
	var srBytes [10]byte
	if d.err = binary.Read(d.r, binary.BigEndian, &srBytes); d.err != nil {
		d.err = fmt.Errorf("sample rate failed to parse - %s", d.err)
		return d.err
	}
	d.SampleRate = audio.IeeeFloatToInt(srBytes)

	if d.Form == aifcID {
		if d.err = binary.Read(d.r, binary.BigEndian, &d.Encoding); d.err != nil {
			d.err = fmt.Errorf("AIFC encoding failed to parse - %s", d.err)
			return d.err
		}
		// pascal style string with the description of the encoding
		var size uint8
		if d.err = binary.Read(d.r, binary.BigEndian, &size); d.err != nil {
			d.err = fmt.Errorf("AIFC encoding failed to parse - %s", d.err)
			return d.err
		}

		desc := make([]byte, size)
		if d.err = binary.Read(d.r, binary.BigEndian, &desc); d.err != nil {
			d.err = fmt.Errorf("AIFC encoding failed to parse - %s", d.err)
			return d.err
		}
		d.EncodingName = string(desc)
	}

	return nil
}

// jumpTo advances the reader to the amount of bytes provided
func (d *Decoder) jumpTo(bytesAhead int) error {
	var err error
	if bytesAhead > 0 {
		_, err = io.CopyN(ioutil.Discard, d.r, int64(bytesAhead))
	}
	return err
}

func sampleDecodeFunc(bitDepth int) (func(io.Reader) (int, error), error) {
	switch bitDepth {
	case 8:
		// 8bit values are unsigned
		return func(r io.Reader) (int, error) {
			var v uint8
			err := binary.Read(r, binary.BigEndian, &v)
			return int(v), err
		}, nil
	case 16:
		return func(r io.Reader) (int, error) {
			var v int16
			err := binary.Read(r, binary.BigEndian, &v)
			return int(v), err
		}, nil
	case 24:
		return func(r io.Reader) (int, error) {
			// TODO: check if the conversion might not be inversed depending on
			// the encoding (BE vs LE)
			var output int32
			d := make([]byte, 3)
			_, err := r.Read(d)
			if err != nil {
				return 0, err
			}
			output |= int32(d[2]) << 0
			output |= int32(d[1]) << 8
			output |= int32(d[0]) << 16
			return int(output), nil
		}, nil
	case 32:
		return func(r io.Reader) (int, error) {
			var v int32
			err := binary.Read(r, binary.BigEndian, &v)
			return int(v), err
		}, nil
	default:
		return nil, fmt.Errorf("%v bit depth not supported", bitDepth)
	}
}

func sampleFloat64DecodeFunc(bitDepth int) (func(io.Reader) (float64, error), error) {
	switch bitDepth {
	case 8:
		// 8bit values are unsigned
		return func(r io.Reader) (float64, error) {
			var v uint8
			err := binary.Read(r, binary.BigEndian, &v)
			return float64(v), err
		}, nil
	case 16:
		return func(r io.Reader) (float64, error) {
			var v int16
			err := binary.Read(r, binary.BigEndian, &v)
			return float64(v), err
		}, nil
	case 24:
		return func(r io.Reader) (float64, error) {
			// TODO: check if the conversion might not be inversed depending on
			// the encoding (BE vs LE)
			var output int32
			d := make([]byte, 3)
			_, err := r.Read(d)
			if err != nil {
				return 0, err
			}
			output |= int32(d[2]) << 0
			output |= int32(d[1]) << 8
			output |= int32(d[0]) << 16
			return float64(output), nil
		}, nil
	case 32:
		return func(r io.Reader) (float64, error) {
			var v float32
			err := binary.Read(r, binary.BigEndian, &v)
			return float64(v), err
		}, nil
	default:
		return nil, fmt.Errorf("%v bit depth not supported", bitDepth)
	}
}
