package aiff_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mattetti/audio/aiff"
)

func ExampleDecoder_Duration() {
	path, _ := filepath.Abs("fixtures/kick.aif")
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c := aiff.NewDecoder(f)
	d, _ := c.Duration()
	fmt.Printf("kick.aif has a duration of %f seconds\n", d.Seconds())
	// Output:
	// kick.aif has a duration of 0.203356 seconds
}

func ExampleDecoder_IsValidFile() {
	f, err := os.Open("fixtures/kick.aif")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Printf("is this file valid: %t", aiff.NewDecoder(f).IsValidFile())
	// Output: is this file valid: true
}
