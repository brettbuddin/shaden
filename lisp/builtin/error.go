package builtin

import (
	"buddin.us/shaden/errors"
	"buddin.us/shaden/lisp"
)

func errorfFn(args lisp.List) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.Errorf("errorf expects at least 1 argument")
	}
	format, ok := args[0].(string)
	if !ok {
		return nil, errors.Errorf("errorf expects string for argument 1")
	}
	return errors.Errorf(format, args[1:]...), nil
}

func checkArityEqual(l lisp.List, name string, expected int) error {
	actual := len(l)
	if actual != expected {
		return errors.Errorf("%s expects %d argument; %d given", name, expected, actual)
	}
	return nil
}

func checkArityAtLeast(l lisp.List, name string, expected int) error {
	actual := len(l)
	if actual < expected {
		var plural string
		if expected > 1 {
			plural = "s"
		}
		return errors.Errorf("%s expects at least %d argument%s; %d given", name, expected, plural, actual)
	}
	return nil
}

func argExpectError(name, what string, pos int) error {
	return errors.Errorf("%s expects %s for argument %d", name, what, pos)
}
