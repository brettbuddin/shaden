package musictheory_test

import (
	"reflect"
	"testing"

	mt "buddin.us/musictheory"
	"buddin.us/musictheory/intervals"
)

type scaleTest struct {
	root      mt.Pitch
	intervals []mt.Interval
	octaves   int
	expected  []string
}

var scaleTests []scaleTest

func init() {
	scaleTests = []scaleTest{
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Chromatic, 1, []string{"C4", "Db4", "D4", "Eb4", "E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Ionian, 1, []string{"C4", "D4", "E4", "F4", "G4", "A4", "B4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Dorian, 1, []string{"C4", "D4", "Eb4", "F4", "G4", "A4", "Bb4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Phrygian, 1, []string{"C4", "Db4", "Eb4", "F4", "G4", "Ab4", "Bb4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Lydian, 1, []string{"C4", "D4", "E4", "Gb4", "G4", "A4", "B4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Mixolydian, 1, []string{"C4", "D4", "E4", "F4", "G4", "A4", "Bb4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Aeolian, 1, []string{"C4", "D4", "Eb4", "F4", "G4", "Ab4", "Bb4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Locrian, 1, []string{"C4", "Db4", "Eb4", "F4", "Gb4", "Ab4", "Bb4"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Major, 2, []string{"C4", "D4", "E4", "F4", "G4", "A4", "B4", "C5", "D5", "E5", "F5", "G5", "A5", "B5"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Minor, 1, []string{"C4", "D4", "Eb4", "F4", "G4", "Ab4", "Bb4"}},
		{mt.NewPitch(mt.E, mt.Natural, 4), intervals.MinorPentatonic, 1, []string{"E4", "G4", "A4", "B4", "D5"}},
		{mt.NewPitch(mt.E, mt.Flat, 4), intervals.MajorPentatonic, 1, []string{"Eb4", "F4", "G4", "Bb4", "C5"}},
	}
}

func TestScales(test *testing.T) {
	for i, t := range scaleTests {
		scale := mt.NewScale(t.root, t.intervals, t.octaves)
		actual := []string{}

		for _, p := range scale {
			actual = append(actual, p.Name(mt.DescNames))
		}

		if !reflect.DeepEqual(actual, t.expected) {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}
