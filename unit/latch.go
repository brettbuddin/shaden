package unit

import (
	"buddin.us/shaden/dsp"
)

func newLatch(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &latch{
		lastTrigger: -1,
		in:          io.NewIn("in", dsp.Float64(0)),
		trigger:     io.NewIn("trigger", dsp.Float64(0)),
		out:         io.NewOut("out"),
		initial:     false,
	}), nil
}

type latch struct {
	in, trigger         *In
	out                 *Out
	lastTrigger, sample float64
	initial             bool
}

func (l *latch) ProcessSample(i int) {
	if !l.initial {
		l.sample = l.in.Read(i)
		l.initial = true
	}
	in := l.in.Read(i)
	trigger := l.trigger.Read(i)
	if isTrig(l.lastTrigger, trigger) {
		l.sample = in
	}
	l.lastTrigger = trigger
	l.out.Write(i, l.sample)
}
