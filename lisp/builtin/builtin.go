// Package builtin provides built-in functionality for the lisp interpreter.
package builtin

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
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
	env.DefineSymbol("begin", beginFn)
	env.DefineSymbol("bool?", isBoolFn)
	env.DefineSymbol("cond", condFn)
	env.DefineSymbol("cons", consFn)
	env.DefineSymbol("define", defineFn)
	env.DefineSymbol("defined?", isDefinedFn)
	env.DefineSymbol("define-macro", defineMacroFn)
	env.DefineSymbol("dotimes", dotimesFn)
	env.DefineSymbol("each", eachFn)
	env.DefineSymbol("empty?", isEmptyFn)
	env.DefineSymbol("error", errorFn)
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
	env.DefineSymbol("read", readFn)
	env.DefineSymbol("reduce", reduceFn)
	env.DefineSymbol("rest", restFn)
	env.DefineSymbol("set!", setFn)
	env.DefineSymbol("sleep", sleepFn)
	env.DefineSymbol("string-split", stringSplitFn)
	env.DefineSymbol("string-join", stringJoinFn)
	env.DefineSymbol("string-has-prefix", stringHasPrefixFn)
	env.DefineSymbol("string-replace", stringReplaceFn)
	env.DefineSymbol("sprintf", sprintfFn)
	env.DefineSymbol("string", stringFn)
	env.DefineSymbol("string?", isStringFn)
	env.DefineSymbol("symbol", symbolFn)
	env.DefineSymbol("symbol?", isSymbolFn)
	env.DefineSymbol("table", tableFn)
	env.DefineSymbol("table-delete!", tdeleteFn)
	env.DefineSymbol("table-exists?", texistsFn)
	env.DefineSymbol("table-get", tgetFn)
	env.DefineSymbol("table-merge", mergeFn)
	env.DefineSymbol("table-set!", tsetFn)
	env.DefineSymbol("table-select", tselectFn)
	env.DefineSymbol("table?", isTableFn)
	env.DefineSymbol("type", typeFn)
	env.DefineSymbol("undefine", undefineFn)
	env.DefineSymbol("unless", unlessFn)
	env.DefineSymbol("when", whenFn)
}

func evalFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	return env.Eval(v)
}

func readFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	v, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	str, ok := v.(string)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeString, 1)
	}
	node, err := lisp.Parse(bytes.NewBufferString(str))
	if err != nil {
		return nil, err
	}
	return node, nil
}

func defineFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
		return nil, err
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
				return nil, errors.New("expects all function parameters to be symbols")
			}
		}
		name := v[0].(lisp.Symbol)
		fn := buildFunction(env, v[1:], args[1:])
		env.DefineSymbol(string(name), fn)
		return nil, nil
	default:
		return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeSymbol, lisp.TypeList), 1)
	}
}

func defineMacroFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityAtLeast(args, 2); err != nil {
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

		fn := buildMacroFunction(env, v[1:], processed)
		env.DefineSymbol(string(name), fn)
		return nil, nil
	default:
		return nil, lisp.ArgExpectError(lisp.TypeList, 1)
	}
}

func setFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 2); err != nil {
		return nil, err
	}
	symbol, ok := args[0].(lisp.Symbol)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeSymbol, 1)
	}
	value, err := env.Eval(args[1])
	if err != nil {
		return nil, err
	}
	return nil, env.SetSymbol(string(symbol), value)
}

func undefineFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	symbol, ok := args[0].(lisp.Symbol)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeSymbol, 1)
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
		return nil, lisp.ArgExpectError(lisp.TypeString, 1)
	}
	fmt.Printf(format, args[1:]...)
	return nil, nil
}

func sprintfFn(args lisp.List) (interface{}, error) {
	format, ok := args[0].(string)
	if !ok {
		return nil, lisp.ArgExpectError(lisp.TypeString, 1)
	}
	return fmt.Sprintf(format, args[1:]...), nil
}
