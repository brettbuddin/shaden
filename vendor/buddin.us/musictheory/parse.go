package musictheory

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	pitch    = regexp.MustCompile("([ABCDEFG])(bb|b|#|x)?(\\d+)")
	interval = regexp.MustCompile("([-+])?([PAdMm]|perf|maj|min|aug|dim)(\\d+)")
)

// MustParsePitch parses and returns a Pitch in scientific pitch notation or panics
func MustParsePitch(str string) Pitch {
	pitch, err := ParsePitch(str)
	if err != nil {
		panic(err)
	}
	return pitch
}

// ParsePitch parses and returns a Pitch in scientific pitch notation
func ParsePitch(str string) (Pitch, error) {
	matches := pitch.FindStringSubmatch(str)
	if len(matches) < 1 {
		return Pitch{}, fmt.Errorf("no matches found")
	}

	class := matches[1]
	modifier := matches[2]
	octave, err := strconv.Atoi(matches[3])
	if err != nil {
		return Pitch{}, err
	}

	classIndex, err := classNameIndex(class)
	if err != nil {
		return Pitch{}, err
	}

	modifierOffset, err := modifierNameOffset(modifier)
	if err != nil {
		return Pitch{}, err
	}

	return NewPitch(classIndex+1, modifierOffset, octave), nil
}

func classNameIndex(name string) (int, error) {
	for i, n := range pitchNames {
		if n == name {
			return i, nil
		}
	}

	return 0, fmt.Errorf("unknown class name: %s", name)
}

func modifierNameOffset(name string) (int, error) {
	for i, a := range modifierNames {
		if a == name {
			return i - 2, nil
		}
	}

	return 0, fmt.Errorf("unknown modifier: %s", name)
}

func ParseInterval(str string) (Interval, error) {
	matches := interval.FindStringSubmatch(str)
	if len(matches) < 1 {
		return Interval{}, fmt.Errorf("no matches found")
	}

	quality := matches[2]
	polarity := "+"
	if len(matches[1]) > 0 {
		polarity = matches[1]
	}
	step, _ := strconv.Atoi(polarity + matches[3])

	var interval Interval
	switch quality {
	case "perf":
		fallthrough
	case "P":
		interval = Perfect(step)
	case "aug":
		fallthrough
	case "A":
		interval = Augmented(step)
	case "maj":
		fallthrough
	case "M":
		interval = Major(step)
	case "min":
		fallthrough
	case "m":
		interval = Minor(step)
	case "dim":
		fallthrough
	case "d":
		interval = Diminished(step)
	default:
		return Interval{}, fmt.Errorf("invalid quality")
	}
	return interval, nil
}
