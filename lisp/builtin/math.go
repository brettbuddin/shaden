package builtin

import (
	"math"
	"math/rand"

	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

// foldNumeric applies intOp/floatOp across a list of numeric arguments.
// It preserves the int type when no float64 values are present.
// For commutative ops (sum, mult), set firstIsIdentity=true so the identity
// value seeds both accumulators. For non-commutative ops (diff, div), set
// firstIsIdentity=false so the first arg seeds the accumulators instead.
func foldNumeric(
	args lisp.List,
	intOp func(a, b int) int,
	floatOp func(a, b float64) float64,
	identity int,
	firstIsIdentity bool,
	verb string,
) (any, error) {
	var (
		outi      = identity
		outf      = float64(identity)
		seenFloat bool
	)
	for i, arg := range args {
		switch v := arg.(type) {
		case int:
			if i == 0 && !firstIsIdentity {
				outi = v
				outf = float64(v)
			} else {
				outi = intOp(outi, v)
				outf = floatOp(outf, float64(v))
			}
		case float64:
			if i == 0 && !firstIsIdentity {
				outf = v
			} else {
				outf = floatOp(outf, v)
			}
			seenFloat = true
		default:
			return nil, errors.Errorf("cannot %s %#v", verb, arg)
		}
	}
	if seenFloat {
		return outf, nil
	}
	return outi, nil
}

func sumFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	return foldNumeric(args,
		func(a, b int) int { return a + b },
		func(a, b float64) float64 { return a + b },
		0, true, "add")
}

func diffFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	return foldNumeric(args,
		func(a, b int) int { return a - b },
		func(a, b float64) float64 { return a - b },
		0, false, "subtract")
}

func multFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	return foldNumeric(args,
		func(a, b int) int { return a * b },
		func(a, b float64) float64 { return a * b },
		1, true, "multiply")
}

func divFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	return foldNumeric(args,
		func(a, b int) int { return a / b },
		func(a, b float64) float64 { return a / b },
		0, false, "divide")
}

func powFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	x, err := lisp.ExtractFloat64(args[0], 1)
	if err != nil {
		return nil, err
	}
	y, err := lisp.ExtractFloat64(args[1], 2)
	if err != nil {
		return nil, err
	}
	return math.Pow(x, y), nil
}

func randFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 0); err != nil {
		return nil, err
	}
	return rand.Float64(), nil
}
