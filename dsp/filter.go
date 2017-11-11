package dsp

import "math"

// SVFilter is a state-variable filter that yields simultaneous lowpass, bandpass and highpass outputs
type SVFilter struct {
	Poles             int
	Cutoff, Resonance float64

	lastCutoff, g, state1, state2 float64
}

// Tick advances the operation
func (f *SVFilter) Tick(in float64) (lp, bp, hp float64) {
	cutoff := math.Abs(f.Cutoff)
	if cutoff != f.lastCutoff {
		f.g = Tan(cutoff)
	}
	f.lastCutoff = cutoff

	r := 1 / math.Max(f.Resonance, 1)
	h := 1 / (1 + r*f.g + f.g*f.g)

	for j := 0; j < f.Poles; j++ {
		hp = h * (in - r*f.state1 - f.g*f.state1 - f.state2)
		bp = f.g*hp + f.state1
		lp = f.g*bp + f.state2

		f.state1 = f.g*hp + bp
		f.state2 = f.g*bp + lp
	}
	return
}

// FilterTypes for Filter
const (
	LowPass FilterType = iota
	HighPass
)

// FilterType describes what frequencies will be allowed/eliminated in an outgoing signal.
type FilterType int

// Filter is a simple low-pass or high-pass filter.
type Filter struct {
	Type              FilterType
	Poles             int
	Cutoff, Resonance float64
	after             []float64
}

// NewFilter returns a new Filter.
func NewFilter(typ FilterType, poles int) *Filter {
	return &Filter{
		Type:  typ,
		Poles: poles,
		after: make([]float64, poles),
	}
}

// Tick advances the operation
func (f *Filter) Tick(in float64) float64 {
	var (
		out, res float64
		cutoff   = 2 * math.Pi * math.Abs(f.Cutoff)
	)

	if f.Resonance > 0 {
		res = f.after[f.Poles-1]
	}
	out = in - (res * Clamp(f.Resonance, 0, 4.5))

	clip := out
	if clip > 1 {
		clip = 1
	} else if clip < -1 {
		clip = -1
	}

	out = clip + ((-clip + out) * 0.995)

	for i := 0; i < f.Poles; i++ {
		if i == 0 {
			f.after[0] += (-f.after[0] + out) * cutoff
		} else {
			f.after[i] += (-f.after[i] + f.after[i-1]) * cutoff
		}
	}

	if f.Type == HighPass {
		return out - f.after[f.Poles-1]
	}
	return f.after[f.Poles-1]
}
