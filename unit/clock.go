package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newClock(name string, _ Config) (*Unit, error) {
	io := NewIO()
	c := &clock{
		tempo:   io.NewIn("tempo", dsp.Frequency(1)),
		pw:      io.NewIn("pulse-width", dsp.Float64(0.1)),
		shuffle: io.NewIn("shuffle", dsp.Float64(0)),
		run:     io.NewIn("run", dsp.Float64(1)),
		out:     io.NewOut("out"),
	}
	return NewUnit(io, name, c), nil
}

type clock struct {
	tempo, pw, shuffle, run *In
	out                     *Out

	tick int
	even bool
}

func (c *clock) ProcessSample(i int) {
	var (
		pw      = c.pw.Read(i)
		shuffle = dsp.Clamp(c.shuffle.Read(i), -0.5, 0.5)
		tempo   = c.tempo.Read(i)
		duty    = math.Floor(60/(tempo*60*dsp.SampleRate)*dsp.SampleRate + 0.5)
	)

	if c.run.Read(i) <= 0 {
		c.out.Write(i, -1)
		return
	}

	if c.even {
		if c.tick >= int(duty+(duty*shuffle)) {
			c.tick = 0
			c.even = false
		}
	} else if c.tick >= int(duty) {
		c.tick = 0
		c.even = true
	}

	if float64(c.tick) <= pw*duty {
		c.out.Write(i, 1)
	} else {
		c.out.Write(i, -1)
	}
	c.tick++
}
