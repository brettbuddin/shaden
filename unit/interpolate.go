package unit

import (
	"buddin.us/shaden/dsp"
)

func newInterpolate(io *IO, c Config) (*Unit, error) {
	var config struct {
		SmoothTime dsp.MS
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	return NewUnit(io, &interpolate{
		in:      io.NewIn("in", dsp.Float64(0)),
		min:     io.NewIn("min", dsp.Float64(0)),
		max:     io.NewIn("max", dsp.Float64(1)),
		scale:   io.NewIn("scale", dsp.Float64(1)),
		out:     io.NewOut("out"),
		average: dsp.RollingAverage{Window: int(config.SmoothTime.Float64())},
	}), nil
}

type interpolate struct {
	in, min, max, scale *In
	out                 *Out
	average             dsp.RollingAverage
}

func (itrp *interpolate) ProcessSample(i int) {
	var (
		max   = itrp.max.Read(i)
		min   = itrp.min.Read(i)
		scale = itrp.scale.Read(i)
		in    = dsp.Clamp(itrp.in.Read(i), 0, 1)
	)
	if itrp.average.Window != 0 {
		in = itrp.average.Tick(in)
	}
	itrp.out.Write(i, dsp.Lerp(min, max, in)*scale)
}
