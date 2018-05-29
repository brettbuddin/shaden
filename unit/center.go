package unit

import "github.com/brettbuddin/shaden/dsp"

func newCenter(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &center{
		in:    io.NewIn("in", dsp.Float64(0)),
		out:   io.NewOut("out"),
		block: &dsp.DCBlock{},
	}), nil
}

type center struct {
	in    *In
	out   *Out
	block *dsp.DCBlock
}

func (c *center) ProcessSample(i int) {
	c.out.Write(i, c.block.Tick(c.in.Read(i)))
}
