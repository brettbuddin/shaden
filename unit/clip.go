package unit

import (
	"github.com/brettbuddin/shaden/dsp"
)

func newClip(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &clip{
		in:    io.NewIn("in", dsp.Float64(0)),
		level: io.NewIn("level", dsp.Float64(1)),
		soft:  io.NewIn("soft", dsp.Float64(1)),
		out:   io.NewOut("out"),
	}), nil
}

type clip struct {
	in, level, soft *In
	out             *Out
}

func (c *clip) ProcessSample(i int) {
	var (
		soft  = c.soft.Read(i)
		in    = c.in.Read(i)
		level = c.level.Read(i)
	)
	if soft == 1 {
		c.out.Write(i, dsp.SoftClamp(in, level))
		return
	}
	c.out.Write(i, dsp.Clamp(in, -level, level))
}
