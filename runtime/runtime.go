package runtime

import (
	"bytes"
	"fmt"
	"log"
	"os"

	prompt "github.com/c-bata/go-prompt"
	"github.com/pkg/errors"

	"buddin.us/shaden/engine"
	"buddin.us/shaden/lisp"
	"buddin.us/shaden/lisp/builtin"
	"buddin.us/shaden/unit"
)

// Engine represents the things we need from engine.Engine
type Engine interface {
	Messages() chan<- *engine.Message
	UnitBuilders() map[string]unit.BuildFunc
}

// Runtime represents the runtime execution environment
type Runtime struct {
	base, user *lisp.Environment
	engine     Engine
	logger     *log.Logger
}

// New returns a new Runtime
func New(e Engine, logger *log.Logger) (*Runtime, error) {
	base := lisp.NewEnvironment()
	builtin.Load(base)
	r := &Runtime{
		base:   base,
		user:   base.Branch(),
		engine: e,
		logger: logger,
	}
	if err := loadShaden(r); err != nil {
		return nil, err
	}
	return r, nil
}

// ClearUserspace resets the environment by clearing all user-defined symbols.
func (r *Runtime) ClearUserspace() {
	r.user = r.base.Branch()
}

// REPL runs the REPL.
func (r *Runtime) REPL(done chan struct{}) {
	prompt.New(
		func(in string) {
			if len(in) > 0 && in[0] == ';' {
				return
			}
			result, err := r.Eval([]byte(in))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(result)
		},
		func(in prompt.Document) []prompt.Suggest {
			s := []prompt.Suggest{}
			return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
		},
		prompt.OptionPrefixTextColor(prompt.DefaultColor),
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("shaden"),
	).Run()
	close(done)
}

// Eval parses and evaluates lisp expressions.
func (r *Runtime) Eval(code []byte) (interface{}, error) {
	node, err := lisp.Parse(bytes.NewBuffer(code))
	if err != nil {
		return nil, err
	}
	v, err := r.user.Eval(node)
	if err != nil {
		return v, errors.Wrapf(err, "evaluating <string>")
	}
	return v, nil
}

// Load parses and evaluates lisp expressions in a file.
func (r *Runtime) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	node, err := lisp.Parse(f)
	if err != nil {
		return errors.Wrapf(err, "parsing %q", path)
	}

	if _, err := r.user.Eval(node); err != nil {
		return errors.Wrapf(err, "evaluating %q", path)
	}
	return nil
}

func loadShaden(r *Runtime) error {
	env := r.base

	loadConstants(env)

	// Music Theory
	env.DefineSymbol("theory/pitch", pitchFn)
	env.DefineSymbol("theory/interval", intervalFn)
	env.DefineSymbol("theory/transpose", transposeFn)

	// Values
	env.DefineSymbol("hz", hzFn)
	env.DefineSymbol("ms", msFn)
	env.DefineSymbol("bpm", bpmFn)
	env.DefineSymbol("db", dbFn)

	engine := r.engine
	logger := r.logger

	// Engine
	env.DefineSymbol("emit", emitFn(engine, logger))
	env.DefineSymbol("clear", engineClear(r))

	// Units
	if err := createBuilders(env, engine, logger); err != nil {
		return err
	}
	env.DefineSymbol(nameUnitID, unitIDFn)
	env.DefineSymbol(nameUnitType, unitTypeFn)
	env.DefineSymbol(nameUnitInputs, unitInputsFn)
	env.DefineSymbol(nameUnitOutputs, unitOutputsFn)
	env.DefineSymbol(nameUnitUnmount, unitUnmountFn(engine, logger))
	env.DefineSymbol(nameUnitRemove, unitRemoveFn(engine, logger))
	env.DefineSymbol(nameUnitPatch, patchFn(engine, logger, true))
	env.DefineSymbol(nameUnitPatchOnly, patchFn(engine, logger, false))
	env.DefineSymbol(nameUnitOutput, outFn(engine))

	return nil
}

func loadConstants(env *lisp.Environment) {
	// Basic Modes
	env.DefineSymbol("mode/on", 1)
	env.DefineSymbol("mode/off", 0)
	env.DefineSymbol("mode/high", 1)
	env.DefineSymbol("mode/low", -1)

	// Stage Modes
	env.DefineSymbol("mode/rest", 0)
	env.DefineSymbol("mode/first", 1)
	env.DefineSymbol("mode/last", 2)
	env.DefineSymbol("mode/all", 3)
	env.DefineSymbol("mode/hold", 4)

	// Sequence Modes
	env.DefineSymbol("mode/forward", 0)
	env.DefineSymbol("mode/reverse", 1)
	env.DefineSymbol("mode/pingpong", 2)
	env.DefineSymbol("mode/random", 3)

	// LPG Modes
	env.DefineSymbol("mode/lp", 0)
	env.DefineSymbol("mode/both", 1)
	env.DefineSymbol("mode/amp", 2)

	// Note Qualities
	env.DefineSymbol("quality/perfect", 0)
	env.DefineSymbol("quality/minor", 1)
	env.DefineSymbol("quality/major", 2)
	env.DefineSymbol("quality/diminished", 3)
	env.DefineSymbol("quality/augmented", 4)
}

func engineClear(r *Runtime) func(*lisp.Environment, lisp.List) (interface{}, error) {
	return func(*lisp.Environment, lisp.List) (interface{}, error) {
		msg := engine.NewMessage(engine.Clear)
		r.engine.Messages() <- msg
		reply := <-msg.Reply
		if reply.Error != nil {
			return nil, reply.Error
		}
		r.ClearUserspace()
		return nil, nil
	}
}
