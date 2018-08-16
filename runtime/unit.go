package runtime

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/brettbuddin/shaden/engine"
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/lisp"
	"github.com/brettbuddin/shaden/midi"
	"github.com/brettbuddin/shaden/unit"
)

const (
	nameUnitUnmount   = "unit-unmount"
	nameUnitRemove    = "unit-remove"
	nameUnitInputs    = "unit-inputs"
	nameUnitOutputs   = "unit-outputs"
	nameUnitType      = "unit-type"
	nameUnitID        = "unit-id"
	nameUnitPatch     = "->"
	nameUnitPatchOnly = "=>"
	nameUnitOutput    = "<-"
	nameEmit          = "emit"

	typeUnit      = "unit"
	typeOutputRef = "output reference"
)

var bold = color.New(color.Bold).SprintFunc()

type lazyUnit struct {
	logger *log.Logger
	engine Engine

	created         *unit.Unit
	inputs, outputs []string
	id, typ         string
	mount           bool
}

func (r *lazyUnit) String() string {
	return fmt.Sprintf("%s(mounted=%v)", r.id, r.mount)
}

func (r *lazyUnit) mounted() (*unit.Unit, error) {
	if r.mount {
		return r.created, nil
	}

	m := engine.NewMessage(engine.MountUnit(r.created))

	if err := r.engine.SendMessage(m); err != nil {
		return nil, err
	}
	reply := <-m.Reply
	if reply.Error != nil {
		return nil, reply.Error
	}
	r.logger.Printf("%s\n└ Completed in %s\n", bold("Adding "+r.created.ID), reply.Duration)
	r.mount = true
	return r.created, nil
}

func createBuilders(env *lisp.Environment, e Engine, logger *log.Logger) error {
	builders, err := unitBuilders(e)
	if err != nil {
		return err
	}
	for name, builder := range builders {
		defineBuilders(env, builder, e, logger, "unit/"+name)
	}
	return nil
}

func unitBuilders(e Engine) (map[string]unit.Builder, error) {
	groups := []map[string]unit.Builder{
		unit.Builders(),
		e.UnitBuilders(),
		midi.UnitBuilders(),
	}
	merged := map[string]unit.Builder{}
	for _, g := range groups {
		for name, c := range g {
			if _, ok := merged[name]; ok {
				return nil, fmt.Errorf("duplicate unit named %q", name)
			}
			merged[name] = c
		}
	}
	return merged, nil
}

func defineBuilders(env *lisp.Environment, builder unit.Builder, e Engine, logger *log.Logger, name string) {
	env.DefineSymbol(name, func(args lisp.List) (interface{}, error) {
		if len(args) > 1 {
			return nil, errors.Errorf("expects at most 1 argument")
		}

		config := map[string]interface{}{}
		if len(args) == 1 {
			m, ok := args[0].(lisp.Table)
			if !ok {
				return nil, lisp.ArgExpectError(lisp.TypeTable, 1)
			}
			for k, v := range m {
				switch k := k.(type) {
				case string:
					config[k] = v
				case lisp.Keyword:
					config[string(k)] = v
				default:
					config[fmt.Sprintf("%v", k)] = v
				}
			}
		}

		unit, err := builder(unit.Config{
			Values:     config,
			SampleRate: e.SampleRate(),
			FrameSize:  e.FrameSize(),
		})
		if err != nil {
			return nil, err
		}

		var inputs, outputs []string
		for k := range unit.In {
			inputs = append(inputs, k)
		}
		for k := range unit.Out {
			outputs = append(outputs, k)
		}
		natsort(inputs)
		natsort(outputs)

		return &lazyUnit{
			logger:  logger,
			engine:  e,
			created: unit,
			id:      unit.ID,
			typ:     unit.Type,
			inputs:  inputs,
			outputs: outputs,
		}, nil
	})
}

