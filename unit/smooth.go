package unit

import (
	"buddin.us/shaden/dsp"
)

func newSmooth(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &smooth{
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
	in := s.in.Read(i)
	time := s.time.Read(i)

	if time == 0 {
		s.out.Write(i, in)
		return
	}

	s.average.Window = int(time)
	s.out.Write(i, s.average.Tick(in))
}
