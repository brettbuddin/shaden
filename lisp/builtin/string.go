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
