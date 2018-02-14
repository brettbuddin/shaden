// this is a simple tool to validate the mp3 files in a folder
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattetti/audio/mp3"
)

var (
	pathFlag = flag.String("path", ".", "Path to look for mp3 files to validate")
)

func main() {
	flag.Parse()
	files, err := ioutil.ReadDir(*pathFlag)
	if err != nil {
		panic(err)
	}
	for _, fi := range files {
		if strings.HasSuffix(strings.ToLower(fi.Name()), ".mp3") {
			f, err := os.Open(filepath.Join(*pathFlag, fi.Name()))
			if err != nil {
				fmt.Println("error reading", fi.Name())
				panic(err)
			}
			if !mp3.SeemsValid(f) {
				fmt.Printf("%s is not valid\n", fi.Name())
				f.Seek(0, 0)
				buf := make([]byte, 128)
				f.Read(buf)
				fmt.Println(hex.Dump(buf))
			}
			f.Close()
		}
	}
}
