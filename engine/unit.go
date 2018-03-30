package engine

import (
	"time"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/unit"
)

func newSink(io *unit.IO, fadeIn time.Duration, sampleRate, frameSize int) *sink {
	var (
		fadeInMS      = float64(fadeIn) / 1000.0
		fadeInSamples = dsp.Duration(fadeInMS, sampleRate).Float64()
	)
	return &sink{
		left: &channel{
			fadeIn: fadeInSamples,
			in:     io.NewIn("l", dsp.Float64(0)),
			out:    make([]float64, frameSize),
		},
		right: &channel{
			fadeIn: fadeInSamples,
			in:     io.NewIn("r", dsp.Float64(0)),
			out:    make([]float64, frameSize),
		},
	}
}

type sink struct {
	left, right *channel
}

func (s *sink) ProcessSample(i int) {
	s.left.tick(i)
	s.right.tick(i)
}

type channel struct {
	in        *unit.In
	out       []float64
	level     float64
	hasSignal bool
	fadeIn    float64
}

func (c *channel) tick(i int) {
	in := c.in.Read(i)
	c.out[i] = in * c.level
	if !c.hasSignal && in != 0 {
		c.hasSignal = true
	}
	if c.level < 1 {
		c.level += 1 / c.fadeIn
		if c.level > 1 {
			c.level = 1
		}
	}
}

func newSource(e *Engine) unit.IOBuilder {
	return func(io *unit.IO, _ unit.Config) (*unit.Unit, error) {
		io.NewOutWithFrame("output", e.input)
		return unit.NewUnit(io, nil), nil
	}
}
