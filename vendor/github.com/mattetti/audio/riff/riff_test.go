package riff

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
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
		d, err := Duration(f)
		if err != nil {
			t.Fatal(err)
		}
		if d != exp.dur {
			t.Fatalf("%s of %s didn't match %f, got %f", "Duration", exp.input, exp.dur.Seconds(), d.Seconds())
		}
	}

}

func ExampleDuration() {
	path, _ := filepath.Abs("fixtures/sample.wav")
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	d, err := Duration(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File with a duration of %f seconds", d.Seconds())
	// Output:
	// File with a duration of 0.612177 seconds
}