func unitRemoveFn(e Engine, logger *log.Logger) func(*lisp.Environment, lisp.List) (interface{}, error) {
	return func(env *lisp.Environment, args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}

		symbol, ok := args[0].(lisp.Symbol)
		if !ok {
			return nil, lisp.ArgExpectError(lisp.TypeSymbol, 1)
		}

		values := make(lisp.List, len(args))
		for i, n := range args {
			value, err := env.Eval(n)
			if err != nil {
				return nil, err
			}
			values[i] = value
		}

		if _, err := unitUnmountFn(e, logger)(values); err != nil {
			return nil, errors.Wrap(err, "unmount unit failed")
		}
		env.UnsetSymbol(string(symbol))
		return nil, nil
	}
}

func unitUnmountFn(e Engine, logger *log.Logger) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityEqual(args, 1); err != nil {
			return nil, err
		}

		lazy, ok := args[0].(*lazyUnit)
		if !ok {
			return nil, lisp.ArgExpectError(typeUnit, 1)
		}

		u, err := lazy.mounted()
		if err != nil {
			return nil, errors.Wrap(err, "retrieving mounted unit failed")
		}

		m := engine.NewMessage(engine.UnmountUnit(u))

		if err := e.SendMessage(m); err != nil {
			return nil, err
		}
		reply := <-m.Reply
		if reply.Error != nil {
			return nil, reply.Error
		}
		logger.Printf(bold("Removing %s\n└ Completed in %s\n"), u.ID, reply.Duration)
		lazy.mount = false
		return nil, nil
	}
}

func patchFn(e Engine, logger *log.Logger, forceReset bool) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if err := lisp.CheckArityAtLeast(args, 2); err != nil {
			return nil, err
		}

		lazy, ok := args[0].(*lazyUnit)
		if !ok {
			return nil, lisp.ArgExpectError(typeUnit, 1)
		}

		inputs, err := patchableInputs(args[1:])
		if err != nil {
			return nil, err
		}

		u, err := lazy.mounted()
		if err != nil {
			return nil, errors.Wrap(err, "retrieving mounted unit failed")
		}

		m := engine.NewMessage(engine.PatchInput(u, inputs, forceReset))

		if err := e.SendMessage(m); err != nil {
			return nil, err
		}
		reply := <-m.Reply
		if reply.Error != nil {
			return nil, reply.Error
		}

		names := make([]string, 0, len(inputs))
		for k := range inputs {
			names = append(names, k)
		}
		natsort(names)

		var b bytes.Buffer
		fmt.Fprintf(&b, bold("Patching %s\n"), u.ID)
		tw := tabwriter.NewWriter(&b, 8, 8, 1, ' ', 0)
		for _, name := range names {
			fmt.Fprintf(tw, "│ %v\t-> %s\n", inputs[name], name)
		}
		tw.Flush()
		fmt.Fprintf(&b, "└ Completed in %s\n", reply.Duration)
		logger.Print(b.String())

		return lazy, nil
	}
}

func patchableInputs(args lisp.List) (map[string]interface{}, error) {
	inputs := map[string]interface{}{}

	if len(args) == 2 {
		switch first := args[0].(type) {
		case string:
			inputs[first] = args[1]
			return inputs, nil
		case lisp.Keyword:
			inputs[string(first)] = args[1]
			return inputs, nil
		}
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case lisp.List:
			for i, e := range v {
				if m, ok := e.(lisp.Table); ok {
					for k, w := range m {
						inputs[fmt.Sprintf("%d/%s", i, k)] = patchableValue(w)
					}
				} else {
					inputs[strconv.Itoa(i)] = patchableValue(e)
				}
			}
		case lisp.Table:
			for k, e := range v {
				switch k := k.(type) {
				case string:
					inputs[k] = patchableValue(e)
				case lisp.Keyword:
					inputs[string(k)] = patchableValue(e)
				default:
					inputs[fmt.Sprint(k)] = patchableValue(e)
				}
			}
		default:
			return nil, typeRemainingError("hash or list", 2)
		}
	}
	return inputs, nil
}

