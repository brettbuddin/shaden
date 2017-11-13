package unit

import "buddin.us/shaden/dsp"

func newOverload(name string, _ Config) (*Unit, error) {
	io := NewIO()
	d := &overload{
		in:   io.NewIn("in", dsp.Float64(0)),
		gain: io.NewIn("gain", dsp.Float64(1)),
		out:  io.NewOut("out"),
	}
	return NewUnit(io, name, d), nil
}

type overload struct {
	in, gain *In
	out      *Out
}

func (o *overload) ProcessSample(i int) {
	o.out.Write(i, dsp.Overload(o.in.Read(i)*o.gain.Read(i)))
}
