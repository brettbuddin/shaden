package unit

import (
	"fmt"

	"buddin.us/lumen/dsp"
)

func newSwitch(name string, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}
	if config.Size == 0 {
		config.Size = 4
	}

	io := NewIO()
	inputs := make([]*In, config.Size)
	for i := range inputs {
		inputs[i] = io.NewIn(fmt.Sprintf("%d", i), dsp.Float64(0))
	}

	s := &seqSwitch{
		trigger:   io.NewIn("trigger", dsp.Float64(-1)),
		reset:     io.NewIn("reset", dsp.Float64(0)),
		inputs:    inputs,
		out:       io.NewOut("out"),
		lastClock: -1,
		lastReset: -1,
	}

	return NewUnit(io, name, s), nil
}

type seqSwitch struct {
	trigger, reset *In
	inputs         []*In
	out            *Out

	step                 int
	lastClock, lastReset float64
}

func (s *seqSwitch) ProcessSample(i int) {
	var (
		trigger = s.trigger.Read(i)
		reset   = s.reset.Read(i)
	)
	if s.lastReset < 0 && reset > 0 {
		s.step = 0
	} else if s.lastClock < 0 && trigger > 0 {
		s.step = (s.step + 1) % len(s.inputs)
	}
	s.lastClock = trigger
	s.out.Write(i, s.inputs[s.step].Read(i))
}
