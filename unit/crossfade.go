package unit

import "github.com/brettbuddin/shaden/dsp"

func newCrossfade(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &crossfade{
		a:   io.NewIn("a", dsp.Float64(0)),
		b:   io.NewIn("b", dsp.Float64(0)),
		mix: io.NewIn("mix", dsp.Float64(0)),
		out: io.NewOut("out"),
	}), nil
}

type crossfade struct {
	a, b, mix *In
	out       *Out
}

func (c *crossfade) ProcessSample(i int) {
	c.out.Write(i, dsp.Mix(c.mix.Read(i), c.a.Read(i), c.b.Read(i)))
}
