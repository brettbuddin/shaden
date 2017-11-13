package unit

import (
	"fmt"

	"buddin.us/shaden/dsp"
)

const alphaSeries = "abcdefghijklmnopqrstuvwxyz"

func newChebyshev(name string, c Config) (*Unit, error) {
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

	io := NewIO()
	cheb := &chebyshev{
		in:     io.NewIn("in", dsp.Float64(0)),
		coeffs: make([]*In, config.Size),
		out:    io.NewOut("out"),
	}

	for i := range cheb.coeffs {
		cheb.coeffs[i] = io.NewIn(string(alphaSeries[i]), dsp.Float64(0))
	}

	return NewUnit(io, name, cheb), nil
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
		out += v.Read(i) * c.n(j, x)
	}
	c.out.Write(i, out)
}

func (c *chebyshev) n(n int, x float64) float64 {
	switch n {
	case 0:
		return 1
	case 1:
		return x
	case 2:
		return (2.0 * x * x) - 1.0
	}
	var (
		y1 = (2.0 * x * x) - 1.0
		y2 = x
		y  = y1
	)
	for i := 3; i <= n; i++ {
		y = (2.0 * x * y1) - y2
		y2, y1 = y1, y
	}
	return y
}
