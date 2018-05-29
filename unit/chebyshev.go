package unit

import (
	"fmt"

	"github.com/brettbuddin/shaden/dsp"
)

const alphaSeries = "abcdefghijklmnopqrstuvwxyz"

func newChebyshev(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 3
	}

	if config.Size > len(alphaSeries) {
		return nil, fmt.Errorf("maximum size is %d", len(alphaSeries))
	}

	cheb := &chebyshev{
		in:     io.NewIn("in", dsp.Float64(0)),
		coeffs: make([]*In, config.Size),
		out:    io.NewOut("out"),
	}

	for i := range cheb.coeffs {
		cheb.coeffs[i] = io.NewIn(string(alphaSeries[i]), dsp.Float64(0))
	}

	return NewUnit(io, cheb), nil
}

type chebyshev struct {
	in     *In
	coeffs []*In
	out    *Out
}

func (c *chebyshev) ProcessSample(i int) {
	var (
		x   = c.in.Read(i)
		out float64
	)
	for j, v := range c.coeffs {
		out += v.Read(i) * dsp.Chebyshev(j, x)
	}
	c.out.Write(i, out)
}
