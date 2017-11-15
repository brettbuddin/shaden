package lisp

import "fmt"

// UndefinedSymbolError is an error returned when a symbol name cannot be found in an Environment.
type UndefinedSymbolError struct {
	Name string
}

func (e UndefinedSymbolError) Error() string {
	return fmt.Sprintf("undefined symbol %s", e.Name)
}

// DefinedSymbolError is an error returned when a symbol has already been defined in the current Environment.
type DefinedSymbolError struct {
	Name string
}

func (e DefinedSymbolError) Error() string {
	return fmt.Sprintf("symbol %s already defined", e.Name)
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
