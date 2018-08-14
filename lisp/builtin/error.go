package builtin

import (
	"strings"

	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func errorfFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, 1); err != nil {
		return nil, err
	}
	format, ok := args[0].(string)
	if !ok {
		return nil, argExpectError(typeString, 1)
	}
	return errors.Errorf(format, args[1:]...), nil
}

func errorFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, 1); err != nil {
		return nil, err
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, argExpectError(typeString, 1)
	}
	return errors.New(str), nil
}

func checkArityEqual(l lisp.List, expected int) error {
	actual := len(l)
	if actual != expected {
		return errors.Errorf("expects %d argument; %d given", expected, actual)
	}
	return nil
}

func checkArityAtLeast(l lisp.List, expected int) error {
	actual := len(l)
	if actual < expected {
		var plural string
		if expected > 1 {
			plural = "s"
		}
		return errors.Errorf("expects at least %d argument%s; %d given", expected, plural, actual)
	}
	return nil
}

func argExpectError(what string, pos int) error {
	return errors.Errorf("expects %s for argument %d", what, pos)
}

func acceptTypes(names ...string) string {
	return strings.Replace(strings.Join(names, ", "), ",", "or", len(names)-1)
}
