package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

var log1 = math.Log(0.1)

func newDynamics(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &dynamics{
		in:        io.NewIn("in", dsp.Float64(0)),
		control:   io.NewIn("control", dsp.Float64(0)),
		threshold: io.NewIn("threshold", dsp.Float64(0.5)),
		above:     io.NewIn("above", dsp.Float64(0.3)),
		below:     io.NewIn("below", dsp.Float64(1)),
		clamp:     io.NewIn("clamp", dsp.Duration(10, c.SampleRate)),
		relax:     io.NewIn("relax", dsp.Duration(10, c.SampleRate)),
		out:       io.NewOut("out"),
		dcBlock:   &dsp.DCBlock{},
		lastClamp: -1,
		lastRelax: -1,
		slope:     1 / float64(c.FrameSize),
	}), nil
}

type dynamics struct {
	in, control, threshold, above, below, clamp, relax *In
	out                                                *Out

	clampCoef, relaxCoef float64
	lastClamp, lastRelax float64
	lastGain, lastMax    float64
	slope                float64

	dcBlock *dsp.DCBlock
}

func (d *dynamics) ProcessSample(i int) {
	var (
		in        = d.in
		control   = d.control
		threshold = d.threshold.Read(i)
		above     = d.above.Read(i)
		below     = d.below.Read(i)
		clamp     = d.clamp.Read(i)
		relax     = d.relax.Read(i)
	)

	d.calcCoefs(clamp, relax)

	v := math.Abs(control.Read(i))
	if v < d.lastMax {
		v = v + (d.lastMax-v)*d.relaxCoef
	} else {
		v = v + (d.lastMax-v)*d.clampCoef
	}
	d.lastMax = v

	var nextGain float64
	if d.lastMax < threshold {
		if below == 1 {
			nextGain = 1
		} else {
			nextGain = math.Pow(d.lastMax/threshold, below-1)
			absGain := math.Abs(nextGain)
			if absGain < 1.0e-15 {
				nextGain = 0
			} else if absGain > 1.0e15 {
				nextGain = 1
			}
		}
	} else {
		if above == 1 {
			nextGain = 1
		} else {
			nextGain = math.Pow(d.lastMax/threshold, above-1)
		}
	}

	slope := (nextGain - d.lastGain) * d.slope
	d.out.Write(i, d.dcBlock.Tick(in.Read(i)*d.lastGain))
	d.lastGain += slope
}

func (d *dynamics) calcCoefs(clamp, relax float64) {
	if clamp != d.lastClamp || d.lastClamp == -1 {
		if clamp == 0 {
			d.clampCoef = 0
		} else {
			d.clampCoef = math.Exp(log1 / clamp)
		}
		d.lastClamp = clamp
	}
	if relax != d.lastRelax || d.lastRelax == -1 {
		if relax == 0 {
			d.relaxCoef = 0
		} else {
			d.relaxCoef = math.Exp(log1 / relax)
		}
		d.lastRelax = relax
	}
}
