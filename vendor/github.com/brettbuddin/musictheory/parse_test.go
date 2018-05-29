package musictheory

import (
	"testing"
)

func TestParsePitch(test *testing.T) {
	data := []struct {
		input    string
		expected Pitch
	}{
		{"C4", NewPitch(C, Natural, 4)},
		{"C#4", NewPitch(C, Sharp, 4)},
		{"Ab3", NewPitch(A, Flat, 3)},
		{"Abb3", NewPitch(A, DoubleFlat, 3)},
		{"Cx4", NewPitch(C, DoubleSharp, 4)},
	}

	for i, t := range data {
		actual, err := ParsePitch(t.input)
		if err != nil {
			test.Error(err)
		}

		if !actual.Eq(t.expected) {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}

func TestParseInterval(test *testing.T) {
	data := []struct {
		input    string
		expected Interval
	}{
		{"P4", Perfect(4)},
		{"P5", Perfect(5)},
		{"perf5", Perfect(5)},
		{"-P5", Perfect(-5)},
		{"-perf5", Perfect(-5)},
		{"M3", Major(3)},
		{"maj3", Major(3)},
		{"m3", Minor(3)},
		{"min3", Minor(3)},
		{"A4", Augmented(4)},
		{"aug4", Augmented(4)},
		{"d5", Diminished(5)},
		{"-d5", Diminished(-5)},
		{"-dim5", Diminished(-5)},
	}

	for i, t := range data {
		actual, err := ParseInterval(t.input)
		if err != nil {
			test.Error(err)
		}

		if !actual.Eq(t.expected) {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}
