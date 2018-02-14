package builtin

import (
	"fmt"
	"strings"

	"buddin.us/shaden/lisp"
)

func stringSplitFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "string-split", 2); err != nil {
		return nil, err
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, argExpectError("string-split", "string", 1)
	}
	delim, ok := args[1].(string)
	if !ok {
		return nil, argExpectError("string-split", "string", 2)
	}
	var lst lisp.List
	for _, v := range strings.Split(str, delim) {
		lst = append(lst, v)
	}
	return lst, nil
}

func stringJoinFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "string-join", 2); err != nil {
		return nil, err
	}
	lst, ok := args[0].(lisp.List)
	if !ok {
		return nil, argExpectError("string-join", "list", 1)
	}
	delim, ok := args[1].(string)
	if !ok {
		return nil, argExpectError("string-join", "string", 2)
	}

	var strs []string
	for _, v := range lst {
		strs = append(strs, fmt.Sprintf("%v", v))
	}
	return strings.Join(strs, delim), nil
}

func stringHasPrefixFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "string-has-prefix", 2); err != nil {
		return nil, err
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, argExpectError("string-has-prefix", "string", 1)
	}
	prefix, ok := args[1].(string)
	if !ok {
		return nil, argExpectError("string-has-prefix", "string", 2)
	}
	return strings.HasPrefix(str, prefix), nil
}

func stringReplaceFn(args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "string-replace", 4); err != nil {
		return nil, err
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, argExpectError("string-replace", "string", 1)
	}
	old, ok := args[1].(string)
	if !ok {
		return nil, argExpectError("string-replace", "string", 2)
	}
	replacement, ok := args[2].(string)
	if !ok {
		return nil, argExpectError("string-replace", "string", 3)
	}
	occurances, ok := args[3].(int)
	if !ok {
		return nil, argExpectError("string-replace", "integer", 4)
	}
	return strings.Replace(str, old, replacement, occurances), nil
}
