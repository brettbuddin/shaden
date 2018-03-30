package unit

import (
	"buddin.us/shaden/dsp"
)

func newClock(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &clock{
		tempo:      io.NewIn("tempo", dsp.Frequency(1, c.SampleRate)),
		pw:         io.NewIn("pulse-width", dsp.Float64(0.1)),
		shuffle:    io.NewIn("shuffle", dsp.Float64(0)),
		run:        io.NewIn("run", dsp.Float64(1)),
		out:        io.NewOut("out"),
		sampleRate: float64(c.SampleRate),
	}), nil
}

type clock struct {
	tempo, pw, shuffle, run *In
	out                     *Out

	tick       int
	sampleRate float64
	even       bool
}

func clockShuffle(v float64) float64 { return dsp.Clamp(v, -0.5, 0.5) }

func (c *clock) ProcessSample(i int) {
	var (
		pw      = c.pw.ReadSlow(i, ident)
		shuffle = c.shuffle.ReadSlow(i, clockShuffle)
		tempo   = c.tempo.ReadSlow(i, ident)
		duty    = 1 / (tempo * c.sampleRate) * c.sampleRate
		offset  = duty * shuffle
	)

	if c.run.Read(i) <= 0 {
		c.out.Write(i, -1)
		return
	}

	if !c.even {
		offset *= -1
	}

	c.advance(i, pw, duty+offset)
}

func (c *clock) advance(i int, pw, duty float64) {
	if c.tick < int(duty) {
		c.write(i, pw, duty)
		c.tick++
		return
	}

	c.tick = 0
	c.even = !c.even
	c.write(i, pw, duty)
	c.tick++
}

func (c *clock) write(i int, pw, duty float64) {
	if float64(c.tick) <= pw*duty {
		c.out.Write(i, 1)
	} else {
		c.out.Write(i, -1)
	}
}
