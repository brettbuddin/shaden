package builtin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
)

func loadFn(env *lisp.Environment, args lisp.List) (interface{}, error) {
	if err := checkArityEqual(args, "load", 1); err != nil {
		return nil, err
	}

	raw, err := env.Eval(args[0])
	if err != nil {
		return nil, err
	}
	path, ok := raw.(string)
	if !ok {
		return nil, argExpectError("load", "string", 1)
	}

	loadPath, err := env.GetSymbol("load-path")
	if err != nil {
		return nil, err
	}

	var found string
	for _, segment := range loadPath.(lisp.List) {
		fullpath := filepath.Join(segment.(string), path)
		if _, err := os.Stat(fullpath); err == nil {
			found = fullpath
		}
	}

	if found == "" {
		return nil, fmt.Errorf("%s not found n load-path", path)
	}

	f, err := os.Open(found)
	if err != nil {
		return nil, errors.Wrapf(err, "opening %q", found)
	}
	defer f.Close()
	node, err := lisp.Parse(f)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing %q", found)
	}
	v, err := env.Eval(node)
	if err != nil {
		return v, errors.Wrapf(err, "evaluating %q", found)
	}
	return v, nil
}
