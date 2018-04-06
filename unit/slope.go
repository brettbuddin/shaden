package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newSlope(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &slope{
		state: &slopeState{
			lastTrigger: -1,
		},
		stateFunc: slopeIdle,
		trigger:   io.NewIn("trigger", dsp.Float64(0)),
		gate:      io.NewIn("gate", dsp.Float64(0)),
		rise:      io.NewIn("rise", dsp.Duration(100, c.SampleRate)),
		fall:      io.NewIn("fall", dsp.Duration(100, c.SampleRate)),
		retrigger: io.NewIn("retrigger", dsp.Float64(1)),
		cycle:     io.NewIn("cycle", dsp.Float64(0)),
		ratio:     io.NewIn("ratio", dsp.Float64(0.01)),
		out:       io.NewOut("out"),
		mirror:    io.NewOut("mirror"),
		eoc:       io.NewOut("eoc"),
		eor:       io.NewOut("eor"),
	}), nil
}

type slope struct {
	trigger, retrigger, gate, rise, fall, cycle, ratio *In
	out, mirror, eoc, eor                              *Out
	state                                              *slopeState
	stateFunc                                          slopeStateFunc
}

func (s *slope) ProcessSample(i int) {
	s.state.trigger = s.trigger.Read(i)
	s.state.retrigger = s.retrigger.ReadSlow(i, ident)
	s.state.gate = s.gate.Read(i)
	s.state.rise = math.Abs(s.rise.Read(i))
	s.state.fall = math.Abs(s.fall.Read(i))
	s.state.cycle = s.cycle.Read(i)
	s.state.ratio = s.ratio.Read(i)
	s.stateFunc = s.stateFunc(s.state)
	s.state.lastTrigger = s.state.trigger
	s.state.lastGate = s.state.gate

	s.out.Write(i, s.state.out)
	s.mirror.Write(i, 1-s.state.out)
	s.eoc.Write(i, s.state.eoc)
	s.eor.Write(i, s.state.eor)
}

type slopeStateFunc func(*slopeState) slopeStateFunc

type slopeState struct {
	trigger, retrigger, gate, rise, fall, cycle, ratio float64
	base, multiplier                                   float64
	lastTrigger, lastGate                              float64
	out, eoc, eor                                      float64
}

func slopeIdle(s *slopeState) slopeStateFunc {
	s.out = 0
	s.eoc = -1
	s.eor = -1
	if isTrig(s.lastTrigger, s.trigger) || isTrig(s.lastGate, s.gate) {
		return prepSlopeRise(s)
	}
	return slopeIdle
}

func slopeRise(s *slopeState) slopeStateFunc {
	s.out = s.base + s.out*s.multiplier
	s.eoc = -1
	s.eor = -1
	if s.out >= 1 {
		s.eor = 1
		if s.gate > 0 {
			return slopeHold
		}
		return prepSlopeFall(s)
	}
	return slopeRise
}

func slopeHold(s *slopeState) slopeStateFunc {
	if s.gate <= 0 {
		return prepSlopeFall(s)
	}
	return slopeHold
}

func slopeFall(s *slopeState) slopeStateFunc {
	s.eoc = -1
	s.eor = -1
	if isHigh(s.retrigger) && isTrig(s.lastTrigger, s.trigger) || isTrig(s.lastGate, s.gate) {
		return prepSlopeRise(s)
	}
	s.out = s.base + s.out*s.multiplier
	if s.out < math.SmallestNonzeroFloat64 {
		s.eoc = 1
		s.out = 0
		if s.cycle > 0 {
			return prepSlopeRise(s)
		}
		return slopeIdle
	}
	return slopeFall
}

func prepSlopeRise(s *slopeState) slopeStateFunc {
	s.base, s.multiplier = slopeCoeffs(s.ratio, s.rise, 1, logCurve)
	return slopeRise
}

func prepSlopeFall(s *slopeState) slopeStateFunc {
	s.base, s.multiplier = slopeCoeffs(s.ratio, s.fall, 0, expCurve)
	return slopeFall
}

const (
	expCurve int = iota
	logCurve
)

func slopeCoeffs(ratio, duration, target float64, curve int) (base, multiplier float64) {
	multiplier = dsp.ExpRatio(ratio, duration)
	if curve == expCurve {
		ratio = -ratio
	}
	base = (target + ratio) * (1.0 - multiplier)
	return
}
