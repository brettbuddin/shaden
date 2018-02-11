package builtin

import (
	"time"

	"buddin.us/shaden/lisp"
)

func sleepFn(args lisp.List) (interface{}, error) {
	if err := checkArityAtLeast(args, "sleep", 1); err != nil {
		return nil, err
	}

	var d time.Duration
	switch v := args[0].(type) {
	case int:
		d = time.Duration(v)
	case float64:
		d = time.Duration(v)
	default:
		return nil, argExpectError("sleep", "integer or float", 1)
	}

	time.Sleep(d * time.Second)

	return nil, nil
}
