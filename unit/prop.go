package unit

import (
	"fmt"

	"buddin.us/shaden/dsp"
)

// PropSetterFunc is a function that will be used when setting the value of a Prop. It provides Modules a point a
// control for the values that are given to its Props.
type PropSetterFunc func(*Prop, interface{}) error

// Prop is a module property
type Prop struct {
	name   string
	setter PropSetterFunc
	value  interface{}
}

// Value returns the Prop's value
func (p *Prop) Value() interface{} {
	return p.value
}

// SetValue sets the Prop's value using its internal PropSetterFunc (if it has one)
func (p *Prop) SetValue(v interface{}) error {
	if p.setter == nil {
		p.value = v
		return nil
	}
	return p.setter(p, v)
}

// InvalidPropValueError is an error that indicates a Prop cannot handle a value that's been given to it
type InvalidPropValueError struct {
	Prop  *Prop
	Value interface{}
}

func (e InvalidPropValueError) Error() string {
	return fmt.Sprintf("invalid value %v (%T) for property %s", e.Value, e.Value, e.Prop.name)
}

func inStringList(l []string) PropSetterFunc {
	return func(p *Prop, v interface{}) error {
		s, ok := v.(string)
		if !ok {
			return InvalidPropValueError{Prop: p, Value: v}
		}
		for _, k := range l {
			if k == s {
				p.value = s
				return nil
			}
		}
		return InvalidPropValueError{Prop: p, Value: v}
	}
}

func clampRange(min, max float64) PropSetterFunc {
	return func(p *Prop, raw interface{}) error {
		switch v := raw.(type) {
		case int:
			p.value = dsp.Clamp(float64(v), min, max)
		case float64:
			p.value = dsp.Clamp(v, min, max)
		default:
			return InvalidPropValueError{Prop: p, Value: v}
		}
		return nil
	}
}
