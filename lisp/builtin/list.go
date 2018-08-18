package builtin

import (
	"github.com/brettbuddin/shaden/lisp"
)

func listFn(args lisp.List) (interface{}, error) {
	return args, nil
}

func consFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	list := lisp.List{args[0]}

	if args[1] == nil {
		return list, nil
	}

	switch v := args[1].(type) {
	case lisp.List:
		for _, e := range v {
			list = append(list, e)
		}
	default:
		list = append(list, v)
	}
	return list, nil
}

func firstFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	if args[0] == nil {
		return lisp.List{}, nil
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func restFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	if args[0] == nil {
		return lisp.List{}, nil
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
	if len(list) == 0 {
		return lisp.List{}, nil
	}
	return list[1:], nil
}

func appendFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
	return append(list, args[1:]...), nil
}

func prependFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 1); err != nil {
		return nil, err
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
	return append(args[1:], list...), nil
}

func lenFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	if args[0] == nil {
		return 0, nil
	}
	switch v := args[0].(type) {
	case string:
		return len(v), nil
	case lisp.List:
		return len(v), nil
	case lisp.Table:
		return len(v), nil
	default:
		return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeList, lisp.TypeTable, lisp.TypeString), 1)
	}
}
