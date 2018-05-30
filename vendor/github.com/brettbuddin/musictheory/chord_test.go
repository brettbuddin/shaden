package musictheory

import (
	"testing"
)

func TestChordInversion(test *testing.T) {
	M3 := Major(3)
	M7 := Major(7)
	P1 := Perfect(1)
	P5 := Perfect(5)
	d5 := Diminished(5)
	m3 := Minor(3)

	majorTriad := []Interval{P1, M3, P5}
	minorTriad := []Interval{P1, m3, P5}
	diminishedMajorSeventh := []Interval{P1, m3, d5, M7}

	data := []struct {
		root      Pitch
		intervals []Interval
		degree    int
		expected  Chord
	}{
		{NewPitch(C, Natural, 0), majorTriad, 1, Chord{
			NewPitch(E, Natural, 0),
			NewPitch(G, Natural, 0),
			NewPitch(C, Natural, 1),
		}},
		{NewPitch(C, Natural, 0), majorTriad, 2, Chord{
			NewPitch(G, Natural, 0),
			NewPitch(C, Natural, 1),
			NewPitch(E, Natural, 1),
		}},
		{NewPitch(C, Natural, 0), majorTriad, 3, Chord{
			NewPitch(C, Natural, 1),
			NewPitch(E, Natural, 1),
			NewPitch(G, Natural, 1),
		}},
		{NewPitch(C, Natural, 0), majorTriad, 4, Chord{
			NewPitch(E, Natural, 1),
			NewPitch(G, Natural, 1),
			NewPitch(C, Natural, 2),
		}},
		{NewPitch(C, Natural, 0), majorTriad, 3, Chord{
			NewPitch(C, Natural, 1),
			NewPitch(E, Natural, 1),
			NewPitch(G, Natural, 1),
		}},
		{NewPitch(C, Sharp, 0), minorTriad, 2, Chord{
			NewPitch(G, Sharp, 0),
			NewPitch(C, Sharp, 1),
			NewPitch(E, Natural, 1),
		}},
		{NewPitch(E, Flat, 2), diminishedMajorSeventh, 2, Chord{
			NewPitch(A, Natural, 2),
			NewPitch(D, Natural, 3),
			NewPitch(E, Flat, 3),
			NewPitch(G, Flat, 3),
		}},
	}

	for i, t := range data {
		chord := NewChord(t.root, t.intervals).Invert(t.degree)

		for j := range chord {
			actual := chord[j]
			if !actual.Eq(t.expected[j]) {
				test.Errorf("index=%d step=%d actual=%s expected=%s", i, j, actual, t.expected[j])
			}
		}
	}
}
