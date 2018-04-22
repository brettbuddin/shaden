package unit

import (
	"fmt"

	"buddin.us/shaden/dsp"
)

func newRCD(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 8
	}
	var outs []*Out
	for i := 1; i <= config.Size; i++ {
		outs = append(outs, io.NewOut(fmt.Sprintf("%d", i)))
	}
	return NewUnit(io, &rcd{
		clock:  io.NewIn("clock", dsp.Float64(-1)),
		rotate: io.NewIn("rotate", dsp.Float64(-1)),
		reset:  io.NewIn("reset", dsp.Float64(-1)),
		ticks:  make([]int, config.Size),
		outs:   outs,
	}), nil
}

type rcd struct {
	clock, rotate, reset *In
	ticks                []int
	outs                 []*Out

	rotation                         int
	lastClock, lastRotate, lastReset float64
}

func (d *rcd) ProcessSample(i int) {
	var (
		clock  = d.clock.Read(i)
		rotate = d.rotate.Read(i)
		reset  = d.reset.Read(i)
		size   = len(d.outs)
	)

	if isTrig(d.lastRotate, rotate) {
		d.rotation = (d.rotation + 1) % size
	}
	if isTrig(d.lastReset, reset) {
		d.rotation = 0
	}

	clockTrig := isTrig(d.lastClock, clock)

	for j := 0; j < size; j++ {
		count := j - d.rotation
		if count < 0 {
			count = size + count
		}
		if clockTrig {
			d.ticks[j] = (d.ticks[j] + 1) % (count + 1)
		}
		if d.ticks[j] == 0 && count == 0 {
			d.outs[j].Write(i, clock)
		} else if d.ticks[j] == 0 {
			d.outs[j].Write(i, 1)
		} else {
			d.outs[j].Write(i, -1)
		}
	}

	d.lastClock = clock
	d.lastRotate = rotate
	d.lastReset = reset
}
