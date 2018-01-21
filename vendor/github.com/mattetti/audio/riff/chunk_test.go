package riff

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWavNextChunk(t *testing.T) {
	path, _ := filepath.Abs("fixtures/sample.wav")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	c := New(f)
	if err := c.ParseHeaders(); err != nil {
		t.Fatal(err)
	}
	// fmt
	ch, err := c.NextChunk()
	if err != nil {
		t.Fatal(err)
	}
	if ch.ID != FmtID {
		t.Fatalf("Expected the next chunk to have an ID of %q but got %q", FmtID, ch.ID)
	}
	if ch.Size != 16 {
		t.Fatalf("Expected the next chunk to have a size of %d but got %d", 16, ch.Size)
	}
	ch.Done()
	//
	ch, err = c.NextChunk()
	if err != nil {
		t.Fatal(err)
	}
	if ch.ID != DataFormatID {
		t.Fatalf("Expected the next chunk to have an ID of %q but got %q", DataFormatID, ch.ID)
	}
	if ch.Size != 53958 {
		t.Fatalf("Expected the next chunk to have a size of %d but got %d", 53958, ch.Size)
	}
	if int(c.Size) != (ch.Size + 36) {
		t.Fatal("Looks like we have some extra data in this wav file?")
	}
}

func TestNextChunk(t *testing.T) {
	path, _ := filepath.Abs("fixtures/sample.wav")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	c := New(f)
	if err := c.ParseHeaders(); err != nil {
		t.Fatal(err)
	}
	ch, err := c.NextChunk()
	if err != nil {
		t.Fatal(err)
	}
	ch.DecodeWavHeader(c)

	ch, err = c.NextChunk()
	if err != nil {
		t.Fatal(err)
	}

	nextSample := func() []byte {
		var s = make([]byte, c.BlockAlign)
		if err := ch.ReadLE(s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	firstSample := nextSample()
	if ch.Pos != int(c.BlockAlign) {
		t.Fatal("Chunk position wasn't moved as expected")
	}
	expectedSample := []byte{0, 0}
	if bytes.Compare(firstSample, expectedSample) != 0 {
		t.Fatalf("First sample doesn't seem right, got %q, expected %q", firstSample, expectedSample)
	}

	desideredPos := 1541
	bytePos := desideredPos * 2
	for ch.Pos < bytePos {
		nextSample()
	}
	s := nextSample()
	expectedSample = []byte{0xfe, 0xff}
	if bytes.Compare(s, expectedSample) != 0 {
		t.Fatalf("1542nd sample doesn't seem right, got %q, expected %q", s, expectedSample)
	}

}

func ExampleParser_NextChunk() {
	// Example showing how to access the sound data
	path, _ := filepath.Abs("fixtures/sample.wav")
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	c := New(f)
	if err := c.ParseHeaders(); err != nil {
		panic(err)
	}

	var chunk *Chunk
	for err == nil {
		chunk, err = c.NextChunk()
		if err != nil {
			panic(err)
		}
		if chunk.ID == FmtID {
			chunk.DecodeWavHeader(c)
		} else if chunk.ID == DataFormatID {
			break
		}
		chunk.Done()
	}
	soundData := chunk

	nextSample := func() []byte {
		s := make([]byte, c.BlockAlign)
		if err := soundData.ReadLE(&s); err != nil {
			panic(err)
		}
		return s
	}

	// jump to a specific sample since first samples are blank
	desideredPos := 1541
	bytePos := desideredPos * 2
	for i := 0; soundData.Pos < bytePos; i++ {
		nextSample()
		if i > soundData.Size {
			panic(fmt.Errorf("%+v read way too many bytes, we're out of bounds", soundData))
		}
	}

	sample := nextSample()
	fmt.Printf("1542nd sample: %#X %#X\n", sample[0], sample[1])
	// Output:
	// 1542nd sample: 0XFE 0XFF
}
