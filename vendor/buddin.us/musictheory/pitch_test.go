package musictheory

import (
	"math"
	"testing"
)

const tolerance = 0.000001

func TestSharpPitchNames(test *testing.T) {
	data := []struct {
		pitch    int
		expected string
	}{
		{C, "C#4"},
		{D, "D#4"},
		{E, "F4"},
		{F, "F#4"},
		{G, "G#4"},
		{A, "A#4"},
		{B, "C4"},
	}

	for i, t := range data {
		actual := NewPitch(t.pitch, Sharp, 4).Name(AscNames)
		if actual != t.expected {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}

func TestFlatPitchNames(test *testing.T) {
	data := []struct {
		pitch    int
		expected string
	}{
		{C, "B4"},
		{D, "Db4"},
		{E, "Eb4"},
		{F, "E4"},
		{G, "Gb4"},
		{A, "Ab4"},
		{B, "Bb4"},
	}

	for i, t := range data {
		actual := NewPitch(t.pitch, Flat, 4).Name(DescNames)
		if actual != t.expected {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}

func TestDoubleSharpPitchNames(test *testing.T) {
	data := []struct {
		pitch    int
		expected string
	}{
		{C, "D4"},
		{D, "E4"},
		{E, "F#4"},
		{F, "G4"},
		{G, "A4"},
		{A, "B4"},
		{B, "C#4"},
	}

	for i, t := range data {
		actual := NewPitch(t.pitch, DoubleSharp, 4).Name(AscNames)
		if actual != t.expected {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}

func TestDoubleFlatPitchNames(test *testing.T) {
	data := []struct {
		pitch    int
		expected string
	}{
		{C, "Bb4"},
		{D, "C4"},
		{E, "D4"},
		{F, "Eb4"},
		{G, "F4"},
		{A, "G4"},
		{B, "A4"},
	}

	for i, t := range data {
		actual := NewPitch(t.pitch, DoubleFlat, 4).Name(DescNames)
		if actual != t.expected {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}

func TestFrequency(test *testing.T) {
	data := []struct {
		input    Pitch
		expected float64
	}{
		{NewPitch(C, Natural, 4), 261.625565},
		{NewPitch(A, Natural, 4), 440.0},
	}

	for _, t := range data {
		actual := t.input.Freq()

		if closeEqualFloat64(actual, t.expected) {
			test.Errorf("input=%s output=%f, expected=%f",
				t.input,
				actual,
				t.expected)
		}
	}
}

func TestMIDI(test *testing.T) {
	data := []struct {
		input    Pitch
		expected int
	}{
		{NewPitch(C, Natural, 3), 60},
		{NewPitch(C, Natural, 4), 72},
		{NewPitch(A, Natural, 4), 81},
	}

	for _, t := range data {
		actual := t.input.MIDI()

		if actual != t.expected {
			test.Errorf("input=%s output=%f, expected=%f",
				t.input,
				actual,
				t.expected)
		}
	}
}

func closeEqualFloat64(actual, expected float64) bool {
	return math.Abs(actual-expected) >= tolerance
}

func TestNearestPitch(test *testing.T) {
	data := []struct {
		input    float64
		expected Pitch
	}{
		{74, NewPitch(D, Natural, 2)},
		{190, NewPitch(F, Sharp, 3)},
		{400, NewPitch(G, Natural, 4)},
		{350, NewPitch(F, Natural, 4)},
		{concertFrequency, NewPitch(A, Natural, 4)},
		{800, NewPitch(G, Natural, 5)},
		{32000, NewPitch(B, Natural, 10)},
	}

	for _, t := range data {
		actual := NearestPitch(t.input)
		if !actual.Eq(t.expected) {
			test.Errorf("input=%f output=%v, expected=%v",
				t.input,
				actual,
				t.expected)
		}
	}
}
