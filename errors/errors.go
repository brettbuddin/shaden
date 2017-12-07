package errors

import (
	"bytes"
	"errors"
	"fmt"
)

// Separator between errors and their causes. Overridable for situations where a newline and tab aren't desired.
var Separator = ":\n\t"

// New creates a new error
func New(msg string) error {
	return &Error{
		error: errors.New(msg),
	}
}

// Errorf creates a new error with a specific formatting
func Errorf(format string, args ...interface{}) error {
	return &Error{
		error: fmt.Errorf(format, args...),
	}
}

// Wrap creates a new error with another error as its cause
func Wrap(err error, msg string) error {
	return &Error{
		error: errors.New(msg),
		cause: err,
	}
}

// Wrapf creates a new error with a specific formatting with another error as its cause
func Wrapf(err error, format string, args ...interface{}) error {
	return &Error{
		error: fmt.Errorf(format, args...),
		cause: err,
	}
}

// Error is an error within shaden
type Error struct {
	error
	cause error
}

func (e *Error) isZero() bool {
	return e.error == nil && e.cause == nil
}

func (e *Error) Error() string {
	b := bytes.NewBuffer(nil)
	b.WriteString(e.error.Error())
	if e.cause != nil {
		writeString(b, Separator)
		b.WriteString(e.cause.Error())
	}
	return b.String()
}

func writeString(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
