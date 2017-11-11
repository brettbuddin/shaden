// Package wav provides WAV file decoding
package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
)

const (
	formatPCM       = 1
	formatIEEEFloat = 3

	preambleFormat = "fmt "
	preambleData   = "data"
)

// Wav is a WAV file
type Wav struct {
	Header
	Samples int
	Reader  io.ReadCloser
}

// ReadAll reads all samples from the WAV file
func (w *Wav) ReadAll() ([]float32, error) {
	return w.Read(w.Samples)
}

// Read reads a specific number of samples of the WAV file
func (w *Wav) Read(n int) ([]float32, error) {
	var data interface{}
	switch w.AudioFormat {
	case formatIEEEFloat:
		data = make([]float32, n)
	case formatPCM:
		switch w.BitsPerSample {
		case 8:
			data = make([]uint8, n)
		case 16:
			data = make([]int16, n)
		default:
			return nil, fmt.Errorf("invalid bits per sample: %v", w.BitsPerSample)
		}
	default:
		return nil, fmt.Errorf("unknown sample type")
	}

	if err := binary.Read(w.Reader, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	var final []float32
	switch d := data.(type) {
	case []uint8:
		final = make([]float32, len(d))
		for i, v := range d {
			final[i] = float32(v) / math.MaxUint8
		}
	case []uint16:
		final = make([]float32, len(d))
		for i, v := range d {
			final[i] = (float32(v) - math.MinInt16) / (math.MaxInt16 - math.MinInt16)
		}
	case []float32:
		final = d
	}

	return final, nil
}

// Close closes the WAV file
func (w *Wav) Close() error {
	return w.Reader.Close()
}

// Header is a WAV file header
type Header struct {
	AudioFormat    uint16
	NumChannels    uint16
	SampleRate     uint32
	BytesPerSecond uint32
	BytesPerBlock  uint16
	BitsPerSample  uint16
}

// Open opens and reads the header of a WAV file
func Open(path string) (*Wav, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	wav, err := load(f)
	if err != nil {
		return nil, err
	}
	return wav, nil
}

func load(r io.ReadCloser) (*Wav, error) {
	var wav Wav

	header := make([]byte, 16)
	if _, err := io.ReadFull(r, header[:12]); err != nil {
		return nil, err
	}
	if string(header[0:4]) != "RIFF" {
		return nil, fmt.Errorf("no RIFF preamble")
	}
	if string(header[8:12]) != "WAVE" {
		return nil, fmt.Errorf("no WAVE preamble")
	}

	var format bool
	for {
		if _, err := io.ReadFull(r, header[:8]); err != nil {
			return nil, err
		}
		size := binary.LittleEndian.Uint32(header[4:])

		switch string(header[:4]) {
		case preambleFormat:
			if err := readFormat(r, &wav, size); err != nil {
				return nil, err
			}
			format = true
		case preambleData:
			if !format {
				return nil, fmt.Errorf("premature data block; no format preamble")
			}
			establishReader(r, &wav, size)
			return &wav, nil
		default:
			if _, err := io.CopyN(ioutil.Discard, r, int64(size)); err != nil {
				return nil, err
			}
		}
	}
}

func readFormat(r io.Reader, wav *Wav, size uint32) error {
	if size < 16 {
		return fmt.Errorf("invalid format size")
	}
	b := make([]byte, size)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	if err := binary.Read(buf, binary.LittleEndian, &wav.Header); err != nil {
		return err
	}
	if wav.AudioFormat != formatPCM && wav.AudioFormat != formatIEEEFloat {
		return fmt.Errorf("unknown format: %v", wav.AudioFormat)
	}
	return nil
}

func establishReader(r io.ReadCloser, wav *Wav, size uint32) {
	wav.Samples = int(size) / int(wav.BitsPerSample) * 8
	wav.Reader = &LimitReadCloser{io.LimitReader(r, int64(size)), r}
}

// LimitReadCloser is a LimitReader that can close
type LimitReadCloser struct {
	io.Reader
	c io.Closer
}

// Close closes the LimitReadCloser
func (lrc *LimitReadCloser) Close() error {
	return lrc.c.Close()
}
