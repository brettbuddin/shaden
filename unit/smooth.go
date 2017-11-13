package unit

import (
	"buddin.us/shaden/dsp"
)

func newSmooth(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &smooth{
		in:      io.NewIn("in", dsp.Float64(0)),
		time:    io.NewIn("time", dsp.Duration(100)),
		out:     io.NewOut("out"),
		average: dsp.RollingAverage{},
	}), nil
}

type smooth struct {
	in, time *In
	out      *Out
	average  dsp.RollingAverage
}

func (s *smooth) ProcessSample(i int) {
	s.average.Window = int(s.time.Read(i))
	s.out.Write(i, s.average.Tick(s.in.Read(i)))
}
