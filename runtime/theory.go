package runtime

import (
	"github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/musictheory/intervals"
	"github.com/brettbuddin/shaden/dsp"
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

const (
	typePitch = "pitch"
)

func loadTheory(env *lisp.Environment, sampleRate int) {
	env.DefineSymbol("quality/perfect", 0)
	env.DefineSymbol("quality/minor", 1)
	env.DefineSymbol("quality/major", 2)
	env.DefineSymbol("quality/diminished", 3)
	env.DefineSymbol("quality/augmented", 4)

	env.DefineSymbol("theory/interval", intervalFn)
	env.DefineSymbol("theory/transpose", transposeFn(sampleRate))
	env.DefineSymbol("theory/scale", scaleFn(sampleRate))
	env.DefineSymbol("theory/chord", chordFn(sampleRate))
}

func intervalFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}

	var quality string
	switch v := args[0].(type) {
	case string:
		quality = v
	case lisp.Keyword:
		quality = string(v)
	default:
		return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeString, lisp.TypeKeyword), 1)
	}

	step, ok := args[1].(int)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeInt, 2)
	}

	switch quality {
	case "semitone":
		return musictheory.Semitones(step), nil
	case "perfect":
		return musictheory.Perfect(step), nil
	case "minor", "min":
		return musictheory.Minor(step), nil
	case "major", "maj":
		return musictheory.Major(step), nil
	case "augmented", "aug":
		return musictheory.Augmented(step), nil
	case "diminished", "dim":
		return musictheory.Diminished(step), nil
	case "octave":
		return musictheory.Octave(step), nil
	default:
		return nil, errors.Errorf("unknown interval %s", quality)
	}
}

func transposeFn(sampleRate int) func(args lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 2); err != nil {
			return nil, err
		}
		pitch, ok := args[0].(dsp.Pitch)
		if !ok {
			return nil, lisp.ArgExpectError(typePitch, 1)
		}
		interval, ok := args[1].(musictheory.Interval)
		if !ok {
			return nil, lisp.ArgExpectError("interval", 2)
		}
		return pitch.Transpose(interval, sampleRate), nil
	}
}

func scaleFn(sampleRate int) func(args lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 3); err != nil {
			return nil, err
		}
		root, ok := args[0].(dsp.Pitch)
		if !ok {
			return nil, lisp.ArgExpectError(typePitch, 1)
		}

		var name string
		switch v := args[1].(type) {
		case string:
			name = v
		case lisp.Keyword:
			name = string(v)
		default:
			return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeString, lisp.TypeKeyword), 2)
		}

		octaves, ok := args[2].(int)
		if !ok {
			return nil, lisp.ArgExpectError(lisp.TypeInt, 3)
		}

		itvls, err := nameToScale(name)
		if err != nil {
			return nil, err
		}

		var list lisp.List
		for _, p := range musictheory.NewScale(root.Pitch, itvls, octaves) {
			list = append(list, dsp.NewPitch(p, sampleRate))
		}
		return list, nil
	}
}

func chordFn(sampleRate int) func(args lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 2); err != nil {
			return nil, err
		}
		root, ok := args[0].(dsp.Pitch)
		if !ok {
			return nil, lisp.ArgExpectError(typePitch, 1)
		}

		var name string
		switch v := args[1].(type) {
		case string:
			name = v
		case lisp.Keyword:
			name = string(v)
		default:
			return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeString, lisp.TypeKeyword), 2)
		}

		itvls, err := nameToChord(name)
		if err != nil {
			return nil, err
		}

		var list lisp.List
		for _, p := range musictheory.NewChord(root.Pitch, itvls) {
			list = append(list, dsp.NewPitch(p, sampleRate))
		}
		return list, nil
	}
}

func nameToScale(name string) ([]musictheory.Interval, error) {
	switch name {
	case "aeolian":
		return intervals.Aeolian, nil
	case "chromatic":
		return intervals.Chromatic, nil
	case "dominant-bebop":
		return intervals.DominantBebop, nil
	case "dorian":
		return intervals.Dorian, nil
	case "double-harmonic":
		return intervals.DoubleHarmonic, nil
	case "harmonic-minor":
		return intervals.HarmonicMinor, nil
	case "in-sen":
		return intervals.InSen, nil
	case "ionian":
		return intervals.Ionian, nil
	case "locrian":
		return intervals.Locrian, nil
	case "lydian":
		return intervals.Lydian, nil
	case "major":
		return intervals.Major, nil
	case "major-bebop":
		return intervals.MajorBebop, nil
	case "major-pentatonic":
		return intervals.MajorPentatonic, nil
	case "melodic-minor-bebop":
		return intervals.MelodicMinorBebop, nil
	case "minor":
		return intervals.Minor, nil
	case "minor-pentatonic":
		return intervals.MinorPentatonic, nil
	case "mixolydian":
		return intervals.Mixolydian, nil
	case "phrygian":
		return intervals.Phrygian, nil
	case "whole-tone":
		return intervals.WholeTone, nil
	}
	return nil, errors.Errorf("unknown scale %s", name)
}

func nameToChord(name string) ([]musictheory.Interval, error) {
	switch name {
	case "augmented-major-seventh", "augM7":
		return intervals.AugmentedMajorSeventh, nil
	case "augmented-seventh", "aug7":
		return intervals.AugmentedSeventh, nil
	case "augmented-sixth", "aug6":
		return intervals.AugmentedSixth, nil
	case "augmented", "aug":
		return intervals.AugmentedTriad, nil
	case "diminished-major-seventh", "dimM7":
		return intervals.DiminishedMajorSeventh, nil
	case "diminished-seventh", "dim7":
		return intervals.DiminishedSeventh, nil
	case "diminished", "dim":
		return intervals.DiminishedTriad, nil
	case "dominant-seventh", "7":
		return intervals.DominantSeventh, nil
	case "half-diminished-seventh", "min7b5", "m7b5":
		return intervals.HalfDiminishedSeventh, nil
	case "major-seventh", "maj7", "M7":
		return intervals.MajorSeventh, nil
	case "major-sixth", "maj6", "M6":
		return intervals.MajorSixth, nil
	case "major", "maj", "M":
		return intervals.MajorTriad, nil
	case "minor-seventh", "min7", "m7":
		return intervals.MinorSeventh, nil
	case "minor-sixth", "min6", "m6":
		return intervals.MinorSixth, nil
	case "minor", "min", "m":
		return intervals.MinorTriad, nil
	}
	return nil, errors.Errorf("unknown chord %s", name)
}
