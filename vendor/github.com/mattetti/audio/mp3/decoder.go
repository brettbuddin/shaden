package mp3

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/mattetti/audio/mp3/id3v1"
	"github.com/mattetti/audio/mp3/id3v2"
)

// Decoder operates on a reader and extracts important information
// See http://www.mp3-converter.com/mp3codec/mp3_anatomy.htm
type Decoder struct {
	r         io.Reader
	NbrFrames int

	ID3v2tag *id3v2.Tag
}

// NewDecoder creates a new reader reading the given reader and parsing its data.
// It is the caller's responsibility to call Close on the reader when done.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// SeemsValid checks if the mp3 file looks like a valid mp3 file by looking at the first few bytes.
// The data can be corrupt but at least the header seems alright.
// It is the caller's responsibility to rewind/close the reader when done.
func SeemsValid(r io.Reader) bool {
	d := New(r)
	fr := &Frame{}
	var frameDuration time.Duration
	var duration time.Duration
	var err error
	var badFrames int
	for {
		err = d.Next(fr)
		if err != nil {
			badFrames++
			if err == ErrInvalidHeader {
				continue
			}
			break
		}
		// garbage needing to be skipped probably means bad frame
		if fr.SkippedBytes > 20 {
			badFrames++
		}
		frameDuration = fr.Duration()

		if frameDuration > 0 {
			duration += frameDuration
		}
		d.NbrFrames++
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF || err == io.ErrShortBuffer {
		err = nil
	}
	if d.NbrFrames <= 0 {
		return false
	}
	percentBadFrames := (float64(badFrames) * 100) / float64(d.NbrFrames)
	// more than 10% frames with issues or a zero/negative duration means bad file
	if percentBadFrames > 10 {
		return false
	}
	if duration <= 0 {
		return false
	}
	return true
}

// Duration returns the time duration for the current mp3 file
// The entire reader will be consumed, the consumer might want to rewind the reader
// if they want to read more from the feed.
// Note that this is an estimated duration based on how the frames look. An invalid file might have
// a duration.
func (d *Decoder) Duration() (time.Duration, error) {
	if d == nil {
		return 0, errors.New("can't calculate the duration of a nil pointer")
	}
	fr := &Frame{}
	var frameDuration time.Duration
	var duration time.Duration
	var err error
	for {
		err = d.Next(fr)
		if err != nil {
			// bad headers can be ignored and hopefully skipped
			if err == ErrInvalidHeader {
				continue
			}
			break
		}
		frameDuration = fr.Duration()
		if frameDuration > 0 {
			duration += frameDuration
		}
		d.NbrFrames++
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF || err == io.ErrShortBuffer {
		err = nil
	}

	return duration, err
}

// Next decodes the next frame into the provided frame structure.
func (d *Decoder) Next(f *Frame) error {
	if f == nil {
		return fmt.Errorf("can't decode to a nil Frame")
	}

	var n int
	f.SkippedBytes = 0
	f.Counter++

	hLen := 4
	if f.buf == nil {
		f.buf = make([]byte, hLen)
	} else {
		f.buf = f.buf[:hLen]
	}

	_, err := io.ReadAtLeast(d.r, f.buf, hLen)
	if err != nil {
		return err
	}

	// ID3v1 tag at the beggining
	if bytes.Compare(f.buf[:3], id3v1.HeaderTagID) == 0 {
		// the ID3v1 tag is always 128 bytes long, we already read 4 bytes
		// so we need to read the rest.
		buf := make([]byte, 124)
		// TODO: parse the actual header
		if _, err := io.ReadAtLeast(d.r, buf, 124); err != nil {
			return ErrInvalidHeader
		}
		buf = append(f.buf, buf...)
		// that wasn't a frame
		f = &Frame{}
		return nil
	}

	// ID3v2 tag
	if bytes.Compare(f.buf[:3], id3v2.HeaderTagID) == 0 {
		d.ID3v2tag = &id3v2.Tag{}
		// we already read 4 bytes, an id3v2 tag header is of size 10, read the rest
		// and append it to what we already have.
		buf := make([]byte, 6)
		n, err := d.r.Read(buf)
		if err != nil || n != 6 {
			return ErrInvalidHeader
		}
		buf = append(f.buf, buf...)

		th := id3v2.TagHeader{}
		copy(th[:], buf)
		if err = d.ID3v2tag.ReadHeader(th); err != nil {
			return err
		}
		// TODO: parse the actual tag
		// Skip the tag for now
		bytesToSkip := int64(d.ID3v2tag.Header.Size)
		var cn int64
		if cn, err = io.CopyN(ioutil.Discard, d.r, bytesToSkip); cn != bytesToSkip {
			return ErrInvalidHeader
		}
		f = &Frame{}
		return err
	}

	f.Header = FrameHeader(f.buf)
	if !f.Header.IsValid() {
		f.Header, n, err = d.skipToNextFrame()
		if err != nil {
			return err
		}
		f.SkippedBytes = n
	}

	dataSize := f.Header.Size()
	if dataSize > 4 {
		// substract the 4 bytes we already read
		dataSize -= 4
		f.buf = append(f.buf, make([]byte, dataSize)...)
		_, err = io.ReadAtLeast(d.r, f.buf[4:], int(dataSize))
	}
	return err
}

// skipToSyncWord reads until it finds a frame header
func (d *Decoder) skipToNextFrame() (fh FrameHeader, readN int, err error) {
	if d == nil {
		return nil, readN, errors.New("nil decoder")
	}
	buf := make([]byte, 1)
	lookAheadBuf := make([]byte, 1)
	var n int
	for {
		n, err = d.r.Read(buf)
		readN += n
		if err != nil {
			return nil, readN, err
		}
		readN++
		if buf[0] == 0xFF {
			if _, err := d.r.Read(lookAheadBuf); err != nil {
				return nil, readN, err
			}
			readN++
			if lookAheadBuf[0]&0xE0 == 0xE0 {
				buf = []byte{0xff, lookAheadBuf[0], 0, 0}
				n, err := d.r.Read(buf[2:])
				if err != nil {
					return nil, readN + n, err
				}
				if n != 2 {
					return nil, readN + n, io.ErrUnexpectedEOF
				}
				readN += 2
			}
			return buf, readN, err
		}
	}
}
