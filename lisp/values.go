package lisp

import "fmt"

// Func is a object that can act as a function invocation in the lisp. It receives fully evaluated arguments.
type Func interface {
	Name() string
	Func(List) (interface{}, error)
}

// EnvFunc is a object that can act as a function invocation in the lisp. It recieves the current Environment and
// unevaluated arguments. This leaves it up to the function to describe how arguments should be evaluated.
type EnvFunc interface {
	Name() string
	EnvFunc(*Environment, List) (interface{}, error)
}

// Symbol represents a symbol lisp type
type Symbol string

// Keyword represents a keyword lisp type
type Keyword string

func (k Keyword) Name() string { return fmt.Sprintf("keyword function %s", k) }

// Func implements the Func interface. It allows Keywords to be called like functions that accept a Table as their first
// argument to return the associated value for that key in the Table. An optional third argument can be provided as a
// default value to be returned if there is no value for the key.
func (k Keyword) Func(args List) (interface{}, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("keyword function %s expects 1 or 2 arguments", k)
	}
	m, ok := args[0].(Table)
	if !ok {
		return nil, fmt.Errorf("keyword function %s expects hash for argument 1", k)
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
type List []interface{}

func (List) Name() string { return "list function" }

// Func implements the Func interface. It allows Lists to be called like functions that accept an integer as their first
// argument to perform offset indexing.
func (l List) Func(args List) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("list expects 1 argument")
	}
	idx, ok := args[0].(int)
	if !ok {
		return nil, fmt.Errorf("list expects integer for argument 1")
	}
	if idx > len(l)-1 {
		return nil, fmt.Errorf("index out of range")
	}
	return l[idx], nil
}

// Table represents a table lisp type
type Table map[interface{}]interface{}

func (Table) Name() string { return "table function" }

// Func implements the Func interface. It allows Tables to be called like functions that accept a key as their first
// argument to return the associated value in the Table. An optional third argument can be provided as a default value
// to be returned if there is no value for the key.
func (t Table) Func(args List) (interface{}, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("table expects 1 or 2 arguments")
	}
	if v, ok := t[args[0]]; ok {
		return v, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return nil, nil
}
