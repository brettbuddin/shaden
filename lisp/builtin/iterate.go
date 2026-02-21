package builtin

import (
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func mapFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	fn, err := lisp.ExtractArg[func(lisp.List) (any, error)](args, 0, lisp.TypeFunction)
	if err != nil {
		return nil, err
	}
	var out lisp.List
	if err := lisp.ForEachEntry(args[1], 2, func(k, v any) error {
		r, err := fn(lisp.List{k, v})
		if err != nil {
			return err
		}
		out = append(out, r)
		return nil
	}); err != nil {
		return nil, err
	}
	return out, nil
}

func eachFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	fn, err := lisp.ExtractArg[func(lisp.List) (any, error)](args, 0, lisp.TypeFunction)
	if err != nil {
		return nil, err
	}
	if err := lisp.ForEachEntry(args[1], 2, func(k, v any) error {
		_, err := fn(lisp.List{k, v})
		return err
	}); err != nil {
		return nil, err
	}
	return args[1], nil
}

func dotimesFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	binding, err := lisp.ExtractArg[lisp.List](args, 0, lisp.TypeList)
	if err != nil {
		return nil, err
	}
	if len(binding) != 2 {
		return nil, lisp.ArgExpectError("name/value pair binding", 1)
	}

	body, err := lisp.ExtractArg[lisp.List](args, 1, lisp.TypeList)
	if err != nil {
		return nil, err
	}

	env = env.Branch()
	name, err := lisp.Extract[lisp.Symbol](binding[0], 1, lisp.TypeSymbol)
	if err != nil {
		return nil, errors.Errorf("expects binding name to be a symbol")
	}
	value, err := env.Eval(binding[1])
	if err != nil {
		return nil, err
	}
	n, ok := value.(int)
	if !ok {
		return nil, errors.Errorf("expects an int for binding value")
	}

	env.DefineSymbol(string(name), 0)

	for i := 0; i < n; i++ {
		if i > 0 {
			if err := env.SetSymbol(string(name), i); err != nil {
				return nil, err
			}
		}
		if _, err = env.Eval(body); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func reduceFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 3); err != nil {
		return nil, err
	}
	fn, err := lisp.ExtractArg[func(lisp.List) (any, error)](args, 0, lisp.TypeFunction)
	if err != nil {
		return nil, err
	}
	acc := args[1]
	if err := lisp.ForEachEntry(args[2], 3, func(k, v any) error {
		var err error
		acc, err = fn(lisp.List{acc, k, v})
		return err
	}); err != nil {
		return nil, err
	}
	return acc, nil
}
