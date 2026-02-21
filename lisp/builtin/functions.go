package builtin

import (
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

const (
	symbolUnderscore = lisp.Symbol("_")
	symbolAmpersand  = lisp.Symbol("&")
)

func beginFn(env *lisp.Environment, args lisp.List) (any, error) {
	env = env.Branch()
	if len(args) == 0 {
		return nil, nil
	}
	for _, arg := range args[:len(args)-1] {
		if _, err := env.Eval(arg); err != nil {
			return nil, err
		}
	}
	return lisp.TailCall{Node: args[len(args)-1], Env: env}, nil
}

func letFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}

	bindings, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}

	env = env.Branch()
	for _, n := range bindings {
		if list, ok := n.(lisp.List); ok {
			if len(list) != 2 {
				return nil, errors.Errorf("expects bindings to be list pairs")
			}
			name, ok := list[0].(lisp.Symbol)
			if !ok {
				return nil, errors.Errorf("expects binding names to be symbols")
			}
			value, err := env.Eval(list[1])
			if err != nil {
				return nil, err
			}
			env.DefineSymbol(string(name), value)
		}
	}
	body := args[1:]
	for _, arg := range body[:len(body)-1] {
		if _, err := env.Eval(arg); err != nil {
			return nil, err
		}
	}
	return lisp.TailCall{Node: body[len(body)-1], Env: env}, nil
}

func fnFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	params, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
	for _, n := range params {
		if _, ok := n.(lisp.Symbol); !ok {
			return nil, errors.Errorf("expects all function parameters to be symbols")
		}
	}
	return buildFunction(env, params, args[1:]), nil
}

func buildFunction(env *lisp.Environment, defArgs, body lisp.List) func(lisp.List) (any, error) {
	return func(args lisp.List) (any, error) {
		callEnv := env.Branch()
		return functionEvaluate(callEnv, args, defArgs, body)
	}
}

func buildMacroFunction(env *lisp.Environment, defArgs, body lisp.List) func(*lisp.Environment, lisp.List) (any, error) {
	return func(env *lisp.Environment, args lisp.List) (any, error) {
		env = env.Branch()
		v, err := functionEvaluate(env, args, defArgs, body)
		// functionEvaluate may return a TailCall; resolve it to get the actual
		// macro expansion form before evaluating it.
		v, err = lisp.ResolveTailCalls(v, err)
		if err != nil {
			return nil, err
		}
		return lisp.TailCall{Node: v, Env: env}, nil
	}
}

func functionArityError(defCount, givenCount int) error {
	switch defCount {
	case 0:
		return errors.Errorf("expects 0 arguments; %d given", givenCount)
	case 1:
		return errors.Errorf("expects 1 argument; %d given", givenCount)
	default:
		return errors.Errorf("expects %d arguments; %d given", defCount, givenCount)
	}
}

func functionEvaluate(env *lisp.Environment, args, defArgs, body lisp.List) (any, error) {
	// Locate the variadic symbol "&" position
	var (
		variadicAt     = -1
		variadicSymbol lisp.Symbol
		variadicArgs   = lisp.List{}
	)
	for i, arg := range defArgs {
		if arg.(lisp.Symbol) == symbolAmpersand {
			variadicAt = i
			break
		}
	}

	if variadicAt < 0 {
		if len(args) != len(defArgs) {
			return nil, functionArityError(len(defArgs), len(args))
		}
	} else if len(args) < variadicAt {
		return nil, functionArityError(variadicAt, len(args))
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
			if symbol == symbolUnderscore {
				continue
			}
			env.DefineSymbol(string(symbol), arg)
		}
	}

	if variadicSymbol != "" && variadicSymbol != symbolUnderscore {
		env.DefineSymbol(string(variadicSymbol), variadicArgs)
	}

	for _, arg := range body[:len(body)-1] {
		if _, err := env.Eval(arg); err != nil {
			return nil, err
		}
	}
	return lisp.TailCall{Node: body[len(body)-1], Env: env}, nil
}

func applyFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
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
		return lisp.ResolveTailCalls(fn.Func(flat))
	case func(lisp.List) (any, error):
		return lisp.ResolveTailCalls(fn(flat))
	default:
		return nil, lisp.ArgExpectError(lisp.TypeFunction, 1)
	}
}
