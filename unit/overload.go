package unit

import "buddin.us/shaden/dsp"

func newOverload(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &overload{
		in:   io.NewIn("in", dsp.Float64(0)),
		gain: io.NewIn("gain", dsp.Float64(1)),
		out:  io.NewOut("out"),
	}), nil
}

type overload struct {
	in, gain *In
	out      *Out
}

func (o *overload) ProcessSample(i int) {
	var (
		in   = o.in.Read(i)
		gain = o.gain.ReadSlow(i, ident)
	)
	o.out.Write(i, dsp.Overload(in*gain))
}
