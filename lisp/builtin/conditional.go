package builtin

import (
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func equalFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	return args[0] == args[1], nil
}

func notEqualFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	return args[0] != args[1], nil
}

// compareOrdered compares two values of the same ordered type.
func compareOrdered[T ~int | ~float64](x, y T, less bool) bool {
	if less {
		return x < y
	}
	return x > y
}

// compareValues dispatches a comparison over the lisp numeric types.
func compareValues(args lisp.List, less bool) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	switch x := args[0].(type) {
	case int:
		if y, ok := args[1].(int); ok {
			return compareOrdered(x, y, less), nil
		}
		return nil, errors.Errorf("cannot compare %T and %T", args[0], args[1])
	case float64:
		if y, ok := args[1].(float64); ok {
			return compareOrdered(x, y, less), nil
		}
		return nil, errors.Errorf("cannot compare %T and %T", args[0], args[1])
	default:
		return nil, errors.Errorf("cannot compare %T and %T", args[0], args[1])
	}
}

func lessThanFn(args lisp.List) (any, error) {
	return compareValues(args, true)
}

func greaterThanFn(args lisp.List) (any, error) {
	return compareValues(args, false)
}

func andFn(env *lisp.Environment, args lisp.List) (any, error) {
	arity := len(args)
	if arity == 0 {
		return true, nil
	} else if arity == 1 {
		return args[0], nil
	}
	var (
		value any
		err   error
	)
	for _, arg := range args {
		value, err = env.Eval(arg)
		if err != nil {
			return nil, err
		}
		if value == nil || value == false {
			return false, nil
		}
	}
	return true, nil
}

func orFn(env *lisp.Environment, args lisp.List) (any, error) {
	arity := len(args)
	if arity == 0 {
		return false, nil
	} else if arity == 1 {
		return args[0], nil
	}
	var (
		value any
		err   error
	)
	for _, arg := range args {
		value, err = env.Eval(arg)
		if err != nil {
			return nil, err
		}
		if value != nil && value != false {
			return true, nil
		}
	}
	return false, nil
}

func notFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	var condition bool
	switch v := args[0].(type) {
	case bool:
		condition = !v
	default:
		if v == nil {
			return false, nil
		}
		return true, nil
	}
	return condition, nil
}

func ifFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 3); err != nil {
		return nil, err
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	if v == nil {
		return env.Eval(args[2])
	}
	condition, ok := v.(bool)
	if !ok {
		return env.Eval(args[1])
	}
	if condition {
		return env.Eval(args[1])
	}
	return env.Eval(args[2])
}

func whenFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	condition, ok := v.(bool)
	if ok && !condition {
		return nil, nil
	}
	var value any
	for _, arg := range args[1:] {
		value, err = env.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func unlessFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	if v == nil {
		return env.Eval(args[1])
	}
	condition, ok := v.(bool)
	if !ok || condition {
		return nil, nil
	}
	var value any
	for _, arg := range args[1:] {
		value, err = env.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func condFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 1); err != nil {
		return nil, err
	}
	for _, n := range args {
		if list, ok := n.(lisp.List); ok {
			if len(list) != 2 {
				return nil, errors.Errorf("expects conditions to be list pairs")
			}
			test, err := env.Eval(list[0])
			if err != nil {
				return nil, err
			}
			if test == nil {
				continue
			} else if testBool, ok := test.(bool); ok && !testBool {
				continue
			}
			value, err := env.Eval(list[1])
			if err != nil {
				return nil, err
			}
			return value, nil
		}
	}
	return nil, nil
}
