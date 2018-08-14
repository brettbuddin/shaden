package builtin

import (
	"time"

	"github.com/brettbuddin/shaden/lisp"
)

func sleepFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, 1); err != nil {
		return nil, err
	}

	var d time.Duration
	switch v := args[0].(type) {
	case int:
		d = time.Duration(v)
	case float64:
		d = time.Duration(v)
	default:
		return nil, argExpectError(acceptTypes(typeInt, typeFloat), 1)
	}

	time.Sleep(d * time.Second)

	return nil, nil
}
