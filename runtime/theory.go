package runtime

import (
	"github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func pitchFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeString, 1)
	}
	return musictheory.ParsePitch(str)
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
		return nil, errors.Errorf("unknown interval quality %s", quality)
	}
}

func transposeFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	pitch, ok := args[0].(musictheory.Pitch)
	if !ok {
		return nil, lisp.ArgExpectError("pitch", 1)
	}
	interval, ok := args[1].(musictheory.Interval)
	if !ok {
		return nil, lisp.ArgExpectError("interval", 2)
	}
	return pitch.Transpose(interval), nil
}
