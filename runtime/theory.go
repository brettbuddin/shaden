package runtime

import (
	"github.com/pkg/errors"

	"buddin.us/musictheory"
	"buddin.us/shaden/lisp"
)

func pitchFn(args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.Errorf("pitch expects one argument")
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, errors.Errorf("pitch expects a string for argument 1")
	}
	return musictheory.ParsePitch(str)
}

func intervalFn(args lisp.List) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.Errorf("interval expects two arguments")
	}

	var quality string
	switch v := args[0].(type) {
	case string:
		quality = v
	case lisp.Keyword:
		quality = string(v)
	default:
		return nil, errors.Errorf("interval expects a string or keyword for argument 1")
	}

	step, ok := args[1].(int)
	if !ok {
		return nil, errors.Errorf("interval expects an integer for argument 2")
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
	default:
		return nil, errors.Errorf("unknown interval quality %s", quality)
	}
}

func transposeFn(args lisp.List) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.Errorf("transpose expects two arguments")
	}
	pitch, ok := args[0].(*musictheory.Pitch)
	if !ok {
		return nil, errors.Errorf("transpose expects a pitch for argument 1")
	}
	interval, ok := args[1].(musictheory.Interval)
	if !ok {
		return nil, errors.Errorf("interval expects an interval for argument 2")
	}
	return pitch.Transpose(interval).(*musictheory.Pitch), nil
}
