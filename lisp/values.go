package lisp

import (
	"fmt"

	"github.com/brettbuddin/shaden/errors"
)

// Lisp type names
const (
	TypeBool     = "bool"
	TypeFloat    = "float"
	TypeFunction = "function"
	TypeInt      = "int"
	TypeKeyword  = "keyword"
	TypeList     = "list"
	TypeString   = "string"
	TypeSymbol   = "symbol"
	TypeTable    = "table"
)

// Func is a object that can act as a function invocation in the lisp. It receives fully evaluated arguments.
type Func interface {
	Name() string
	Func(List) (any, error)
}

// EnvFunc is a object that can act as a function invocation in the lisp. It receives the current Environment and
// unevaluated arguments. This leaves it up to the function to describe how arguments should be evaluated.
type EnvFunc interface {
	Name() string
	EnvFunc(*Environment, List) (any, error)
}

// Symbol represents a symbol lisp type
type Symbol string

// Keyword represents a keyword lisp type
type Keyword string

// Name returns the name used when calling Keyword as a function.
func (k Keyword) Name() string { return fmt.Sprintf("keyword function %s", k) }

// Func implements the Func interface. It allows Keywords to be called like functions that accept a Table as their first
// argument to return the associated value for that key in the Table. An optional third argument can be provided as a
// default value to be returned if there is no value for the key.
func (k Keyword) Func(args List) (any, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.Errorf("expects 1 or 2 arguments")
	}
	m, ok := args[0].(Table)
	if !ok {
		return nil, ArgExpectError(TypeTable, 1)
	}
	if v, ok := m[k]; ok {
		return v, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return nil, nil
}

// List represents a list lisp type
type List []any

// Name returns the name used when calling List as a function.
func (List) Name() string { return "list function" }

// Func implements the Func interface. It allows Lists to be called like functions that accept an integer as their first
// argument to perform offset indexing.
func (l List) Func(args List) (any, error) {
	if err := CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	idx, ok := args[0].(int)
	if !ok {
		return nil, ArgExpectError(TypeInt, 1)
	}
	if idx > len(l)-1 {
		return nil, errors.Errorf("index out of range")
	}
	return l[idx], nil
}

// Table represents a table lisp type
type Table map[any]any

// Name returns the name used when calling Table as a function.
func (Table) Name() string { return "table function" }

// Func implements the Func interface. It allows Tables to be called like functions that accept a key as their first
// argument to return the associated value in the Table. An optional third argument can be provided as a default value
// to be returned if there is no value for the key.
func (t Table) Func(args List) (any, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.Errorf("expects 1 or 2 arguments")
	}
	if v, ok := t[args[0]]; ok {
		return v, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return nil, nil
}
