package builtin

import (
	"github.com/brettbuddin/shaden/lisp"
)

func listFn(args lisp.List) (any, error) {
	return args, nil
}

func reverseFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch target := args[0].(type) {
	case lisp.List:
		var list lisp.List
		for i := len(target) - 1; i >= 0; i-- {
			list = append(list, target[i])
		}
		return list, nil
	case string:
		r := []rune(target)
		for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r), nil
	default:
		return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeList, lisp.TypeString), 1)
	}
}

func consFn(args lisp.List) (any, error) {
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

func firstFn(args lisp.List) (any, error) {
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

func restFn(args lisp.List) (any, error) {
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

func appendFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
	return append(list, args[1:]...), nil
}

func prependFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 1); err != nil {
		return nil, err
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
	return append(args[1:], list...), nil
}

func lenFn(args lisp.List) (any, error) {
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
