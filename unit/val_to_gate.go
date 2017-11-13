package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newValToGate(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &valToGate{
		in:  io.NewIn("in", dsp.Float64(0)),
		out: io.NewOut("out"),
	}), nil
}

type valToGate struct {
	in  *In
	out *Out
}

func (g *valToGate) ProcessSample(i int) {
	in := g.in.Read(i)
	g.out.Write(i, (math.Min(in, 1)-0.5)*2)
}
