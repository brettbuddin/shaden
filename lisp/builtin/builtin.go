package builtin

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"buddin.us/lumen/lisp"
)

func Load(env *lisp.Environment) {
	lispPath := os.Getenv("LUMEN_LISP_PATH")
	loadPath := strings.Split(lispPath, string(filepath.ListSeparator))

	env.DefineSymbol("!=", notEqualFn)
	env.DefineSymbol("*", multFn)
	env.DefineSymbol("+", sumFn)
	env.DefineSymbol("-", diffFn)
	env.DefineSymbol("/", divFn)
	env.DefineSymbol("=", equalFn)
	env.DefineSymbol("and", andFn)
	env.DefineSymbol("append", appendFn)
	env.DefineSymbol("apply", applyFn)
	env.DefineSymbol("boolean?", isBoolFn)
	env.DefineSymbol("cond", condFn)
	env.DefineSymbol("cons", consFn)
	env.DefineSymbol("define", defineFn)
	env.DefineSymbol("define-macro", defineMacroFn)
	env.DefineSymbol("do", doFn)
	env.DefineSymbol("dotimes", dotimesFn)
	env.DefineSymbol("each", eachFn)
	env.DefineSymbol("empty?", isEmptyFn)
	env.DefineSymbol("errorf", errorfFn)
	env.DefineSymbol("eval", evalFn)
	env.DefineSymbol("first", firstFn)
	env.DefineSymbol("float", floatFn)
	env.DefineSymbol("float?", isFloatFn)
	env.DefineSymbol("fn", fnFn)
	env.DefineSymbol("fn?", isFnFn)
	env.DefineSymbol("if", ifFn)
	env.DefineSymbol("int", intFn)
	env.DefineSymbol("int?", isIntFn)
	env.DefineSymbol("keyword", keywordFn)
	env.DefineSymbol("keyword?", isKeywordFn)
	env.DefineSymbol("len", lenFn)
	env.DefineSymbol("let", letFn)
	env.DefineSymbol("list", listFn)
	env.DefineSymbol("list?", isListFn)
	env.DefineSymbol("load", loadFn)
	env.DefineSymbol("load-path", loadPath)
	env.DefineSymbol("map", mapFn)
	env.DefineSymbol("nil?", isNilFn)
	env.DefineSymbol("not", notFn)
	env.DefineSymbol("number?", isNumberFn)
	env.DefineSymbol("or", orFn)
	env.DefineSymbol("pow", powFn)
	env.DefineSymbol("prepend", prependFn)
	env.DefineSymbol("printf", printfFn)
	env.DefineSymbol("println", printlnFn)
	env.DefineSymbol("quasiquote", quasiquoteFn)
	env.DefineSymbol("quote", quoteFn)
	env.DefineSymbol("rand", randFn)
	env.DefineSymbol("rand-intn", randIntnFn)
	env.DefineSymbol("read", readFn)
	env.DefineSymbol("reduce", reduceFn)
	env.DefineSymbol("rest", restFn)
	env.DefineSymbol("set", setFn)
	env.DefineSymbol("sprintf", sprintfFn)
	env.DefineSymbol("string?", isStringFn)
	env.DefineSymbol("symbol?", isSymbolFn)
	env.DefineSymbol("table", tableFn)
	env.DefineSymbol("table-delete", tdeleteFn)
	env.DefineSymbol("table-exists?", texistsFn)
	env.DefineSymbol("table-get", tgetFn)
	env.DefineSymbol("table-merge", mergeFn)
	env.DefineSymbol("table-set", tsetFn)
	env.DefineSymbol("table?", isTableFn)
	env.DefineSymbol("undefine", undefineFn)
	env.DefineSymbol("unless", unlessFn)
	env.DefineSymbol("when", whenFn)
	env.DefineSymbol("zero-value?", isZeroValueFn)
}

func evalFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("eval expects 1 argument")
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	return env.Eval(v)
}

func readFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("read expects 1 argument")
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	str, ok := v.(string)
	if !ok {
		return nil, errors.New("read expects string for argument 1")
	}
	node, err := lisp.Parse(bytes.NewBufferString(str))
	if err != nil {
		return nil, err
	}
	return node, nil
}

func defineFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("define expects at least 2 arguments")
	}

	switch v := args[0].(type) {
	case lisp.Symbol:
		value, err := env.Eval(args[1])
		if err != nil {
			return nil, err
		}
		env.DefineSymbol(string(v), value)
		return nil, nil
	case lisp.List:
		for _, n := range v {
			if _, ok := n.(lisp.Symbol); !ok {
				return nil, errors.New("define expects all function parameters to be symbols")
			}
		}
		name := v[0].(lisp.Symbol)
		fn := buildFunction(env, string(name), v[1:], args[1:])
		env.DefineSymbol(string(name), fn)
		return nil, nil
	default:
		return nil, errors.New("define expects symbol or list for argument 1")
	}
}

func defineMacroFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("define-macro expects at least 2 arguments")
	}

	switch v := args[0].(type) {
	case lisp.List:
		for _, n := range v {
			if _, ok := n.(lisp.Symbol); !ok {
				return nil, errors.New("define-macro expects all function parameters to be symbols")
			}
		}
		name := v[0].(lisp.Symbol)

		var processed lisp.List
		for _, n := range args[1:] {
			p, err := env.QuasiQuoteEval(n)
			if err != nil {
				return nil, err
			}
			processed = append(processed, p)
		}

		fn := buildMacroFunction(env, string(name), v[1:], processed)
		env.DefineSymbol(string(name), fn)
		return nil, nil
	default:
		return nil, errors.New("define expects symbol or list for argument 1")
	}
}

func setFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("set expects 2 arguments")
	}
	symbol, ok := args[0].(lisp.Symbol)
	if !ok {
		return nil, errors.New("set expects symbol for argument 1")
	}
	value, err := env.Eval(args[1])
	if err != nil {
		return nil, err
	}
	return nil, env.SetSymbol(string(symbol), value)
}

func undefineFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("undef expects 1 argument")
	}
	symbol, ok := args[0].(lisp.Symbol)
	if !ok {
		return nil, errors.New("undef expects symbol for argument 1")
	}
	env.UnsetSymbol(string(symbol))
	return nil, nil
}

func printlnFn(args lisp.List) (interface{}, error) {
	fmt.Println(args...)
	return nil, nil
}

func printfFn(args lisp.List) (interface{}, error) {
	format, ok := args[0].(string)
	if !ok {
		return nil, errors.New("printf expects string for argument 1")
	}
	fmt.Printf(format, args[1:]...)
	return nil, nil
}

func sprintfFn(args lisp.List) (interface{}, error) {
	format, ok := args[0].(string)
	if !ok {
		return nil, errors.New("sprintf expects string for argument 1")
	}
	return fmt.Sprintf(format, args[1:]...), nil
}
