package unit

import (
	"buddin.us/shaden/dsp"
)

func newCount(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &count{
		trigger:   io.NewIn("trigger", dsp.Float64(-1)),
		reset:     io.NewIn("reset", dsp.Float64(-1)),
		limit:     io.NewIn("limit", dsp.Float64(32)),
		step:      io.NewIn("step", dsp.Float64(1)),
		offset:    io.NewIn("offset", dsp.Float64(0)),
		out:       io.NewOut("out"),
		resetOut:  io.NewOut("reset"),
		lastReset: -1,
	}), nil
}

type count struct {
	trigger, reset, limit, step, offset *In
	out, resetOut                       *Out
	count                               int
	lastTrigger, lastReset              float64
}

func (c *count) ProcessSample(i int) {
	var (
		offset           = int(c.offset.Read(i))
		limit            = int(c.limit.Read(i))
		step             = c.step.Read(i)
		trigger          = c.trigger.Read(i)
		reset            = c.reset.Read(i)
		resetOut float64 = -1
	)

	if isTrig(c.lastReset, reset) {
		c.count = 0
		resetOut = 1
	} else if isTrig(c.lastTrigger, trigger) {
		c.count = (c.count + int(step) + limit) % limit
		if c.count == 0 {
			resetOut = 1
		}
	}

	c.lastReset = reset
	c.lastTrigger = trigger

	c.out.Write(i, float64(offset+c.count))
	c.resetOut.Write(i, resetOut)
}
