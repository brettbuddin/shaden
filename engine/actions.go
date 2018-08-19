package engine

import (
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/unit"
)

// Clear is an action that resets the Engine's state.
func Clear(e *Engine) (interface{}, error) {
	return nil, e.Reset()
}

// MountUnit mounts a Unit into the audio graph.
func MountUnit(u *unit.Unit) func(*Graph) (interface{}, error) {
	return func(g *Graph) (interface{}, error) {
		return nil, g.Mount(u)
	}
}

// UnmountUnit removes a Unit from the audio graph.
func UnmountUnit(u *unit.Unit) func(*Graph) (interface{}, error) {
	return func(g *Graph) (interface{}, error) {
		return nil, g.Unmount(u)
	}
}

// EmitOutputs sinks 1 or 2 outputs to the Engine.
func EmitOutputs(left, right unit.OutRef) func(*Graph) (interface{}, error) {
	return func(g *Graph) (interface{}, error) {
		leftOut, ok := left.Unit.Out[left.Output]
		if !ok {
			return nil, errors.Errorf("unit %q has no output %q", left.Unit.ID, left.Output)
		}
		var rightOut unit.Output
		if right.Unit == nil {
			rightOut = leftOut
		} else {
			var ok bool
			rightOut, ok = right.Unit.Out[right.Output]
			if !ok {
				return nil, errors.Errorf("unit %s has no output %q", right.Unit.ID, right.Output)
			}
		}
		if err := g.Patch(leftOut, g.sink.In["l"]); err != nil {
			return nil, errors.Wrap(err, "patch")
		}
		return nil, g.Patch(rightOut, g.sink.In["r"])
	}
}

// PatchInput patches values into a Unit's Ins. If `forceReset` is set to `true` all Ins on that Unit that haven't been
// referenced in `inputs` will be reset to their default values.
func PatchInput(u *unit.Unit, inputs map[string]interface{}, forceReset bool) func(*Graph) (interface{}, error) {
	seen := make(map[string]struct{}, len(u.In))
	return func(g *Graph) (interface{}, error) {
		for k, v := range inputs {
			in, ok := u.In[k]
			if !ok {
				prop, ok := u.Prop[k]
				if !ok {
					return nil, errors.Errorf("unit %q has no input or property %q", u.ID, k)
				}
				if err := prop.SetValue(v); err != nil {
					return nil, err
				}
				continue
			}
			seen[k] = struct{}{}

			if err := g.Patch(v, in); err != nil {
				return nil, err
			}
		}
		if forceReset {
			for k, v := range u.In {
				if _, ok := seen[k]; !ok {
					v.Reset()
				}
			}
		}
		return nil, nil
	}
}

// SwapUnit swaps one unit out in the graph for another. If the two units or of
// different types, the original is just removed and nothing is done. Otherwise,
// it tries its best to patch sources and destinatinos from the original unit to
// the new unit.
func SwapUnit(u1, u2 *unit.Unit) func(*Graph) (interface{}, error) {
	return func(g *Graph) (interface{}, error) {
		if u1.Type != u2.Type {
			return nil, g.Unmount(u1)
		}

		for k, u1in := range u1.In {
			u2in, ok := u2.In[k]
			if !ok {
				continue
			}
			if u1in.HasSource() {
				if err := g.Patch(u1in.Source(), u2in); err != nil {
					return nil, err
				}
				if err := g.Unpatch(u1in); err != nil {
					return nil, err
				}
			} else {
				if err := g.Patch(u1in.Constant(), u2in); err != nil {
					return nil, err
				}
			}
		}

		for k, u1out := range u1.Out {
			u2out, ok := u2.Out[k]
			if !ok {
				continue
			}
			for _, in := range u1out.Out().Destinations() {
				if err := g.Patch(u2out, in); err != nil {
					return nil, err
				}
			}
		}

		return nil, g.Unmount(u1)
	}
}
