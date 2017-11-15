package builtin

import (
	"github.com/pkg/errors"

	"buddin.us/shaden/lisp"
)

const (
	underscoreSymbol = lisp.Symbol("_")
	ampersandSymbol  = lisp.Symbol("&")
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
			if err := env.DefineSymbol(string(name), value); err != nil {
				return nil, err
			}
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

func functionArityError(name string, defArgCount int) error {
	switch defArgCount {
	case 0:
		return errors.Errorf("%s expects 0 arguments", name)
	case 1:
		return errors.Errorf("%s expects 1 argument", name)
	default:
		return errors.Errorf("%s expects %d arguments", name, defArgCount)
	}
}

func functionEvaluate(env *lisp.Environment, name string, args, defArgs, body lisp.List) (interface{}, error) {
	// Locate the variadic symbol "&" position
	var (
		variadicAt     = -1
		variadicSymbol lisp.Symbol
		variadicArgs   = lisp.List{}
	)
	for i, arg := range defArgs {
		if arg.(lisp.Symbol) == ampersandSymbol {
			variadicAt = i
			break
		}
	}

	// TODO: Ensure the variadic symbol is right next to the last varibale. If there are more than 1 more symbol, we
	// should error.

	if variadicAt < 0 {
		if len(args) != len(defArgs) {
			return nil, functionArityError(name, len(defArgs))
		}
	} else if len(args) < variadicAt {
		return nil, functionArityError(name, variadicAt)
	} else if variadicAt >= 0 {
		if len(defArgs)-2 != variadicAt {
			return nil, errors.New("definition has too many arguments after variadic symbol &")
		}
		variadicSymbol = defArgs[variadicAt+1].(lisp.Symbol)
	}

	for i, arg := range args {
		if variadicAt >= 0 && i >= variadicAt {
			variadicArgs = append(variadicArgs, arg)
		} else {
			symbol := defArgs[i].(lisp.Symbol)
			if symbol == underscoreSymbol {
				continue
			}
			if err := env.DefineSymbol(string(symbol), arg); err != nil {
				return nil, err
			}
		}
	}

	if variadicSymbol != "" && variadicSymbol != underscoreSymbol {
		if err := env.DefineSymbol(string(variadicSymbol), variadicArgs); err != nil {
			return nil, err
		}
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
