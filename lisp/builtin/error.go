package builtin

import (
	"fmt"

	"github.com/pkg/errors"

	"buddin.us/shaden/lisp"
)

func errorfFn(args lisp.List) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("errorf expects at least 1 argument")
	}
	format, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("errorf expects string for argument 1")
	}
	return fmt.Errorf(format, args[1:]...), nil
}

// func attemptAllFn(env *lisp.Environment, args []lisp.Node) (interface{}, error) {
// 	return nil, nil
// }

// func ifFailureFn(env *lisp.Environment, args []lisp.Node) (interface{}, error) {
// 	return nil, nil
// }

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

func checkArityRange(l lisp.List, name string, min, max int) error {
	actual := len(l)
	if actual < min || max > actual {
		return errors.Errorf("%s expects between %d and %d arguments; %d given", name, min, max, actual)
	}
	return nil
}

func argExpectError(name, what string, pos int) error {
	return errors.Errorf("%s expects %s for argument %d", name, what, pos)
}
