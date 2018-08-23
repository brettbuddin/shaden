package musictheory_test

import (
	"reflect"
	"testing"

	mt "github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/musictheory/intervals"
)

type scaleTest struct {
	root      mt.Pitch
	intervals []mt.Interval
	octaves   int
	expected  []string
}

func TestScales(t *testing.T) {
	tests := []scaleTest{
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Chromatic, 1, []string{"C4", "Db4", "D4", "Eb4", "E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4"}},
		{mt.NewPitch(mt.C, mt.Sharp, 4), intervals.Chromatic, 1, []string{"Db4", "D4", "Eb4", "E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5"}},
		{mt.NewPitch(mt.D, mt.Natural, 4), intervals.Chromatic, 1, []string{"D4", "Eb4", "E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5"}},
		{mt.NewPitch(mt.E, mt.Flat, 4), intervals.Chromatic, 1, []string{"Eb4", "E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5"}},
		{mt.NewPitch(mt.E, mt.Natural, 4), intervals.Chromatic, 1, []string{"E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5"}},
		{mt.NewPitch(mt.F, mt.Natural, 4), intervals.Chromatic, 1, []string{"F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5"}},
		{mt.NewPitch(mt.F, mt.Sharp, 4), intervals.Chromatic, 1, []string{"Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5", "F5"}},
		{mt.NewPitch(mt.G, mt.Natural, 4), intervals.Chromatic, 1, []string{"G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5", "F5", "Gb5"}},
		{mt.NewPitch(mt.A, mt.Flat, 4), intervals.Chromatic, 1, []string{"Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5", "F5", "Gb5", "G5"}},
		{mt.NewPitch(mt.A, mt.Natural, 4), intervals.Chromatic, 1, []string{"A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5", "F5", "Gb5", "G5", "Ab5"}},
		{mt.NewPitch(mt.B, mt.Flat, 4), intervals.Chromatic, 1, []string{"Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5", "F5", "Gb5", "G5", "Ab5", "A5"}},
		{mt.NewPitch(mt.B, mt.Natural, 4), intervals.Chromatic, 1, []string{"B4", "C5", "Db5", "D5", "Eb5", "E5", "F5", "Gb5", "G5", "Ab5", "A5", "Bb5"}},

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
		{mt.NewPitch(mt.E, mt.Flat, 4), intervals.WholeTone, 1, []string{"Eb4", "F4", "G4", "A4", "B4", "Db5"}},

		{mt.NewPitch(mt.E, mt.Flat, 4), intervals.WholeTone, -1, []string{"Eb4", "Db4", "B3", "A3", "G3", "F3"}},
		{mt.NewPitch(mt.C, mt.Natural, 4), intervals.Ionian, -2, []string{"C4", "B3", "A3", "G3", "F3", "E3", "D3", "C3", "B2", "A2", "G2", "F2", "E2", "D2"}},
	}
	for i, tt := range tests {
		scale := mt.NewScale(tt.root, tt.intervals, tt.octaves)
		actual := []string{}

		for _, p := range scale {
			actual = append(actual, p.Name(mt.DescNames))
		}

		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("index=%d actual=%s expected=%s", i, actual, tt.expected)
		}
	}
}
