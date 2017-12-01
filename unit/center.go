package unit

import "buddin.us/shaden/dsp"

func newCenter(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &center{
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
