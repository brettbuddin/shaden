package builtin

import (
	"errors"
	"fmt"

	"buddin.us/shaden/lisp"
)

func doFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	env = env.Branch()
	var (
		value interface{}
		err   error
	)
	for _, arg := range args {
		value, err = env.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func letFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "let", 2); err != nil {
		return nil, err
	}

	bindings, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError("let", "list", 1)
	}

	env = env.Branch()
	for _, n := range bindings {
		if list, ok := n.(lisp.List); ok {
			if len(list) != 2 {
				return nil, errors.New("let expects bindings to be list pairs")
			}
			name, ok := list[0].(lisp.Symbol)
			if !ok {
				return nil, errors.New("let expects binding names to be symbols")
			}
			value, err := env.Eval(list[1])
			if err != nil {
				return nil, err
			}
			env.DefineSymbol(string(name), value)
		}
	}
	var (
		value interface{}
		err   error
	)
	for _, arg := range args[1:] {
		value, err = env.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func fnFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "fn", 2); err != nil {
		return nil, err
	}
	params, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError("fn", "list", 1)
	}
	for _, n := range params {
		if _, ok := n.(lisp.Symbol); !ok {
			return nil, errors.New("fn expects all function parameters to be symbols")
		}
	}
	return buildFunction(env, "anonymous function", params, args[1:]), nil
}

func buildFunction(env *lisp.Environment, name string, defArgs, body lisp.List) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		env = env.Branch()
		return functionEvaluate(env, name, args, defArgs, body)
	}
}

func buildMacroFunction(env *lisp.Environment, name string, defArgs, body lisp.List) func(*lisp.Environment, lisp.List) (interface{}, error) {
	return func(env *lisp.Environment, args lisp.List) (interface{}, error) {
		env = env.Branch()
		v, err := functionEvaluate(env, name, args, defArgs, body)
		if err != nil {
			return nil, err
		}
		return env.Eval(v)
	}
}

func functionEvaluate(env *lisp.Environment, name string, args, defArgs, body lisp.List) (interface{}, error) {
	if len(args) != len(defArgs) {
		switch len(defArgs) {
		case 0:
			return nil, fmt.Errorf("%s expects 0 arguments", name)
		case 1:
			return nil, fmt.Errorf("%s expects 1 argument", name)
		default:
			return nil, fmt.Errorf("%s expects %d arguments", name, len(defArgs))
		}
	}
	for i, arg := range args {
		name := defArgs[i].(lisp.Symbol)
		env.DefineSymbol(string(name), arg)
	}
	var (
		value interface{}
		err   error
	)
	for _, arg := range body {
		value, err = env.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func applyFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "apply", 2); err != nil {
		return nil, err
	}

	// Ensure that (apply f 1 2 3 (list 4 5 6)) == (apply f 1 2 3 4 5 6)
	flat := lisp.List{}
	for _, arg := range args[1:] {
		if list, ok := arg.(lisp.List); ok {
			for _, e := range list {
				flat = append(flat, e)
			}
		} else {
			flat = append(flat, arg)
		}
	}

	switch fn := args[0].(type) {
	case lisp.Func:
		return fn.Func(flat)
	case func(lisp.List) (interface{}, error):
		return fn(flat)
	default:
		return nil, argExpectError("apply", "function", 1)
	}
}
