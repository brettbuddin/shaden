package builtin

import (
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func mapFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, 2); err != nil {
		return nil, err
	}
	fn, ok := args[0].(func(lisp.List) (interface{}, error))
	if !ok {
		return nil, argExpectError(typeFunction, 1)
	}

	switch v := args[1].(type) {
	case lisp.List:
		var out lisp.List
		for i, e := range v {
			r, err := fn(lisp.List{i, e})
			if err != nil {
				return nil, err
			}
			out = append(out, r)
		}
		return out, nil
	case lisp.Table:
		var out lisp.List
		for k, e := range v {
			r, err := fn(lisp.List{k, e})
			if err != nil {
				return nil, err
			}
			out = append(out, r)
		}
		return out, nil
	default:
		return nil, argExpectError(acceptTypes(typeList, typeTable), 2)
	}
}

func eachFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, 2); err != nil {
		return nil, err
	}
	fn, ok := args[0].(func(lisp.List) (interface{}, error))
	if !ok {
		return nil, argExpectError(typeFunction, 1)
	}

	switch v := args[1].(type) {
	case lisp.List:
		for i, e := range v {
			_, err := fn(lisp.List{i, e})
			if err != nil {
				return nil, err
			}
		}
		return v, nil
	case lisp.Table:
		for k, e := range v {
			_, err := fn(lisp.List{k, e})
			if err != nil {
				return nil, err
			}
		}
		return v, nil
	default:
		return nil, argExpectError(acceptTypes(typeList, typeTable), 2)
	}
}

func dotimesFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, 2); err != nil {
		return nil, err
	}
	binding, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError(typeList, 1)
	}
	if len(binding) != 2 {
		return nil, argExpectError("name/value pair binding", 1)
	}

	body, ok := args[1].(lisp.List)
	if !ok {
		return nil, argExpectError(typeList, 2)
	}

	env = env.Branch()
	name, ok := binding[0].(lisp.Symbol)
	if !ok {
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

	if err := env.DefineSymbol(string(name), 0); err != nil {
		return nil, err
	}
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

func reduceFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, 3); err != nil {
		return nil, err
	}
	fn, ok := args[0].(func(lisp.List) (interface{}, error))
	if !ok {
		return nil, argExpectError(typeFunction, 1)
	}

	switch v := args[2].(type) {
	case lisp.List:
		reduce := args[1]
		var err error
		for i, e := range v {
			reduce, err = fn(lisp.List{reduce, i, e})
			if err != nil {
				return nil, err
			}
		}
		return reduce, nil
	case lisp.Table:
		reduce := args[1]
		var err error
		for k, e := range v {
			reduce, err = fn(lisp.List{reduce, k, e})
			if err != nil {
				return nil, err
			}
		}
		return reduce, nil
	default:
		return nil, argExpectError(acceptTypes(typeList, typeTable), 3)
	}
}
