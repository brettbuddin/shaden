package caf

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestBadFileHeaderData(t *testing.T) {
	r := bytes.NewReader([]byte{'m', 'a', 't', 't', 0, 0, 0})
	d := New(r)
	if err := d.Parse(); err == nil {
		t.Fatalf("Expected bad data to return %s", ErrFmtNotSupported)
	}

	r = bytes.NewReader([]byte{'c', 'a', 'f', 'f', 2, 0, 0})
	d = New(r)
	if err := d.Parse(); err == nil {
		t.Fatalf("Expected bad data to return %s", ErrFmtNotSupported)
	}
}

func TestParsingFile(t *testing.T) {
	expectations := []struct {
		path    string
		format  [4]byte
		version uint16
		flags   uint16
	}{
		{"fixtures/ring.caf", fileHeaderID, 1, 0},
		{"fixtures/bass.caf", fileHeaderID, 1, 0},
	}

	for _, exp := range expectations {
		t.Log(exp.path)
		path, _ := filepath.Abs(exp.path)
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		d := New(f)
		if err := d.Parse(); err != nil {
			t.Fatal(err)
		}
		fmt.Println(d)
		t.Logf("%+v\n", *d)

		if d.Format != exp.format {
			t.Fatalf("%s of %s didn't match %v, got %v", "format", exp.path, exp.format, d.Format)
		}
		if d.Version != exp.version {
			t.Fatalf("%s of %s didn't match %d, got %v", "version", exp.path, exp.version, d.Version)
		}
		if d.Flags != exp.flags {
			t.Fatalf("%s of %s didn't match %d, got %v", "flags", exp.path, exp.flags, d.Flags)
		}

	}
}
