package unit

import "buddin.us/lumen/dsp"

func newAdjust(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &adjust{
		in:     io.NewIn("in", dsp.Float64(0)),
		gain:   io.NewIn("gain", dsp.Float64(1)),
		offset: io.NewIn("offset", dsp.Float64(0)),
		out:    io.NewOut("out"),
	}), nil
}

type adjust struct {
	in, gain, offset *In
	out              *Out
}

func (a *adjust) ProcessSample(i int) {
	var (
		in   = a.in.Read(i)
		mult = a.gain.Read(i)
		add  = a.offset.Read(i)
	)
	a.out.Write(i, in*mult+add)
}
