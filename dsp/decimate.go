package dsp

import "math"

// Decimate reduces the resolution/bandwidth of a signal to produce bitcrushing effects
type Decimate struct {
	SampleRate  float64
	count, last float64
}

// Tick advances the filter's state
func (d *Decimate) Tick(in, rate, bits float64) float64 {
	var (
		step, stepRatio, ratio float64
	)
	if bits >= 31 || bits < 0 {
		step = 0
		stepRatio = 1
	} else {
		step = math.Pow(0.5, bits-0.999)
		stepRatio = 1 / step
	}

	if rate >= d.SampleRate {
		ratio = 1
	} else {
		ratio = rate / d.SampleRate
	}

	d.count += ratio
	if d.count >= 1 {
		d.count--
		var x float64 = 1
		if in < 0 {
			x = -1
		}
		_, frac := math.Modf((in + x*step*0.5) * stepRatio)
		delta := frac * step
		d.last = in - delta
	}
	return d.last
}
