package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newLag(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &lag{
		in:   io.NewIn("in", dsp.Float64(0)),
		rise: io.NewIn("rise", dsp.Duration(5, c.SampleRate)),
		fall: io.NewIn("fall", dsp.Duration(5, c.SampleRate)),
		out:  io.NewOut("out"),
		slew: newSlew(),
	}), nil
}

type lag struct {
	in, rise, fall *In
	out            *Out
	*slew
}

func (g *lag) ProcessSample(i int) {
	var (
		in   = g.in.Read(i)
		rise = g.rise.Read(i)
		fall = g.fall.Read(i)
	)
	g.out.Write(i, g.slew.Tick(in, rise, fall))
}

type slewStateFunc func(*slewState) slewStateFunc

type slewState struct {
	value, in, lastIn float64
	from, to          float64
	rise, fall        float64
}

type slew struct {
	stateFunc slewStateFunc
	state     *slewState
}

func newSlew() *slew {
	return &slew{slewIdle, &slewState{}}
}

func (s *slew) Tick(v, rise, fall float64) float64 {
	s.state.lastIn, s.state.in = s.state.in, v
	s.state.rise, s.state.fall = rise, fall
	s.stateFunc = s.stateFunc(s.state)
	return s.state.value
}

func slewIdle(s *slewState) slewStateFunc {
	if (s.rise == 0 && s.fall == 0) || (s.lastIn == 0 && s.in != 0) {
		s.value = s.in
		s.lastIn = s.in
		return slewIdle
	}
	if s.in != s.lastIn && math.Abs(s.in-s.lastIn) > math.SmallestNonzeroFloat64 {
		s.from, s.to = s.lastIn, s.in
		s.lastIn = s.in
		s.value = s.from
		return slewTransition
	}
	return slewIdle
}

func slewTransition(s *slewState) slewStateFunc {
	var (
		d      = s.to - s.from
		amount float64
	)
	if d < 0 {
		if s.fall == 0 {
			return slewFinish
		}
		amount = d / s.fall
	} else if d > 0 {
		if s.rise == 0 {
			return slewFinish
		}
		amount = d / s.rise
	} else if math.Abs(d) <= math.SmallestNonzeroFloat64 {
		return slewFinish
	}

	s.value += amount
	remain := s.value - s.to
	if (d > 0 && remain >= 0) || (d < 0 && remain <= 0) {
		return slewFinish
	}
	return slewTransition
}

func slewFinish(s *slewState) slewStateFunc {
	s.value = s.to
	return slewIdle
}
