package main

import (
	"flag"
	"fmt"
	"os"

	mt "buddin.us/musictheory"
)

const help = `transpose: tool for music pitch transposition

Usage:

transpose -pitch <pitch> -interval [<polarity>]<interval>

Pitches:

Pitches are expressed in scientific pitch notation: pitch class + octave. The valid pitch classes are:

C, C#, Db, D, D#, Eb, E, F, F#, Gb, G, G#, Ab, A, Bb, B

Double flats (bb) and double sharps (x) are also valid.

Intervals:

Intervals are expressed as a quality + step.

Name	    Values
----------  -------
Perfect     P, perf
Major       M, maj
Minor       m, min
Augmented   A, aug
Diminished  d, dim

Providing a negative interval (e.g. "-d5") will transpose down by that interval.

Examples:

transpose -pitch C3 -interval m3
transpose -pitch F#4 -interval +P5
transpose -pitch Bb5 -interval -m2

Flags:
`

var (
	pitchFlag    = flag.String("pitch", "", "pitch")
	intervalFlag = flag.String("interval", "", "interval")
	namingFlag   = flag.String("naming", "asc", "ascending or descending naming strategy")
	helpFlag     = flag.Bool("help", false, "show usage message")
)

func main() {
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "transpose: too many arguments")
		os.Exit(1)
	}

	if *helpFlag || *pitchFlag == "" || *intervalFlag == "" {
		fmt.Println(help)
		flag.PrintDefaults()
		return
	}

	tonic, err := mt.ParsePitch(*pitchFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse pitch: %v\n", err)
		os.Exit(1)
	}
	interval, err := mt.ParseInterval(*intervalFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse interval: %v\n", err)
		os.Exit(1)
	}

	var strategy mt.ModifierStrategy
	switch *namingFlag {
	case "asc":
		strategy = mt.AscNames
	case "desc":
		strategy = mt.DescNames
	default:
		fmt.Fprintf(os.Stderr, "unknown naming strategy: %v\n", *namingFlag)
		os.Exit(1)
	}

	fmt.Println(tonic.Transpose(interval).Name(strategy))
}
