package aiff

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/mattetti/audio"
)

// Encoder encodes LPCM data into an aiff content.
type Encoder struct {
	w          io.WriteSeeker
	SampleRate int
	BitDepth   int
	NumChans   int

	WrittenBytes    int
	frames          int
	pcmChunkStarted bool
	pcmChunkSizePos int
}

// NewEncoder creates a new encoder to create a new aiff file.
// Don't forget to add Frames to the encoder before writing.
func NewEncoder(w io.WriteSeeker, sampleRate, bitDepth, numChans int) *Encoder {
	return &Encoder{
		w:          w,
		SampleRate: sampleRate,
		BitDepth:   bitDepth,
		NumChans:   numChans,
	}
}

// AddBE serializes and adds the passed value using big endian
func (e *Encoder) AddBE(src interface{}) error {
	e.WrittenBytes += binary.Size(src)
	return binary.Write(e.w, binary.BigEndian, src)
}

// AddLE serializes and adds the passed value using little endian
func (e *Encoder) AddLE(src interface{}) error {
	e.WrittenBytes += binary.Size(src)
	return binary.Write(e.w, binary.LittleEndian, src)
}

func (e *Encoder) addBuffer(buf *audio.PCMBuffer) error {
	if buf == nil {
		return fmt.Errorf("can't add a nil buffer")
	}

	frameCount := buf.Size()
	buf.CacheInts()
	for i := 0; i < frameCount; i++ {
		// TODO(mattetti): support float encoded wav files
		for j := 0; j < buf.Format.NumChannels; j++ {
			v := buf.Ints[i*buf.Format.NumChannels+j]
			switch e.BitDepth {
			case 8:
				if err := e.AddBE(uint8(v)); err != nil {
					return err
				}
			case 16:
				if err := e.AddBE(uint16(v)); err != nil {
					return err
				}
			case 24:
				if err := e.AddBE(audio.Uint32toUint24Bytes(uint32(v))); err != nil {
					return err
				}
			case 32:
				if err := e.AddBE(uint32(v)); err != nil {
					return err
				}
			default:
				return fmt.Errorf("can't add frames of bit size %d", e.BitDepth)
			}
		}
		e.frames++
	}

	return nil
}

func (e *Encoder) writeHeader() error {
	if e == nil {
		return fmt.Errorf("can't write a nil encoder")
	}
	if e.w == nil {
		return fmt.Errorf("can't write to a nil writer")
	}

	if e.WrittenBytes > 0 {
		return nil
	}

	// ID
	if err := e.AddBE(formID); err != nil {
		return fmt.Errorf("%v when writing FORM header", err)
	}
	// size, will need to be updated later on (total size - 8)
	if err := e.AddBE(uint32(0)); err != nil {
		return fmt.Errorf("%v when writing size header", err)
	}
	// Format
	if err := e.AddBE(aiffID); err != nil {
		return fmt.Errorf("%v when writing format header", err)
	}
	// comm chunk
	if err := e.AddBE(COMMID); err != nil {
		return fmt.Errorf("%v when writing comm chunk ID header", err)
	}
	// blocksize uint32
	if err := e.AddBE(uint32(18)); err != nil {
		return fmt.Errorf("%v when writing comm chunk size header", err)
	}
	if err := e.AddBE(uint16(e.NumChans)); err != nil {
		return fmt.Errorf("%v when writing comm chan numbers", err)
	}
	// number of sample frames (unknown at this point)
	// will have to come back and edit
	if err := e.AddBE(uint32(42)); err != nil {
		return fmt.Errorf("%v when writing comm num sample frames", err)
	}
	if err := e.AddBE(uint16(e.BitDepth)); err != nil {
		return fmt.Errorf("%v when writing comm chan numbers", err)
	}
	// sample rate in IeeeFloat (10 bytes)
	if err := e.AddBE(audio.IntToIeeeFloat(int(e.SampleRate))); err != nil {
		return fmt.Errorf("%v when writing comm sample rate", err)
	}
	return nil
}

func (e *Encoder) Write(buf *audio.PCMBuffer) error {
	if err := e.writeHeader(); err != nil {
		return err
	}

	if !e.pcmChunkStarted {
		// other chunks
		// audio frames
		if err := e.AddBE([]byte("SSND")); err != nil {
			return fmt.Errorf("%v when writing SSND chunk ID header", err)
		}

		e.pcmChunkSizePos = e.WrittenBytes
		e.pcmChunkStarted = true

		// temporary blocksize uint32
		//chunksize := uint32((int(e.BitDepth)/8)*int(e.NumChans)*len(e.Frames) + 8)
		if err := e.AddBE(uint32(84)); err != nil {
			return fmt.Errorf("%v when writing SSND chunk size header", err)
		}

		if err := e.AddBE(uint32(0)); err != nil {
			return fmt.Errorf("%v when writing SSND offset", err)
		}
		if err := e.AddBE(uint32(0)); err != nil {
			return fmt.Errorf("%v when writing SSND block size", err)
		}
	}

	return e.addBuffer(buf)
}

// Close flushes the content to disk, make sure the headers are up to date
// Note that the underlying writter is NOT being closed.
func (e *Encoder) Close() error {
	// go back and write total size
	if _, err := e.w.Seek(4, 0); err != nil {
		return err
	}
	if err := e.AddBE(uint32(e.WrittenBytes) - 8); err != nil {
		return fmt.Errorf("%v when writing the total written bytes", err)
	}
	if _, err := e.w.Seek(22, 0); err != nil {
		return err
	}
	if err := e.AddBE(uint32(e.frames)); err != nil {
		return fmt.Errorf("%v when writing the total of frames", err)
	}
	// rewrite the audio chunk length header
	if e.pcmChunkSizePos > 0 {
		if _, err := e.w.Seek(int64(e.pcmChunkSizePos), 0); err != nil {
			return err
		}
		chunksize := uint32((int(e.BitDepth)/8)*int(e.NumChans)*e.frames + 8)
		if err := e.AddBE(uint32(chunksize)); err != nil {
			return fmt.Errorf("%v when writing wav data chunk size header", err)
		}
	}
	// jump to the end of the file.
	e.w.Seek(0, 2)
	switch e.w.(type) {
	case *os.File:
		e.w.(*os.File).Sync()
	}
	return nil
}
