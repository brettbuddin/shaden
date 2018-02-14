// riffinfo is a command line tool used to gather information about riff files.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattetti/audio/riff"
)

var pathToParse = flag.String("path", ".", "Where to find wav files")
var fileToParse = flag.String("file", "", "The wav file to analyze (instead of a path)")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: \n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *fileToParse != "" {
		analyze(*fileToParse)
	} else {
		err := filepath.Walk(*pathToParse, walkFn)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func walkFn(path string, fi os.FileInfo, err error) (e error) {
	if err != nil {
		log.Fatal(err)
	}
	if fi.IsDir() {
		filepath.Walk(path, walkFolder)
		return
	}
	if !strings.HasSuffix(fi.Name(), ".wav") || fi.IsDir() {
		return
	}
	return nil
}

func walkFolder(path string, fi os.FileInfo, err error) (e error) {
	if !strings.HasSuffix(fi.Name(), ".wav") || fi.IsDir() {
		return
	}
	analyze(path)
	return nil
}

func analyze(path string) {
	fmt.Println(path)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	c := riff.New(f)
	if err := c.ParseHeaders(); err != nil {
		log.Fatalf("Can't parse the headers of %s - %s\n", path, err)
	}
	fmt.Println(c)
}