func patchableValue(v interface{}) interface{} {
	switch value := v.(type) {
	case lisp.List:
		s := make([]interface{}, len(value))
		for i, v := range value {
			s[i] = patchableValue(v)
		}
		return s
	case lisp.Table:
		m := map[string]interface{}{}
		for k, v := range value {
			m[fmt.Sprint(k)] = patchableValue(v)
		}
		return m
	default:
		return value
	}
}

func emitFn(e Engine, logger *log.Logger) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, errors.Errorf("expects 1 or 2 arguments")
		}

		left, ok := args[0].(unit.OutRef)
		if !ok {
			return nil, lisp.ArgExpectError(typeOutputRef, 1)
		}

		var right unit.OutRef
		if len(args) > 1 {
			var ok bool
			right, ok = args[1].(unit.OutRef)
			if !ok {
				return nil, lisp.ArgExpectError(typeOutputRef, 2)
			}
		} else {
			right = left
		}

		msg := engine.NewMessage(engine.EmitOutputs(left, right))
		if err := e.SendMessage(msg); err != nil {
			return nil, err
		}
		reply := <-msg.Reply

		var b bytes.Buffer
		fmt.Fprintln(&b, bold("Emitting"))
		tw := tabwriter.NewWriter(&b, 8, 8, 1, ' ', 0)
		fmt.Fprintf(tw, "│ %s\t-> left\n", left)
		fmt.Fprintf(tw, "│ %s\t-> right\n", right)
		tw.Flush()
		fmt.Fprintf(&b, "└ Completed in %s", reply.Duration)
		logger.Print(b.String())

		return reply.Data, reply.Error
	}
}

func outFn(e Engine) func(lisp.List) (interface{}, error) {
	return func(args lisp.List) (interface{}, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, errors.Errorf("expects 1 or 2 arguments")
		}

		lazy, ok := args[0].(*lazyUnit)
		if !ok {
			return nil, lisp.ArgExpectError(typeUnit, 1)
		}

		var output string
		if len(args) == 1 {
			output = "out"
		} else {
			switch arg := args[1].(type) {
			case string:
				output = arg
			case lisp.Keyword:
				output = string(arg)
			default:
				return nil, lisp.ArgExpectError(lisp.AcceptTypes(lisp.TypeString, lisp.TypeKeyword), 2)
			}
		}

		var found bool
		for _, v := range lazy.outputs {
			if v == output {
				found = true
			}
		}
		if !found {
			return nil, errors.Errorf("unit %q has no output %q", lazy.id, output)
		}

		u, err := lazy.mounted()
		if err != nil {
			return nil, err
		}
		return unit.OutRef{Unit: u, Output: output}, nil
	}
}

func unitInputsFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	lazy, ok := args[0].(*lazyUnit)
	if !ok {
		return nil, lisp.ArgExpectError(typeUnit, 1)
	}
	var inputs lisp.List
	for _, in := range lazy.inputs {
		inputs = append(inputs, in)
	}
	return inputs, nil
}

func unitOutputsFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	lazy, ok := args[0].(*lazyUnit)
	if !ok {
		return nil, lisp.ArgExpectError(typeUnit, 1)
	}
	var outputs lisp.List
	for _, out := range lazy.outputs {
		outputs = append(outputs, out)
	}
	return outputs, nil
}

func unitTypeFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	lazy, ok := args[0].(*lazyUnit)
	if !ok {
		return nil, lisp.ArgExpectError(typeUnit, 1)
	}
	return lazy.typ, nil
}

func unitIDFn(args lisp.List) (interface{}, error) {
	if err := lisp.CheckArityEqual(args, 1); err != nil {
		return nil, err
	}
	lazy, ok := args[0].(*lazyUnit)
	if !ok {
		return nil, lisp.ArgExpectError(typeUnit, 1)
	}
	return lazy.id, nil
}
