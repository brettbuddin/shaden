package runtime

import (
	"math"

	"github.com/brettbuddin/shaden/dsp"
	"github.com/brettbuddin/shaden/lisp"
)

func hzFn(sampleRate int) func(lisp.List) (any, error) {
	return func(args lisp.List) (any, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}
		switch v := args[0].(type) {
		case float64:
			return dsp.Frequency(v, sampleRate), nil
		case int:
			return dsp.Frequency(float64(v), sampleRate), nil
		case dsp.Pitch:
			return v, nil
		case string:
			return dsp.ParsePitch(v, sampleRate)
		case lisp.Keyword:
			return dsp.ParsePitch(string(v), sampleRate)
		default:
			return 0, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeString, lisp.TypeInt, lisp.TypeFloat), 1)
		}
	}
}

func msFn(sampleRate int) func(lisp.List) (any, error) {
	return func(args lisp.List) (any, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}
		f, err := lisp.ExtractFloat64(args[0], 1)
		if err != nil {
			return 0, err
		}
		return dsp.Duration(f, sampleRate), nil
	}
}

func bpmFn(sampleRate int) func(lisp.List) (any, error) {
	return func(args lisp.List) (any, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}
		f, err := lisp.ExtractFloat64(args[0], 1)
		if err != nil {
			return 0, err
		}
		return dsp.BPM(f, sampleRate), nil
	}
}

func dbFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	f, err := lisp.ExtractFloat64(args[0], 1)
	if err != nil {
		return 0, err
	}
	return math.Pow(10, 0.05*f), nil
}
