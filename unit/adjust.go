package unit

import "github.com/brettbuddin/shaden/dsp"

func newAdjust(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &adjust{
		in:   io.NewIn("in", dsp.Float64(0)),
		mult: io.NewIn("mult", dsp.Float64(1)),
		add:  io.NewIn("add", dsp.Float64(0)),
		out:  io.NewOut("out"),
	}), nil
}

type adjust struct {
	in, mult, add *In
	out           *Out
}

func (a *adjust) ProcessSample(i int) {
	var (
		in   = a.in.Read(i)
		mult = a.mult.Read(i)
		add  = a.add.Read(i)
	)
	a.out.Write(i, in*mult+add)
}
