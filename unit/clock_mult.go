package unit

import (
	"buddin.us/lumen/dsp"
)

func newClockMult(name string, c Config) (*Unit, error) {
	var config struct {
		Mult int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Mult == 0 {
		config.Mult = 2
	}

	io := NewIO()
	cm := &clockMult{
		in:   io.NewIn("in", dsp.Float64(0)),
		mult: io.NewIn("mult", dsp.Float64(config.Mult)),
		out:  io.NewOut("out"),
	}
	return NewUnit(io, name, cm), nil
}

type clockMult struct {
	in, mult *In
	out      *Out

	learn struct {
		rate, last float64
	}
	rate float64
	tick int
}

func (m *clockMult) ProcessSample(i int) {
	var (
		mult = m.mult.Read(i)
		in   = m.in.Read(i)
	)

	if m.learn.last < 0 && in > 0 {
		m.rate = (m.rate + m.learn.rate) * 0.5
		m.learn.rate = 0
		m.tick = 0
	}
	m.learn.rate++
	m.learn.last = in

	if m.tick == 0 || float64(m.tick) >= m.rate/mult {
		m.out.Write(i, 1)
		m.tick = 0
	} else {
		m.out.Write(i, -1)
	}
	m.tick++
}
