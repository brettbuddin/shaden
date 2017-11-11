package unit

import "buddin.us/lumen/dsp"

func newDecimate(name string, _ Config) (*Unit, error) {
	io := NewIO()
	d := &decimate{
		decimate: &dsp.Decimate{},
		in:       io.NewIn("in", dsp.Float64(0)),
		rate:     io.NewIn("rate", dsp.Float64(dsp.SampleRate)),
		bits:     io.NewIn("bits", dsp.Float64(24)),
		out:      io.NewOut("out"),
	}
	return NewUnit(io, name, d), nil
}

type decimate struct {
	in, rate, bits *In
	out            *Out
	decimate       *dsp.Decimate
}

func (d *decimate) ProcessSample(i int) {
	d.out.Write(i, d.decimate.Tick(d.in.Read(i), d.rate.Read(i), d.bits.Read(i)))
}
