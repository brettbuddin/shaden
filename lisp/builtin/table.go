package builtin

import (
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func tableFn(args lisp.List) (any, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("expects an even number of arguments")
	}
	m := lisp.Table{}
	for i := 0; i < len(args); i += 2 {
		m[args[i]] = args[i+1]
	}
	return m, nil
}

func tgetFn(args lisp.List) (any, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.Errorf("expects 2 or 3 arguments")
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeTable, 1)
	}
	if v, ok := m[args[1]]; ok {
		return v, nil
	}
	if len(args) == 3 {
		return args[2], nil
	}
	return nil, nil
}

func tsetFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 3); err != nil {
		return nil, err
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeTable, 1)
	}
	m[args[1]] = args[2]
	return nil, nil
}

func tdeleteFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeTable, 1)
	}
	delete(m, args[1])
	return nil, nil
}

func texistsFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeTable, 1)
	}
	_, ok = m[args[1]]
	return ok, nil
}

func mergeFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	m := lisp.Table{}
	for _, arg := range args {
		if argm, ok := arg.(lisp.Table); ok {
			for k, v := range argm {
				m[k] = v
			}
		} else {
			return nil, errors.Errorf("expects arguments to be tables")
		}
	}
	return m, nil
}

func tselectFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	filtered := lisp.Table{}

	t, ok := args[0].(lisp.Table)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeTable, 1)
	}

	fn, ok := args[1].(func(lisp.List) (any, error))
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeFunction, 2)
	}

	for k, v := range t {
		result, err := lisp.ResolveTailCalls(fn(lisp.List{k, v}))
		if err != nil {
			return nil, err
		}
		b, ok := result.(bool)
		if !ok {
			return nil, errors.New("expects function to return boolean value")
		}
		if b {
			filtered[k] = v
		}
	}

	return filtered, nil
}
