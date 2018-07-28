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
		err := g.Unmount(u)
		return nil, err
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
		if err := unit.Patch(g.graph, leftOut, g.sink.In["l"]); err != nil {
			return nil, errors.Wrap(err, "patch")
		}
		return nil, unit.Patch(g.graph, rightOut, g.sink.In["r"])
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
