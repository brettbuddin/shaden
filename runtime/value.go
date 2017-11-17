package runtime

import (
	"fmt"
	"math"

	"buddin.us/musictheory"
	"buddin.us/shaden/dsp"
	"buddin.us/shaden/lisp"
)

func hzFn(args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("hz expects a string or number")
	}
	switch v := args[0].(type) {
	case float64:
		return dsp.Frequency(v), nil
	case int:
		return dsp.Frequency(float64(v)), nil
	case *musictheory.Pitch:
		return dsp.Frequency(v.Freq()), nil
	case string:
		return dsp.ParsePitch(v)
	case lisp.Keyword:
		return dsp.ParsePitch(string(v))
	default:
		return 0, fmt.Errorf("hz expects a number")
	}
}

func msFn(args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ms expects a number")
	}
	switch v := args[0].(type) {
	case float64:
		return dsp.Duration(v), nil
	case int:
		return dsp.Duration(float64(v)), nil
	default:
		return 0, fmt.Errorf("ms expects a number")
	}
}

func bpmFn(args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bpm expects a number")
	}
	switch v := args[0].(type) {
	case float64:
		return dsp.BPM(v), nil
	case int:
		return dsp.BPM(float64(v)), nil
	default:
		return 0, fmt.Errorf("bpm expects a number")
	}
}

func dbFn(args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("db expects a number")
	}

	switch v := args[0].(type) {
	case float64:
		return math.Pow(10, 0.05*v), nil
	case int:
		return math.Pow(10, 0.05*float64(v)), nil
	default:
		return 0, fmt.Errorf("db expects a number")
	}
}
