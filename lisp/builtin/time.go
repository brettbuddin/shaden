package builtin

import (
	"time"

	"github.com/brettbuddin/shaden/lisp"
)

func sleepFn(args lisp.List) (any, error) {
	if err := lisp.CheckArityAtLeast(args, 1); err != nil {
		return nil, err
	}
	f, err := lisp.ExtractFloat64(args[0], 1)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(f) * time.Second)
	return nil, nil
}
