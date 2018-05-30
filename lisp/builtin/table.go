package builtin

import (
	"errors"
	"fmt"

	"github.com/brettbuddin/shaden/lisp"
)

func tableFn(args lisp.List) (interface{}, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("table expects an even number of arguments")
	}
	m := lisp.Table{}
	for i := 0; i < len(args); i += 2 {
		m[args[i]] = args[i+1]
	}
	return m, nil
}

func tgetFn(args lisp.List) (interface{}, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("table-get expects 2 or 3 arguments")
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, argExpectError("table-get", "table", 1)
	}
	if v, ok := m[args[1]]; ok {
		return v, nil
	}
	if len(args) == 3 {
		return args[2], nil
	}
	return nil, nil
}

func tsetFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "table-set", 3); err != nil {
		return nil, err
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, argExpectError("table-set", "table", 1)
	}
	m[args[1]] = args[2]
	return nil, nil
}

func tdeleteFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "table-delete", 2); err != nil {
		return nil, err
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, argExpectError("table-delete", "table", 1)
	}
	delete(m, args[1])
	return nil, nil
}

func texistsFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "table-exists?", 2); err != nil {
		return nil, err
	}
	m, ok := args[0].(lisp.Table)
	if !ok {
		return nil, argExpectError("table-exists?", "table", 1)
	}
	_, ok = m[args[1]]
	return ok, nil
}

func mergeFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "table-merge", 2); err != nil {
		return nil, err
	}
	m := lisp.Table{}
	for _, arg := range args {
		if argm, ok := arg.(lisp.Table); ok {
			for k, v := range argm {
				m[k] = v
			}
		} else {
			return nil, fmt.Errorf("table-merge expects arguments to be tables")
		}
	}
	return m, nil
}

func tselectFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "table-select", 2); err != nil {
		return nil, err
	}
	filtered := lisp.Table{}

	t, ok := args[0].(lisp.Table)
	if !ok {
		return nil, argExpectError("table-select", "table", 1)
	}

	fn, ok := args[1].(func(lisp.List) (interface{}, error))
	if !ok {
		return nil, argExpectError("table-select", "function", 2)
	}

	for k, v := range t {
		result, err := fn(lisp.List{k, v})
		if err != nil {
			return nil, err
		}
		b, ok := result.(bool)
		if !ok {
			return nil, errors.New("table-select expects function to return boolean value")
		}
		if b {
			filtered[k] = v
		}
	}

	return filtered, nil
}
