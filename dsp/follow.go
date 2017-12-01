package dsp

import "math"

// Follow is an envelope follower
type Follow struct {
	Rise, Fall float64
	env        float64
}

// Tick advances the follower's state
func (f *Follow) Tick(in float64) float64 {
	in = math.Abs(in)
	if in == f.env {
		return f.env
	}
	slope := f.Fall
	if in > f.env {
		slope = f.Rise
	}
	f.env = math.Pow(0.01, 1.0/slope)*(f.env-in) + in
	return f.env
}
