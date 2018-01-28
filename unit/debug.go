package unit

import (
	"fmt"

	"buddin.us/shaden/dsp"
)

func newDebug(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &debug{
		fmt: io.NewProp("fmt", "%.8f", func(p *Prop, v interface{}) error {
			p.value = v
			return nil
		}),
		in:   io.NewIn("in", dsp.Float64(0)),
		rate: io.NewIn("rate", dsp.Float64(0.1)),
		out:  io.NewOut("out"),
	}), nil
}

type debug struct {
	fmt      *Prop
	in, rate *In
	out      *Out
	tick     int
	lastIn   float64
}

func (d *debug) ProcessSample(i int) {
	var (
		in   = d.in.Read(i)
		rate = dsp.Clamp(d.rate.Read(i), 0.01, 1)
	)

	if d.tick%int(dsp.SampleRate*rate) == 0 {
		if d.lastIn != in {
			fmt.Printf(d.fmt.Value().(string)+"\n", in)
			d.lastIn = in
		}
		d.tick = 0
	}
	d.tick++
	d.out.Write(i, in)
}
