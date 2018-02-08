package unit

import (
	"fmt"

	"buddin.us/shaden/dsp"
)

func newGateSeries(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 4
	}

	outs := make([]*Out, config.Size)
	for i := range outs {
		outs[i] = io.NewOut(fmt.Sprintf("%d", i))
	}

	return NewUnit(io, &gateSeries{
		clock:   io.NewIn("clock", dsp.Float64(-1)),
		advance: io.NewIn("advance", dsp.Float64(-1)),
		reset:   io.NewIn("reset", dsp.Float64(-1)),
		outs:    outs,
		target:  -1,
	}), nil
}

type gateSeries struct {
	clock, advance, reset  *In
	outs                   []*Out
	target                 int
	lastAdvance, lastReset float64
}

func (g *gateSeries) ProcessSample(i int) {
	var (
		clk     = g.clock.Read(i)
		reset   = g.reset.Read(i)
		advance float64
	)

	if g.advance.HasSource() {
		advance = g.advance.Read(i)
	} else {
		advance = clk
	}
	if isTrig(g.lastAdvance, advance) {
		g.target = (g.target + 1) % len(g.outs)
	}
	g.lastAdvance = advance

	if isTrig(g.lastReset, reset) {
		g.target = 0
	}
	g.lastReset = reset

	for j, out := range g.outs {
		if j == g.target {
			out.Write(i, clk)
		} else {
			out.Write(i, -1)
		}
	}
}
