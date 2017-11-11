package lisp

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

const unquote = Symbol("unquote")

// Environment is a registry of symbols for a lexical scope.
type Environment struct {
	sync.RWMutex
	parent  *Environment
	symbols map[string]interface{}
}

// NewEnvironment returns a new Environment.
func NewEnvironment() *Environment {
	return &Environment{
		symbols: map[string]interface{}{},
	}
}

// Branch returns a new child context Environment.
func (e *Environment) Branch() *Environment {
	e.RLock()
	defer e.RUnlock()
	return &Environment{
		parent:  e,
		symbols: map[string]interface{}{},
	}
}

// DefineSymbol defines or sets the value of a symbol.
func (e *Environment) DefineSymbol(symbol string, v interface{}) {
	e.Lock()
	defer e.Unlock()
	e.symbols[symbol] = v
}

// SetSymbol sets the value of a symbol. Like GetSymbol, it advances to parent Environments if the symbol cannot be
// found in the current context. If the value has not been defined yet, it errors.
func (e *Environment) SetSymbol(symbol string, v interface{}) error {
	e.Lock()
	defer e.Unlock()
	for e != nil {
		if _, ok := e.symbols[symbol]; ok {
			e.symbols[symbol] = v
			return nil
		}
		e = e.parent
	}
	return UndefinedSymbolError{symbol}
}

// UnsetSymbol removes a symbol definition. It only operates on the current context; no parent Environments will be
// affected.
func (e *Environment) UnsetSymbol(symbol string) error {
	e.Lock()
	defer e.Unlock()
	_, ok := e.symbols[symbol]
	if !ok {
		return UndefinedSymbolError{symbol}
	}
	delete(e.symbols, symbol)
	return nil
}

// GetSymbol performs a symbol lookup and returns the value if its present. If the symbol cannot be found in the current
// Environment context, it advances to the parent to see if can be found there.
func (e *Environment) GetSymbol(symbol string) (interface{}, error) {
	e.RLock()
	defer e.RUnlock()
	for e != nil {
		if v, ok := e.symbols[symbol]; ok {
			return v, nil
		}
		e = e.parent
	}
	return nil, UndefinedSymbolError{symbol}
}

// Eval evaluates expressions obtained via the Parser.
func (e *Environment) Eval(node interface{}) (interface{}, error) {
	switch node := node.(type) {
	case *root:
		var (
			value interface{}
			err   error
		)
		for _, node := range node.Nodes {
			value, err = e.Eval(node)
			if err != nil {
				return nil, err
			}
		}
		return value, nil
	case List:
		return e.call(node)
	case string, float64, int, bool, Keyword:
		return node, nil
	case Symbol:
		return e.GetSymbol(string(node))
	default:
		if node == nil {
			return nil, nil
		}
		return nil, newLineError(errors.Errorf("unknown node type %T", node), 0)
	}
}

// QuasiQuoteEval is a special evaluation mode that is used in quasiquoting. Like `quote` it doesn't evaluate any
// expressions. However, unlike `quote` it can evaluation some nested expressions marked by the special `unquote`
// function invocation.
func (e *Environment) QuasiQuoteEval(node interface{}) (interface{}, error) {
	switch node := node.(type) {
	case List:
		var result List
		for _, n := range node {
			var (
				v        interface{}
				err      error
				list, ok = n.(List)
			)

			if ok {
				switch size := len(list); {
				case size == 2 && list[0] == unquote:
					v, err = e.Eval(list[1])
				case size > 0 && list[0] == unquote:
					err = errors.Errorf("unquote expects 1 argument")
				default:
					v, err = e.QuasiQuoteEval(n)
				}
			} else {
				v, err = e.QuasiQuoteEval(n)
			}
			if err != nil {
				return nil, err
			}
			result = append(result, v)
		}
		return result, nil
	case string, float64, int, bool, Keyword, Symbol:
		return node, nil
	default:
		if node == nil {
			return nil, nil
		}
		return nil, newLineError(errors.Errorf("unknown node type %T", node), 0)
	}
}

func (e *Environment) call(nodes List) (interface{}, error) {
	if len(nodes) == 0 {
		return List{}, nil
	}

	head, rest := nodes[0], nodes[1:]
	fn, err := e.Eval(head)
	if err != nil {
		return nil, err
	}

	name := fnName(head)

	switch fn := fn.(type) {
	case Func:
		vargs, err := e.evalArgs(rest)
		if err != nil {
			return nil, errors.Wrap(err, evalArgsErrorMsg(name))
		}
		return e.callFunc(name, fn.Func, vargs)
	case func(List) (interface{}, error):
		vargs, err := e.evalArgs(rest)
		if err != nil {
			return nil, errors.Wrap(err, evalArgsErrorMsg(name))
		}
		return e.callFunc(name, fn, vargs)
	case EnvFunc:
		return e.callEnvFunc(name, fn.EnvFunc, rest)
	case func(*Environment, List) (interface{}, error):
		return e.callEnvFunc(name, fn, rest)
	}
	return nil, errors.Errorf("uncallable function %#v", fn)
}

func (e *Environment) callFunc(name string, fn func(List) (interface{}, error), args List) (interface{}, error) {
	result, err := fn(args)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("error calling %q", name))
	}
	return result, nil
}

func (e *Environment) callEnvFunc(name string, fn func(*Environment, List) (interface{}, error), args List) (interface{}, error) {
	result, err := fn(e, args)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("error calling %q", name))
	}
	return result, nil
}

func fnName(v interface{}) string {
	switch v := v.(type) {
	case Symbol:
		return string(v) + " function"
	case Keyword:
		return string(v) + " function"
	case List:
		return "list function"
	case Table:
		return "table function"
	default:
		return fmt.Sprintf("uncallable value %T", v)
	}
}

func (e *Environment) evalArgs(nodes List) (List, error) {
	values := make(List, len(nodes))
	for i, n := range nodes {
		value, err := e.Eval(n)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func evalArgsErrorMsg(name string) string {
	return fmt.Sprintf("error evaluating arguments of %s", name)
}
