package unit

import "buddin.us/shaden/dsp"

func newClockDiv(name string, c Config) (*Unit, error) {
	var config struct {
		Div int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Div == 0 {
		config.Div = 2
	}

	io := NewIO()
	cd := &clockDiv{
		in:  io.NewIn("in", dsp.Float64(0)),
		div: io.NewIn("div", dsp.Float64(config.Div)),
		out: io.NewOut("out"),
	}
	return NewUnit(io, name, cd), nil
}

type clockDiv struct {
	in, div *In
	out     *Out

	tick int
	last float64
}

func (d *clockDiv) ProcessSample(i int) {
	var (
		div = d.div.Read(i)
		in  = d.in.Read(i)
	)

	if float64(d.tick) >= div {
		d.out.Write(i, 1)
		d.tick = 0
	} else {
		d.out.Write(i, -1)
	}
	if d.last < 0 && in > 0 {
		d.tick++
	}
	d.last = in
}
