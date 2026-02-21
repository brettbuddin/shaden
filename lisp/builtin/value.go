package builtin

import (
	"fmt"
	"reflect"

	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

// isType returns a lisp builtin that checks whether its single argument is of type T.
func isType[T any]() func(lisp.List) (any, error) {
	return func(args lisp.List) (any, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}
		_, ok := args[0].(T)
		return ok, nil
	}
}

func typeFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case string:
		return lisp.TypeString, nil
	case int:
		return lisp.TypeInt, nil
	case float64:
		return lisp.TypeFloat, nil
	case lisp.Keyword:
		return lisp.TypeKeyword, nil
	case lisp.Symbol:
		return lisp.TypeSymbol, nil
	case lisp.List:
		return lisp.TypeList, nil
	case lisp.Table:
		return lisp.TypeTable, nil
	case func(lisp.List) (any, error):
		return lisp.TypeFunction, nil
	case func(*lisp.Environment, lisp.List) (any, error):
		return lisp.TypeFunction, nil
	default:
		return fmt.Sprintf("%T", v), nil
	}
}

func quoteFn(env *lisp.Environment, args lisp.List) (any, error) {
	if len(args) == 0 {
		return nil, nil
	}
	return args[0], nil
}

func quasiquoteFn(env *lisp.Environment, args lisp.List) (any, error) {
	if len(args) == 0 {
		return nil, nil
	}
	return env.QuasiQuoteEval(args[0])
}

func keywordFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case string:
		return lisp.Keyword(v), nil
	case lisp.Keyword:
		return v, nil
	default:
		return lisp.Keyword(fmt.Sprintf("%v", v)), nil
	}
}

func stringFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case string:
		return v, nil
	case lisp.Keyword:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func isNilFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	return args[0] == nil, nil
}

func isNumberFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch args[0].(type) {
	case float64, int:
		return true, nil
	default:
		return false, nil
	}
}

func isFnFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch args[0].(type) {
	case lisp.Func, lisp.EnvFunc:
		return true, nil
	}
	return reflect.TypeOf(args[0]).Kind() == reflect.Func, nil
}

func symbolFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case string:
		return lisp.Symbol(v), nil
	case lisp.Keyword:
		return lisp.Symbol(v), nil
	default:
		return lisp.Symbol(fmt.Sprintf("%v", v)), nil
	}
}

func isEmptyFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}

	if args[0] == nil {
		return true, nil
	}

	switch v := args[0].(type) {
	case lisp.Table:
		return len(v) == 0, nil
	case lisp.List:
		return len(v) == 0, nil
	case string:
		return len(v) == 0, nil
	default:
		return nil, errors.New("expects table, list or string for argument 1")
	}
}

func intFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return nil, errors.Errorf("expects numeric type for argument 1")
	}
}

func floatFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	f, err := lisp.ExtractFloat64(args[0], 1)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func isDefinedFn(env *lisp.Environment, args lisp.List) (any, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}

	v, ok := args[0].(string)
	if !ok {
		return nil, errors.Errorf("expects a string for argument 1")
	}

	_, err := env.GetSymbol(v)
	if err != nil {
		return false, nil
	}
	return true, nil
}
