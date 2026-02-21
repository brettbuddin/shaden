package lisp

import "github.com/brettbuddin/shaden/errors"

// TailCall is a sentinel value returned by functions to signal that the result
// is a deferred evaluation. Instead of recursing into Eval, the caller's
// trampoline loop picks up the Node and Env and evaluates them iteratively.
// This turns recursive tail calls into a flat loop, preventing stack overflow.
type TailCall struct {
	Node    any
	Env     *Environment
	ErrWrap string // if non-empty, wrap errors with this context
}

// ResolveTailCalls evaluates any chain of TailCall values to produce a final
// result. This is needed when a caller requires a fully-evaluated value (e.g.
// macro expansion, higher-order function callbacks) rather than allowing the
// TailCall to propagate upward to the Eval trampoline.
func ResolveTailCalls(result any, err error) (any, error) {
	var errWrap string
	for {
		if err != nil {
			if errWrap != "" {
				err = errors.Wrap(err, errWrap)
			}
			return nil, err
		}
		tc, ok := result.(TailCall)
		if !ok {
			return result, nil
		}
		if errWrap == "" && tc.ErrWrap != "" {
			errWrap = tc.ErrWrap
		}
		result, err = tc.Env.eval(tc.Node)
	}
}
