package builtin

import (
	"math"
	"math/rand"

	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func multFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	var (
		outi      = int(1)
		outf      = float64(1)
		seenFloat bool
	)
	for _, arg := range args {
		switch arg := arg.(type) {
		case int:
			outi *= arg
			outf *= float64(arg)
		case float64:
			outf *= arg
			seenFloat = true
		default:
			return nil, errors.Errorf("cannot multiply %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func divFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	var (
		arity     = len(args)
		outi      int
		outf      float64
		seenFloat bool
	)
	for i, arg := range args {
		switch arg := arg.(type) {
		case int:
			if i == 0 && arity > 1 {
				outi = arg
				outf = float64(arg)
			} else {
				outi /= arg
				outf /= float64(arg)
			}
		case float64:
			if i == 0 && arity > 1 {
				outf = float64(arg)
			} else {
				outf /= arg
			}
			seenFloat = true
		default:
			return nil, errors.Errorf("cannot divide %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func sumFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	var (
		outi      int
		outf      float64
		seenFloat bool
	)
	for _, arg := range args {
		switch arg := arg.(type) {
		case int:
			outi += arg
			outf += float64(arg)
		case float64:
			outf += arg
			seenFloat = true
		default:
			return nil, errors.Errorf("cannot add %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func diffFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	var (
		arity     = len(args)
		outi      int
		outf      float64
		seenFloat bool
	)
	for i, arg := range args {
		switch arg := arg.(type) {
		case int:
			if i == 0 && arity > 1 {
				outi = arg
				outf = float64(arg)
			} else {
				outi -= arg
				outf -= float64(arg)
			}
		case float64:
			if i == 0 && arity > 1 {
				outf = float64(arg)
			} else {
				outf -= arg
			}
			seenFloat = true
		default:
			return nil, errors.Errorf("cannot subtract %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func powFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}

	var x, y float64

	if v, ok := args[0].(int); ok {
		x = float64(v)
	} else if f, ok := args[0].(float64); ok {
		x = f
	} else {
		return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeInt, lisp.TypeFloat), 1)
	}

	if v, ok := args[1].(int); ok {
		y = float64(v)
	} else if f, ok := args[1].(float64); ok {
		y = f
	} else {
		return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeInt, lisp.TypeFloat), 2)
	}

	return math.Pow(x, y), nil
}

func randFn(args lisp.List) (value interface{}, err error) {
	if err := lisp.CheckArityEqual(args, 0); err != nil {
		return nil, err
	}
	return rand.Float64(), nil
}
