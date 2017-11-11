package unit

import (
	"buddin.us/lumen/dsp"
)

func newInterpolate(name string, c Config) (*Unit, error) {
	var config struct {
		SmoothTime int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.SmoothTime == 0 {
		config.SmoothTime = 10
	}

	io := NewIO()
	return NewUnit(io, name, &interpolate{
		in:      io.NewIn("in", dsp.Float64(0)),
		min:     io.NewIn("min", dsp.Float64(0)),
		max:     io.NewIn("max", dsp.Float64(1)),
		scale:   io.NewIn("scale", dsp.Float64(1)),
		out:     io.NewOut("out"),
		average: dsp.RollingAverage{Window: int(dsp.DurationInt(config.SmoothTime).Float64())},
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
	)
	in := itrp.average.Tick(dsp.Clamp(itrp.in.Read(i), 0, 1))
	itrp.out.Write(i, dsp.Lerp(min, max, in)*scale)
}
