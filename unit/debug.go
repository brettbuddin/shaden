package unit

import (
	"fmt"

	"github.com/brettbuddin/shaden/dsp"
)

func newDebug(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &debug{
		fmt: io.NewProp("fmt", "%.8f", func(p *Prop, v interface{}) error {
			p.value = v
			return nil
		}),
		in:         io.NewIn("in", dsp.Float64(0)),
		rate:       io.NewIn("rate", dsp.Float64(0.1)),
		out:        io.NewOut("out"),
		sampleRate: c.SampleRate,
	}), nil
}

type debug struct {
	fmt              *Prop
	in, rate         *In
	out              *Out
	sampleRate, tick int
	lastIn           float64
}

func (d *debug) ProcessSample(i int) {
	var (
		in   = d.in.Read(i)
		rate = dsp.Clamp(d.rate.Read(i), 0.01, 1)
	)

	if d.tick%d.sampleRate*int(rate) == 0 {
		if d.lastIn != in {
			fmt.Printf(d.fmt.Value().(string)+"\n", in)
			d.lastIn = in
		}
		d.tick = 0
	}
	d.tick++
	d.out.Write(i, in)
}
