package runtime

import (
	"log"
	"os"
	"testing"

	"github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/musictheory/intervals"
	"github.com/brettbuddin/shaden/dsp"
	"github.com/brettbuddin/shaden/engine"
	"github.com/brettbuddin/shaden/lisp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterval(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(newBackend(0), frameSize, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	tests := []struct {
		input    string
		expected musictheory.Interval
	}{
		{"(theory/interval :perfect 1)", musictheory.Perfect(1)},
		{"(theory/interval :minor 2)", musictheory.Minor(2)},
		{"(theory/interval :major 3)", musictheory.Major(3)},
		{"(theory/interval :augmented 4)", musictheory.Augmented(4)},
		{"(theory/interval :diminished 5)", musictheory.Diminished(5)},
		{`(theory/interval "diminished" 5)`, musictheory.Diminished(5)},
	}

	for _, test := range tests {
		v, err := run.Eval([]byte(test.input))
		require.NoError(t, err)
		require.Equal(t, test.expected, v)
	}

	_, err = run.Eval([]byte(`(theory/interval :unknown)`))
	require.Error(t, err)

	_, err = run.Eval([]byte(`(theory/interval)`))
	require.Error(t, err)
}

func TestTranspose(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(newBackend(0), sampleRate, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	v, err := run.Eval([]byte(`(theory/transpose (hz "A4") (theory/interval :minor 3))`))
	require.NoError(t, err)

	actualPitch, err := musictheory.ParsePitch("C5")
	require.NoError(t, err)
	require.True(t, actualPitch.Eq(v.(dsp.Pitch).Pitch))
}

func TestScale(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(newBackend(0), frameSize, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	root, _ := musictheory.ParsePitch("C4")

	newScale := func(intvls []musictheory.Interval, octaves int) lisp.List {
		var list lisp.List
		for _, v := range musictheory.NewScale(root, intvls, octaves) {
			list = append(list, v)
		}
		return list
	}

	tests := []struct {
		input    string
		expected lisp.List
	}{
		{`(theory/scale (hz "C4") "major" 1)`, newScale(intervals.Major, 1)},
		{`(theory/scale (hz "C4") "minor" 1)`, newScale(intervals.Minor, 1)},
		{`(theory/scale (hz "C4") "aeolian" 1)`, newScale(intervals.Aeolian, 1)},
		{`(theory/scale (hz "C4") "chromatic" 1)`, newScale(intervals.Chromatic, 1)},
		{`(theory/scale (hz "C4") "dominant-bebop" 1)`, newScale(intervals.DominantBebop, 1)},
		{`(theory/scale (hz "C4") "dorian" 1)`, newScale(intervals.Dorian, 1)},
		{`(theory/scale (hz "C4") "double-harmonic" 1)`, newScale(intervals.DoubleHarmonic, 1)},
		{`(theory/scale (hz "C4") "in-sen" 1)`, newScale(intervals.InSen, 1)},
		{`(theory/scale (hz "C4") "ionian" 1)`, newScale(intervals.Ionian, 1)},
		{`(theory/scale (hz "C4") "locrian" 1)`, newScale(intervals.Locrian, 1)},
		{`(theory/scale (hz "C4") "lydian" 1)`, newScale(intervals.Lydian, 1)},
		{`(theory/scale (hz "C4") "major-bebop" 1)`, newScale(intervals.MajorBebop, 1)},
		{`(theory/scale (hz "C4") "major-pentatonic" 1)`, newScale(intervals.MajorPentatonic, 1)},
		{`(theory/scale (hz "C4") "melodic-minor-bebop" 1)`, newScale(intervals.MelodicMinorBebop, 1)},
		{`(theory/scale (hz "C4") "minor-pentatonic" 1)`, newScale(intervals.MinorPentatonic, 1)},
		{`(theory/scale (hz "C4") "mixolydian" 1)`, newScale(intervals.Mixolydian, 1)},
		{`(theory/scale (hz "C4") "phrygian" 1)`, newScale(intervals.Phrygian, 1)},
		{`(theory/scale (hz "C4") "whole-tone" 1)`, newScale(intervals.WholeTone, 1)},
		{`(theory/scale (hz "C4") "whole-tone" 2)`, newScale(intervals.WholeTone, 2)},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			iactual, err := run.Eval([]byte(test.input))
			require.NoError(t, err)

			actual := iactual.(lisp.List)

			for i, expected := range test.expected {
				if i >= len(actual) {
					t.Fail()
				}
				assert.Equal(t, expected, actual[i].(dsp.Pitch).Pitch)
			}
		})
	}
}

func TestChord(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(newBackend(0), frameSize, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	root, _ := musictheory.ParsePitch("C4")

	newChord := func(intvls []musictheory.Interval) lisp.List {
		var list lisp.List
		for _, v := range musictheory.NewChord(root, intvls) {
			list = append(list, v)
		}
		return list
	}

	tests := []struct {
		input    string
		expected lisp.List
	}{
		{`(theory/chord (hz "C4") "major")`, newChord(intervals.MajorTriad)},
		{`(theory/chord (hz "C4") "minor")`, newChord(intervals.MinorTriad)},
		{`(theory/chord (hz "C4") "augmented-major-seventh")`, newChord(intervals.AugmentedMajorSeventh)},
		{`(theory/chord (hz "C4") "augmented-seventh")`, newChord(intervals.AugmentedSeventh)},
		{`(theory/chord (hz "C4") "augmented-sixth")`, newChord(intervals.AugmentedSixth)},
		{`(theory/chord (hz "C4") "augmented")`, newChord(intervals.AugmentedTriad)},
		{`(theory/chord (hz "C4") "diminished-major-seventh")`, newChord(intervals.DiminishedMajorSeventh)},
		{`(theory/chord (hz "C4") "diminished-seventh")`, newChord(intervals.DiminishedSeventh)},
		{`(theory/chord (hz "C4") "diminished")`, newChord(intervals.DiminishedTriad)},
		{`(theory/chord (hz "C4") "dominant-seventh")`, newChord(intervals.DominantSeventh)},
		{`(theory/chord (hz "C4") "half-diminished-seventh")`, newChord(intervals.HalfDiminishedSeventh)},
		{`(theory/chord (hz "C4") "major-seventh")`, newChord(intervals.MajorSeventh)},
		{`(theory/chord (hz "C4") "major-sixth")`, newChord(intervals.MajorSixth)},
		{`(theory/chord (hz "C4") "minor-seventh")`, newChord(intervals.MinorSeventh)},
		{`(theory/chord (hz "C4") "minor-sixth")`, newChord(intervals.MinorSixth)},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			iactual, err := run.Eval([]byte(test.input))
			require.NoError(t, err)

			actual := iactual.(lisp.List)

			for i, expected := range test.expected {
				if i >= len(actual) {
					t.Fail()
				}
				assert.Equal(t, expected, actual[i].(dsp.Pitch).Pitch)
			}
		})
	}
}
