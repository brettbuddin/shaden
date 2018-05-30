package unit

import (
	"github.com/brettbuddin/shaden/dsp"
)

func newCond(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &cond{
		cond: io.NewIn("cond", dsp.Float64(0)),
		x:    io.NewIn("x", dsp.Float64(0)),
		y:    io.NewIn("y", dsp.Float64(0)),
		out:  io.NewOut("out"),
	}), nil
}

type cond struct {
	cond, x, y *In
	out        *Out
}

func (c *cond) ProcessSample(i int) {
	var (
		cond = c.cond.Read(i)
		a    = c.x.Read(i)
		b    = c.y.Read(i)
	)

	if cond > 0 {
		c.out.Write(i, a)
	} else {
		c.out.Write(i, b)
	}
}
