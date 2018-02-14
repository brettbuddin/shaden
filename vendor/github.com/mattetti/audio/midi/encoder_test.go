package midi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewEncoder(t *testing.T) {
	w, err := tmpFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		w.Close()
		os.Remove(w.Name())
	}()
	e := NewEncoder(w, SingleTrack, 96)
	tr := e.NewTrack()
	// add a C3 at velocity 99, half a beat/quarter note after the start
	tr.Add(0.5, NoteOn(1, KeyInt("C", 3), 99))
	// turn off the C3
	tr.Add(1, NoteOff(1, KeyInt("C", 3)))
	if err := e.Write(); err != nil {
		t.Fatal(err)
	}
	w.Seek(0, 0)
	midiData, err := ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
	}
	expected := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 00, 00, 00, 0x01, 00, 0x60, 0x4d, 0x54,
		0x72, 0x6b, 0x00, 0x00, 0x00, 0x14, 0x00, 0xff, 0x58, 0x04, 0x04, 0x02, 0x24, 0x08, 0x30, 0x91,
		0x3c, 0x63, 0x60, 0x81, 0x3c, 0x40, 0x00, 0xff, 0x2f, 0x00,
	}
	if bytes.Compare(midiData, expected) != 0 {
		t.Logf("\nExpected:\t%#v\nGot:\t\t%#v\n", expected, midiData)
		t.Fatal(fmt.Errorf("Midi binary output didn't match expectations"))
	}
}

func TestNewEncoderWithForcedEnd(t *testing.T) {
	w, err := tmpFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		w.Close()
		os.Remove(w.Name())
	}()
	e := NewEncoder(w, SingleTrack, 96)
	tr := e.NewTrack()
	// add a C3 at velocity 99, half a beat/quarter note after the start
	tr.Add(0.5, NoteOn(1, KeyInt("C", 3), 99))
	// turn off the C3
	tr.Add(1, NoteOff(1, KeyInt("C", 3)))
	// force the end of track to be later
	tr.Add(2, EndOfTrack())
	if err := e.Write(); err != nil {
		t.Fatal(err)
	}

	w.Seek(0, 0)
	midiData, err := ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
	}
	expected := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 00, 00, 00, 0x01, 00, 0x60, 0x4d, 0x54,
		0x72, 0x6b, 0x00, 0x00, 0x00, 0x15, 0x00, 0xff, 0x58, 0x04, 0x04, 0x02, 0x24, 0x08, 0x30, 0x91,
		0x3c, 0x63, 0x60, 0x81, 0x3c, 0x40, 0x81, 0x40, 0xff, 0x2f, 0x0,
	}
	if bytes.Compare(midiData, expected) != 0 {
		t.Logf("\nExpected:\t%#v\nGot:\t\t%#v\n", expected, midiData)
		t.Fatal(fmt.Errorf("Midi binary output didn't match expectations"))
	}
}

func tmpFile() (*os.File, error) {
	f, err := ioutil.TempFile("", "midi-test-")
	if err != nil {
		return nil, err
	}
	return f, nil
}
