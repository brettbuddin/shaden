package runtime

import (
	"math"

	"github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/shaden/dsp"
	"github.com/brettbuddin/shaden/lisp"
)

func hzFn(sampleRate int) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}
		switch v := args[0].(type) {
		case float64:
			return dsp.Frequency(v, sampleRate), nil
		case int:
			return dsp.Frequency(float64(v), sampleRate), nil
		case musictheory.Pitch:
			return dsp.Frequency(v.Freq(), sampleRate), nil
		case string:
			return dsp.ParsePitch(v, sampleRate)
		case lisp.Keyword:
			return dsp.ParsePitch(string(v), sampleRate)
		default:
			return 0, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeString, lisp.TypeInt, lisp.TypeFloat), 1)
		}
	}
}

func msFn(sampleRate int) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}
		switch v := args[0].(type) {
		case float64:
			return dsp.Duration(v, sampleRate), nil
		case int:
			return dsp.Duration(float64(v), sampleRate), nil
		default:
			return 0, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeInt, lisp.TypeFloat), 1)
		}
	}
}

func bpmFn(sampleRate int) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}
		switch v := args[0].(type) {
		case float64:
			return dsp.BPM(v, sampleRate), nil
		case int:
			return dsp.BPM(float64(v), sampleRate), nil
		default:
			return 0, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeInt, lisp.TypeFloat), 1)
		}
	}
}

func dbFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case float64:
		return math.Pow(10, 0.05*v), nil
	case int:
		return math.Pow(10, 0.05*float64(v)), nil
	default:
		return 0, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeInt, lisp.TypeFloat), 1)
	}
}
