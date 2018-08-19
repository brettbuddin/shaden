package lisp

import (
	"fmt"
	"strings"

	"github.com/brettbuddin/shaden/errors"
)

// UndefinedSymbolError is an error returned when a symbol name cannot be found in an Environment.
type UndefinedSymbolError struct {
	Name string
}

func (e UndefinedSymbolError) Error() string {
	return fmt.Sprintf("undefined symbol %s", e.Name)
}

type lineError struct {
	error error
	line  pos
}

func (e lineError) Error() string {
	return fmt.Sprintf("%s (line %d)", e.error, e.line)
}

func (e lineError) GoString() string {
	return e.Error()
}

func newLineError(err error, line pos) error {
	if _, ok := err.(lineError); ok {
		return err
	}
	return lineError{
		error: err,
		line:  line,
	}
}

// CheckArityEqual requires the argument list to be of a certain length
func CheckArityEqual(l List, expected int) error {
	actual := len(l)
	if actual != expected {
		return errors.Errorf("expects %d argument; %d given", expected, actual)
	}
	return nil
}

// CheckArityEqual requires the argument list to be at least certain length
func CheckArityAtLeast(l List, expected int) error {
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

// ArgExpectError returns an error that indicates requirements for an argument
func ArgExpectError(what string, pos int) error {
	return errors.Errorf("expects %s for argument %d", what, pos)
}

// AcceptTypes creates a formatted list of accepted types for humans
func AcceptTypes(names ...string) string {
	if len(names) == 2 {
		return strings.Join(names, " or ")
	}
	return strings.Replace(strings.Join(names, ", "), ",", "or", len(names)-1)
}
