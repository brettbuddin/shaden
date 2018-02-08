package unit

import "buddin.us/shaden/dsp"

func newPan(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &pan{
		in:  io.NewIn("in", dsp.Float64(0)),
		pan: io.NewIn("pan", dsp.Float64(0)),
		a:   io.NewOut("a"),
		b:   io.NewOut("b"),
	}), nil
}

type pan struct {
	in, pan *In
	a, b    *Out
}

func (p *pan) ProcessSample(i int) {
	pan := dsp.Clamp(p.pan.Read(i), -1, 1)
	in := p.in.Read(i)
	if pan > 0 {
		p.a.Write(i, (1-pan)*in)
		p.b.Write(i, in)
	} else if pan < 0 {
		p.a.Write(i, in)
		p.b.Write(i, (1+pan)*in)
	} else {
		p.a.Write(i, in)
		p.b.Write(i, in)
	}
}
