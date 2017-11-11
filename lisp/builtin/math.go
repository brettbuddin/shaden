package builtin

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"buddin.us/lumen/lisp"
)

func multFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "*", 2); err != nil {
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
			return nil, fmt.Errorf("cannot multiply %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func divFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "/", 2); err != nil {
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
			return nil, fmt.Errorf("cannot divide with %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func sumFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "+", 2); err != nil {
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
			return nil, fmt.Errorf("cannot add %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func diffFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "-", 2); err != nil {
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
			return nil, fmt.Errorf("cannot subtract with %#v", arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func powFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "pow", 2); err != nil {
		return nil, err
	}

	var x, y float64

	if v, ok := args[0].(int); ok {
		x = float64(v)
	} else if f, ok := args[0].(float64); ok {
		x = f
	} else {
		return nil, argExpectError("pow", "number", 1)
	}

	if v, ok := args[1].(int); ok {
		y = float64(v)
	} else if f, ok := args[1].(float64); ok {
		y = f
	} else {
		return nil, argExpectError("pow", "number", 2)
	}

	return math.Pow(x, y), nil
}

func randFn(args lisp.List) (value interface{}, err error) {
	if err := checkArityEqual(args, "rand", 0); err != nil {
		return nil, err
	}
	return rand.Float64(), nil
}

func randIntnFn(args lisp.List) (value interface{}, err error) {
	if len(args) > 1 {
		return nil, errors.New("rand-intn expects 0 or 1 arguments")
	}
	n, ok := args[0].(int)
	if !ok {
		return nil, argExpectError("rand-intn", "integer", 1)
	}
	return rand.Intn(n), nil
}
