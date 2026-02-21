package builtin

import (
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func errorfFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 1); err != nil {
		return nil, err
	}
	format, ok := args[0].(string)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeString, 1)
	}
	return errors.Errorf(format, args[1:]...), nil
}

func errorFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeString, 1)
	}
	return errors.New(str), nil
}
