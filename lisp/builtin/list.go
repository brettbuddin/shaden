package builtin

import (
	"buddin.us/shaden/lisp"
)

func listFn(args lisp.List) (interface{}, error) {
	return args, nil
}

func consFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "cons", 2); err != nil {
		return nil, err
	}
	list := lisp.List{args[0]}
	switch v := args[1].(type) {
	case lisp.List:
		for _, e := range v {
			list = append(list, e)
		}
	case lisp.Table:
		for k, e := range v {
			list = append(list, lisp.List{k, e})
		}
	default:
		list = append(list, v)
	}
	return list, nil
}

func firstFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "first", 1); err != nil {
		return nil, err
	}
	if args[0] == nil {
		return lisp.List{}, nil
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError("first", "list", 1)
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func restFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "rest", 1); err != nil {
		return nil, err
	}
	if args[0] == nil {
		return lisp.List{}, nil
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError("rest", "list", 1)
	}
	if len(list) == 0 {
		return lisp.List{}, nil
	}
	return list[1:], nil
}

func appendFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "append", 2); err != nil {
		return nil, err
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError("append", "list", 1)
	}
	return append(list, args[1:]...), nil
}

func prependFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "prepend", 1); err != nil {
		return nil, err
	}
	list, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError("prepend", "list", 1)
	}
	return append(args[1:], list...), nil
}

func lenFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "len", 1); err != nil {
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
		return nil, argExpectError("len", "list, table or string", 1)
	}
}
