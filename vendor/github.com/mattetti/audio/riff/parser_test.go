package riff

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseHeader(t *testing.T) {
	expectations := []struct {
		input  string
		id     [4]byte
		size   uint32
		format [4]byte
	}{
		{"fixtures/sample.rmi", RiffID, 29632, rmiFormatID},
		{"fixtures/sample.wav", RiffID, 53994, WavFormatID},
		{"fixtures/sample.avi", RiffID, 230256, aviFormatID},
	}

	for _, exp := range expectations {
		path, _ := filepath.Abs(exp.input)
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		c := New(f)
		err = c.ParseHeaders()
		if err != nil {
			t.Fatal(err)
		}
		if c.ID != exp.id {
			t.Fatalf("%s of %s didn't match %s, got %s", "ID", exp.input, exp.id, c.ID)
		}
		if c.Size != exp.size {
			t.Fatalf("%s of %s didn't match %d, got %d", "BlockSize", exp.input, exp.size, c.Size)
		}
		if c.Format != exp.format {
			t.Fatalf("%s of %s didn't match %q, got %q", "Format", exp.input, exp.format, c.Format)
		}
	}
}

func TestParseWavHeaders(t *testing.T) {
	expectations := []struct {
		input         string
		headerSize    uint32
		format        uint16
		numChans      uint16
		sampleRate    uint32
		byteRate      uint32
		blockAlign    uint16
		bitsPerSample uint16
	}{
		// mono audio files
		{"fixtures/sample.wav", 16, 1, 1, 44100, 88200, 2, 16},
		// sterep audio files with junk, bext and more headers
		{"fixtures/junkKick.wav", 40, 1, 2, 44100, 176400, 4, 16},
	}

	for _, exp := range expectations {
		path, _ := filepath.Abs(exp.input)
		t.Log(path)

		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		c := New(f)
		if err := c.ParseHeaders(); err != nil {
			t.Fatalf("%s for %s when parsing headers", err, path)
		}
		ch, err := c.NextChunk()
		if err != nil {
			t.Fatal(err)
		}

		for ; ch != nil; ch, err = c.NextChunk() {
			if err != nil {
				if err != io.EOF {
					t.Fatal(err)
				}
				break
			}

			if bytes.Compare(ch.ID[:], FmtID[:]) == 0 {
				ch.DecodeWavHeader(c)
			} else {
				ch.Done()
			}
		}

		if c.wavHeaderSize != exp.headerSize {
			t.Fatalf("%s didn't match %d, got %d", "header size", exp.headerSize, c.wavHeaderSize)
		}
		if c.WavAudioFormat != exp.format {
			t.Fatalf("%s didn't match %d, got %d", "audio format", exp.format, c.WavAudioFormat)
		}
		if c.NumChannels != exp.numChans {
			t.Fatalf("%s didn't match %d, got %d", "# of channels", exp.numChans, c.NumChannels)
		}

		if c.SampleRate != exp.sampleRate {
			t.Fatalf("%s didn't match %d, got %d", "SampleRate", exp.sampleRate, c.SampleRate)
		}
		if c.AvgBytesPerSec != exp.byteRate {
			t.Fatalf("%s didn't match %d, got %d", "ByteRate", exp.byteRate, c.AvgBytesPerSec)
		}
		if c.BlockAlign != exp.blockAlign {
			t.Fatalf("%s didn't match %d, got %d", "BlockAlign", exp.blockAlign, c.BlockAlign)
		}
		if c.BitsPerSample != exp.bitsPerSample {
			t.Fatalf("%s didn't match %d, got %d", "BitsPerSample", exp.bitsPerSample, c.BitsPerSample)
		}
	}

}

func TestContainerDuration(t *testing.T) {
	expectations := []struct {
		input string
		dur   time.Duration
	}{
		{"fixtures/sample.wav", time.Duration(612176870)},
	}

	for _, exp := range expectations {
		path, _ := filepath.Abs(exp.input)
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		c := New(f)
		d, err := c.Duration()
		if err != nil {
			t.Fatal(err)
		}
		if d != exp.dur {
			t.Fatalf("Container duration of %s didn't match %f, got %f", exp.input, exp.dur.Seconds(), d.Seconds())
		}
	}

}
