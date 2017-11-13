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
		in:  io.NewIn("in", dsp.Float64(0)),
		out: io.NewOut("out"),
	}), nil
}

type debug struct {
	fmt  *Prop
	in   *In
	out  *Out
	tick int
}

func (d *debug) ProcessSample(i int) {
	in := d.in.Read(i)
	if d.tick%int(dsp.SampleRate*0.1) == 0 {
		fmt.Printf(d.fmt.Value().(string)+"\n", in)
		d.tick = 0
	}
	d.tick++
	d.out.Write(i, in)
}
