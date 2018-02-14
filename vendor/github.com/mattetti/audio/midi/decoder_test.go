package midi

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestVarint(t *testing.T) {
	expecations := []struct {
		dec   uint32
		bytes []byte
	}{
		{0, []byte{0}},
		{42, []byte{0x2a}},
		{4610, []byte{0xa4, 0x02}},
	}

	for _, exp := range expecations {
		conv := EncodeVarint(exp.dec)
		if bytes.Compare(conv, exp.bytes) != 0 {
			t.Fatalf("%d was converted to %#v didn't match %#v\n", exp.dec, conv, exp.bytes)
		}
	}

	for _, exp := range expecations {
		conv, _ := DecodeVarint(exp.bytes)
		if conv != exp.dec {
			t.Fatalf("%#v was converted to %d didn't match %d\n", exp.bytes, conv, exp.dec)
		}
	}
}

func TestParsingFile(t *testing.T) {
	expectations := []struct {
		path                string
		format              uint16
		numTracks           uint16
		ticksPerQuarterNote uint16
		timeFormat          timeFormat
		trackNames          []string
		bpms                []int
	}{
		{"fixtures/elise.mid", 1, 4, 960, MetricalTF, []string{"Track 0", "F\xfcr Elise", "http://www.forelise.com/", ""}, []int{69, 0, 0, 0}},
		{"fixtures/elise1track.mid", 0, 1, 480, MetricalTF, []string{"F"}, []int{69}},
		{"fixtures/bossa.mid", 0, 1, 96, MetricalTF, []string{"bossa 1"}, []int{0}},
		{"fixtures/closedHat.mid", 0, 1, 96, MetricalTF, []string{"01 4th Hat Closed Side"}, []int{0}},
	}

	for _, exp := range expectations {
		t.Log(exp.path)
		path, _ := filepath.Abs(exp.path)
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		p := NewDecoder(f)
		if err := p.Parse(); err != nil {
			t.Fatal(err)
		}

		if p.Format != exp.format {
			t.Fatalf("%s of %s didn't match %v, got %v", "format", exp.path, exp.format, p.Format)
		}
		if p.NumTracks != exp.numTracks {
			t.Fatalf("%s of %s didn't match %v, got %v", "numTracks", exp.path, exp.numTracks, p.NumTracks)
		}
		if p.TicksPerQuarterNote != exp.ticksPerQuarterNote {
			t.Fatalf("%s of %s didn't match %v, got %v", "ticksPerQuarterNote", exp.path, exp.ticksPerQuarterNote, p.TicksPerQuarterNote)
		}
		if p.TimeFormat != exp.timeFormat {
			t.Fatalf("%s of %s didn't match %v, got %v", "format", exp.path, exp.timeFormat, p.TimeFormat)
		}

		if len(p.Tracks) == 0 {
			t.Fatal("Tracks not parsed")
		}
		t.Logf("%d tracks\n", len(p.Tracks))
		for i, tr := range p.Tracks {
			t.Log("track", i)
			if tName := tr.Name(); tName != exp.trackNames[i] {
				t.Fatalf("expected name of track %d to be %s but got %s (%q)", i, exp.trackNames[i], tName, tName)
			}
			if bpm := tr.Tempo(); bpm != exp.bpms[i] {
				t.Fatalf("expected tempo of track %d to be %d but got %d", i, exp.bpms[i], bpm)
			}
			for _, ev := range tr.Events {
				t.Log(ev)
			}
		}

	}
}
