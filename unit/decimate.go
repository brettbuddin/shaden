package unit

import "github.com/brettbuddin/shaden/dsp"

func newDecimate(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &decimate{
		decimate: &dsp.Decimate{
			SampleRate: float64(c.SampleRate),
		},
		in:   io.NewIn("in", dsp.Float64(0)),
		rate: io.NewIn("rate", dsp.Float64(c.SampleRate)),
		bits: io.NewIn("bits", dsp.Float64(24)),
		out:  io.NewOut("out"),
	}), nil
}

type decimate struct {
	in, rate, bits *In
	out            *Out
	decimate       *dsp.Decimate
}

func (d *decimate) ProcessSample(i int) {
	var (
		in   = d.in.Read(i)
		rate = d.rate.ReadSlow(i, ident)
		bits = d.bits.ReadSlow(i, ident)
	)
	d.out.Write(i, d.decimate.Tick(in, rate, bits))
}
