package musictheory

import (
	"fmt"
	"time"
)

// Dotted makes a dotted duration
func Dotted(d Duration, dots int) Duration {
	return Duration{d.Value, dots, false}
}

// Triplet makes a triplet duration
func Triplet(d Duration) Duration {
	return Duration{d.Value, 0, true}
}

const ns = 1000000000

// Durations
var (
	D1   = Duration{1, 0, false}   // Whole
	D2   = Duration{2, 0, false}   // Half
	D4   = Duration{4, 0, false}   // Quarter
	D8   = Duration{8, 0, false}   // Eighth
	D16  = Duration{16, 0, false}  // Sixteenth
	D32  = Duration{32, 0, false}  // Thirty Second
	D64  = Duration{64, 0, false}  // Sixty Fourth
	D128 = Duration{128, 0, false} // Hundred Twenty Eighth
)

// Duration represents a note's duration
type Duration struct {
	Value   int
	Dots    int
	Triplet bool
}

// Time returns the time in nanoseconds the note's duration lasts.
// Calculated based on what unit gets the beat and what the BPM is.
func (d Duration) Time(unit Duration, bpm int) time.Duration {
	val := (60.0 / float64(bpm)) / (float64(d.Value) / 4.0) / (float64(unit.Value) / 4.0)

	for i := 0; i < d.Dots; i++ {
		val += val / 2.0
	}

	if d.Triplet {
		val = val / 3.0
	}

	return time.Duration(val * ns)
}

func (d Duration) String() string {
	var value string
	switch d.Value {
	case 1:
		value = "whole"
	case 2:
		value = "half"
	case 4:
		value = "quarter"
	case 8:
		value = "eighth"
	case 16:
		value = "sixteenth"
	case 32:
		value = "thirty-second"
	case 64:
		value = "sixty-fourth"
	case 128:
		value = "hundred-twenty-eighth"
	default:
		value = "unknown"
	}

	var dots string
	if d.Dots > 0 {
		dots = fmt.Sprintf("dotted(%d) ", d.Dots)
	}

	var triplet string
	if d.Triplet {
		triplet = "triplet "
	}

	return fmt.Sprintf("%s%s%s", dots, triplet, value)
}

// Note is a pitch with a duration
type Note struct {
	Pitch
	Duration
}

// NewNote creates a new note
func NewNote(pitch Pitch, duration Duration) Note {
	return Note{pitch, duration}
}

// Transpose transposes a note by a given interval
func (n Note) Transpose(i Interval) Note {
	return Note{n.Pitch.Transpose(i), n.Duration}
}

func (n Note) String() string {
	return fmt.Sprintf("%s %s", n.Pitch.String(), n.Duration.String())
}
