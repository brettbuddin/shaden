// Package builtin provides built-in functionality for the lisp interpreter.
package builtin

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"buddin.us/shaden/lisp"
)

// Load populates a lisp.Environment with builtin symbols.
func Load(env *lisp.Environment) {
	lispPath := os.Getenv("SHADEN_LISP_PATH")
	var loadPath lisp.List
	for _, str := range strings.Split(lispPath, string(filepath.ListSeparator)) {
		loadPath = append(loadPath, str)
	}

	env.DefineSymbol("!=", notEqualFn)
	env.DefineSymbol("*", multFn)
	env.DefineSymbol("+", sumFn)
	env.DefineSymbol("-", diffFn)
	env.DefineSymbol("/", divFn)
	env.DefineSymbol("<", lessThanFn)
	env.DefineSymbol("=", equalFn)
	env.DefineSymbol(">", greaterThanFn)
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
	env.DefineSymbol("error?", isErrorFn)
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
	env.DefineSymbol("string-split", stringSplitFn)
	env.DefineSymbol("sprintf", sprintfFn)
	env.DefineSymbol("string", stringFn)
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
}

func evalFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "eval", 1); err != nil {
		return nil, err
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	return env.Eval(v)
}

func readFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "read", 1); err != nil {
		return nil, err
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	str, ok := v.(string)
	if !ok {
		return nil, argExpectError("read", "string", 1)
	}
	node, err := lisp.Parse(bytes.NewBufferString(str))
	if err != nil {
		return nil, err
	}
	return node, nil
}

func defineFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "define", 2); err != nil {
		return nil, err
	}

	switch v := args[0].(type) {
	case lisp.Symbol:
		value, err := env.Eval(args[1])
		if err != nil {
			return nil, err
		}
		return nil, env.DefineSymbol(string(v), value)
	case lisp.List:
		for _, n := range v {
			if _, ok := n.(lisp.Symbol); !ok {
				return nil, errors.New("define expects all function parameters to be symbols")
			}
		}
		name := v[0].(lisp.Symbol)
		fn := buildFunction(env, string(name), v[1:], args[1:])
		return nil, env.DefineSymbol(string(name), fn)
	default:
		return nil, argExpectError("define", "symbol or list", 1)
	}
}

func defineMacroFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "define-macro", 2); err != nil {
		return nil, err
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
		return nil, env.DefineSymbol(string(name), fn)
	default:
		return nil, argExpectError("define-macro", "list", 1)
	}
}

func setFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "set", 2); err != nil {
		return nil, err
	}
	symbol, ok := args[0].(lisp.Symbol)
	if !ok {
		return nil, argExpectError("set", "symbol", 1)
	}
	value, err := env.Eval(args[1])
	if err != nil {
		return nil, err
	}
	return nil, env.SetSymbol(string(symbol), value)
}

func undefineFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "undefine", 1); err != nil {
		return nil, err
	}
	symbol, ok := args[0].(lisp.Symbol)
	if !ok {
		return nil, argExpectError("undefine", "symbol", 1)
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
		return nil, argExpectError("printf", "string", 1)
	}
	fmt.Printf(format, args[1:]...)
	return nil, nil
}

func sprintfFn(args lisp.List) (interface{}, error) {
	format, ok := args[0].(string)
	if !ok {
		return nil, argExpectError("sprintf", "string", 1)
	}
	return fmt.Sprintf(format, args[1:]...), nil
}
