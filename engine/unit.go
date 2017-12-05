package engine

import (
	"buddin.us/shaden/dsp"
	"buddin.us/shaden/unit"
)

func unitBuilders(e *Engine) map[string]unit.BuildFunc {
	return map[string]unit.BuildFunc{
		"source": newSource(e),
	}
}

func newSink(fadeIn bool) (*unit.Unit, *sink) {
	io := unit.NewIO()
	s := &sink{
		left: &channel{
			fadeIn: fadeIn,
			in:     io.NewIn("l", dsp.Float64(0)),
			out:    make([]float64, dsp.FrameSize),
		},
		right: &channel{
			fadeIn: fadeIn,
			in:     io.NewIn("r", dsp.Float64(0)),
			out:    make([]float64, dsp.FrameSize),
		},
	}
	return unit.NewUnit(io, "sink", s), s
}

type sink struct {
	left, right *channel
}

func (s *sink) ProcessSample(i int) {
	s.left.tick(i)
	s.right.tick(i)
}

var fadeSamples = dsp.Duration(100).Float64()

type channel struct {
	in        *unit.In
	out       []float64
	level     float64
	hasSignal bool
	fadeIn    bool
}

func (c *channel) tick(i int) {
	in := c.in.Read(i)
	c.out[i] = in * c.level
	if !c.hasSignal && in != 0 {
		c.hasSignal = true
	}
	if c.level < 1 {
		c.level += 1 / fadeSamples
		if c.level > 1 {
			c.level = 1
		}
	}
}

func newSource(e *Engine) unit.BuildFunc {
	return func(unit.Config) (*unit.Unit, error) {
		io := unit.NewIO()
		io.NewOutWithFrame("output", e.input)
		return unit.NewUnit(io, "source", nil), nil
	}
}
