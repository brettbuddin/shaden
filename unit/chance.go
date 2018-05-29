package unit

import (
	"math/rand"

	"github.com/brettbuddin/shaden/dsp"
)

func newChance(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &chance{
		in:   io.NewIn("in", dsp.Float64(0)),
		bias: io.NewIn("bias", dsp.Float64(0)),
		a:    io.NewOut("a"),
		b:    io.NewOut("b"),
	}), nil
}

type chance struct {
	in, bias *In
	a, b     *Out
	last     float64
}

func (c *chance) ProcessSample(i int) {
	var (
		in   = c.in.Read(i)
		bias = dsp.Clamp(c.bias.Read(i), -1, 1)/2 + 0.5
		a, b = -1.0, -1.0
	)

	if isTrig(c.last, in) {
		if bias == 1 {
			a, b = -1, 1
		} else if bias == 0 {
			a, b = 1, -1
		} else {
			if rand.Float64() > bias {
				a, b = 1, -1
			} else {
				a, b = -1, 1
			}
		}
	}

	c.a.Write(i, a)
	c.b.Write(i, b)
	c.last = in
}
