package lisp

import (
	"fmt"
	"sync"

	"github.com/brettbuddin/shaden/errors"
)

const unquote = Symbol("unquote")

// Environment is a registry of symbols for a lexical scope.
type Environment struct {
	sync.RWMutex
	parent  *Environment
	symbols map[string]any
}

// NewEnvironment returns a new Environment.
func NewEnvironment() *Environment {
	return &Environment{
		symbols: map[string]any{},
	}
}

// Branch returns a new child context Environment.
func (e *Environment) Branch() *Environment {
	e.RLock()
	defer e.RUnlock()
	return &Environment{
		parent:  e,
		symbols: map[string]any{},
	}
}

// DefineSymbol defines a symbol.
func (e *Environment) DefineSymbol(symbol string, v any) error {
	e.RLock()
	exist, ok := e.symbols[symbol]
	e.RUnlock()

	if ok {
		if err := replace(exist, v); err != nil {
			return err
		}
	}

	e.Lock()
	e.symbols[symbol] = v
	e.Unlock()
	return nil
}

// SetSymbol sets the value of a symbol. Like GetSymbol, it advances to parent Environments if the symbol cannot be
// found in the current context. If the value has not been defined yet, it errors.
func (e *Environment) SetSymbol(symbol string, v any) error {
	env := e
	for env != nil {
		env.RLock()
		exist, ok := env.symbols[symbol]
		env.RUnlock()

		if ok {
			if err := replace(exist, v); err != nil {
				return err
			}

			env.Lock()
			env.symbols[symbol] = v
			env.Unlock()
			return nil
		}

		env = env.parent
	}
	return UndefinedSymbolError{symbol}
}

func replace(existing, replacement any) error {
	type replacer interface {
		Replace(v any) error
	}

	switch v := existing.(type) {
	case replacer:
		return v.Replace(replacement)
	case Table:
		el, ok := v[Keyword("__replace")]
		if !ok {
			return nil
		}
		fn, ok := el.(func(List) (any, error))
		if !ok {
			return errors.New("table key :__replace should be a function")
		}
		_, err := fn(List{replacement})
		return err
	default:
		return nil
	}
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
func (e *Environment) GetSymbol(symbol string) (any, error) {
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

// Eval evaluates expressions obtained via the Parser. It uses a trampoline
// loop to iteratively resolve TailCall values, enabling tail call optimization
// without growing the Go call stack.
func (e *Environment) Eval(node any) (any, error) {
	var errWrap string
	for {
		result, err := e.eval(node)
		if err != nil {
			if errWrap != "" {
				err = errors.Wrap(err, errWrap)
			}
			return nil, err
		}
		tc, ok := result.(TailCall)
		if !ok {
			return result, nil
		}
		node = tc.Node
		e = tc.Env
		if errWrap == "" && tc.ErrWrap != "" {
			errWrap = tc.ErrWrap
		}
	}
}

// eval is the internal evaluator. It may return TailCall sentinel values that
// the trampoline in Eval resolves iteratively.
func (e *Environment) eval(node any) (any, error) {
	switch node := node.(type) {
	case *root:
		var (
			value any
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
		return nil, errors.Errorf("unknown node type %T", node)
	}
}

// QuasiQuoteEval is a special evaluation mode that is used in quasiquoting. Like `quote` it doesn't evaluate any
// expressions. However, unlike `quote` it can evaluation some nested expressions marked by the special `unquote`
// function invocation.
func (e *Environment) QuasiQuoteEval(node any) (any, error) {
	switch node := node.(type) {
	case List:
		var result List
		for _, n := range node {
			var (
				v        any
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
		return nil, errors.Errorf("unknown node type %T", node)
	}
}

func (e *Environment) call(nodes List) (any, error) {
	if len(nodes) == 0 {
		return List{}, nil
	}

	head, rest := nodes[0], nodes[1:]
	fn, err := e.Eval(head)
	if err != nil {
		return nil, err
	}

	switch fn := fn.(type) {
	case Func:
		name := fn.Name()
		vargs, err := e.evalArgs(rest)
		if err != nil {
			return nil, errors.Wrap(err, evalArgsErrorMsg(name))
		}
		return e.callFunc(name, fn.Func, vargs)
	case func(List) (any, error):
		name := fnName(head)
		vargs, err := e.evalArgs(rest)
		if err != nil {
			return nil, errors.Wrap(err, evalArgsErrorMsg(name))
		}
		return e.callFunc(name, fn, vargs)
	case EnvFunc:
		return e.callEnvFunc(fn.Name(), fn.EnvFunc, rest)
	case func(*Environment, List) (any, error):
		name := fnName(head)
		return e.callEnvFunc(name, fn, rest)
	}
	return nil, errors.Errorf("uncallable function %#v", fn)
}

func fnName(v any) string {
	if sym, ok := v.(Symbol); ok {
		return string(sym)
	}
	return "anonymous function"
}

func (e *Environment) callFunc(name string, fn func(List) (any, error), args List) (any, error) {
	result, err := fn(args)
	if err != nil {
		return result, errors.Wrapf(err, "failed to call %s", name)
	}
	if tc, ok := result.(TailCall); ok {
		if tc.ErrWrap == "" {
			tc.ErrWrap = fmt.Sprintf("failed to call %s", name)
		}
		return tc, nil
	}
	return result, nil
}

func (e *Environment) callEnvFunc(name string, fn func(*Environment, List) (any, error), args List) (any, error) {
	result, err := fn(e, args)
	if err != nil {
		return result, errors.Wrapf(err, "failed to call %s", name)
	}
	if tc, ok := result.(TailCall); ok {
		if tc.ErrWrap == "" {
			tc.ErrWrap = fmt.Sprintf("failed to call %s", name)
		}
		return tc, nil
	}
	return result, nil
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
