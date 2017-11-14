package engine

import (
	"fmt"

	"github.com/pkg/errors"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/graph"
	"buddin.us/shaden/unit"
)

// Clear is an action that resets the Engine's state.
func Clear(e *Engine) (interface{}, error) {
	return nil, e.Reset()
}

// MountUnit mounts a Unit into the audio graph.
func MountUnit(u *unit.Unit) func(*graph.Graph) (interface{}, error) {
	return func(g *graph.Graph) (interface{}, error) {
		return nil, u.Attach(g)
	}
}

// UnmountUnit removes a Unit from the audio graph.
func UnmountUnit(u *unit.Unit) func(*graph.Graph) (interface{}, error) {
	return func(g *graph.Graph) (interface{}, error) {
		if err := u.Detach(g); err != nil {
			switch err := err.(type) {
			case graph.NotInGraphError:
				return nil, errors.Errorf("unit %q not in graph", u.ID)
			default:
				return nil, err
			}
		}
		return nil, nil
	}
}

// EmitOutputs sinks 1 or 2 outputs to the Engine.
func EmitOutputs(left, right unit.OutRef) func(*Engine) (interface{}, error) {
	return func(e *Engine) (interface{}, error) {
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
		if err := unit.Patch(e.graph, leftOut, e.unit.In["l"]); err != nil {
			return nil, errors.Wrap(err, "patch")
		}
		return nil, unit.Patch(e.graph, rightOut, e.unit.In["r"])
	}
}

// PatchInput patches values into a Unit's Ins. If `forceReset` is set to `true` all Ins on that Unit that haven't been
// referenced in `inputs` will be reset to their default values.
func PatchInput(u *unit.Unit, inputs map[string]interface{}, forceReset bool) func(*graph.Graph) (interface{}, error) {
	seen := make(map[string]struct{}, len(u.In))
	return func(g *graph.Graph) (interface{}, error) {
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

			switch v := v.(type) {
			case float64:
				if err := unit.Unpatch(g, in); err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("unpatch %q", in))
				}
				in.Fill(dsp.Float64(v))
			case int:
				if err := unit.Unpatch(g, in); err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("unpatch %q", in))
				}
				in.Fill(dsp.Float64(v))
			case dsp.Valuer:
				if err := unit.Unpatch(g, in); err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("unpatch %q", in))
				}
				in.Fill(v)
			case unit.OutRef:
				out, ok := v.Unit.Out[v.Output]
				if !ok {
					return nil, errors.Errorf("unit %q has no output %q", v.Unit.ID, v.Output)
				}
				if err := unit.Patch(g, out, in); err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("patch %q into %q", out.Out(), in))
				}
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
