package musictheory

import (
	"testing"
	"time"
)

func TestNoteDuration(test *testing.T) {
	data := []struct {
		duration Duration
		unit     Duration
		bpm      int
		expected time.Duration
	}{
		{D1, D4, 60, 4 * time.Second},
		{D4, D4, 60, 1 * time.Second},
		{D8, D8, 60, 250 * time.Millisecond},
		{D1, D8, 60, 2 * time.Second},
		{D1, D8, 60, 2 * time.Second},
		{Triplet(D4), D4, 60, (1 * time.Second / 3.0)},
		{Dotted(D4, 1), D4, 60, 1500 * time.Millisecond},
		{Dotted(D4, 2), D4, 60, 2250 * time.Millisecond},
	}

	for i, t := range data {
		actual := NewNote(NewPitch(C, Natural, 0), t.duration).Time(t.unit, t.bpm)

		if actual != t.expected {
			test.Errorf("index=%d actual=%s expected=%s", i, actual, t.expected)
		}
	}
}
