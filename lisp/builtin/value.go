package builtin

import (
	"errors"
	"fmt"
	"reflect"

	"buddin.us/shaden/lisp"
)

func quoteFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	return args[0], nil
}

func quasiquoteFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	return env.QuasiQuoteEval(args[0])
}

func keywordFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "keyword", 1); err != nil {
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

func stringFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "string", 1); err != nil {
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

func isErrorFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "error?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(error)
	return ok, nil
}

func isNilFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "nil?", 1); err != nil {
		return nil, err
	}
	return args[0] == nil, nil
}

func isStringFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "string?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(string)
	return ok, nil
}

func isBoolFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "boolean?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(bool)
	return ok, nil
}

func isIntFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "int?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(int)
	return ok, nil
}

func isFloatFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "float?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(float64)
	return ok, nil
}

func isNumberFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "number?", 1); err != nil {
		return nil, err
	}
	switch args[0].(type) {
	case float64, int:
		return true, nil
	default:
		return false, nil
	}
}

func isFnFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "fn?", 1); err != nil {
		return nil, err
	}
	switch args[0].(type) {
	case lisp.Func, lisp.EnvFunc:
		return true, nil
	}
	return reflect.TypeOf(args[0]).Kind() == reflect.Func, nil
}

func isKeywordFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "keyword?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(lisp.Keyword)
	return ok, nil
}

func isSymbolFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "symbol?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(lisp.Symbol)
	return ok, nil
}

func symbolFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "symbol", 1); err != nil {
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

func isListFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "list?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(lisp.List)
	return ok, nil
}

func isTableFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "table?", 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(lisp.Table)
	return ok, nil
}

func isEmptyFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "empty?", 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case lisp.Table:
		return len(v) == 0, nil
	case lisp.List:
		return len(v) == 0, nil
	case string:
		return len(v) == 0, nil
	default:
		return nil, errors.New("empty? expects table, list or string for argument 1")
	}
}

func intFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "int", 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return nil, fmt.Errorf("int expects numeric type for argument 1")
	}
}

func floatFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "float", 1); err != nil {
		return nil, err
	}
	switch v := args[0].(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return nil, fmt.Errorf("float expects numeric type for argument 1")
	}
}

func isDefinedFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "defined?", 1); err != nil {
		return nil, err
	}

	v, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("defined? expects a string for argument 1")
	}

	_, err := env.GetSymbol(v)
	if err != nil {
		return false, nil
	}
	return true, nil
}
