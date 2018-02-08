package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newAdsr(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &adsr{
		state:       &adsrState{},
		stateFunc:   adsrIdle,
		gate:        io.NewIn("gate", dsp.Float64(0)),
		attack:      io.NewIn("attack", dsp.Duration(50)),
		decay:       io.NewIn("decay", dsp.Duration(50)),
		sustain:     io.NewIn("sustain", dsp.Float64(0.5)),
		sustainHold: io.NewIn("sustain-hold", dsp.Duration(0)),
		release:     io.NewIn("release", dsp.Duration(50)),
		cycle:       io.NewIn("cycle", dsp.Float64(0)),
		ratio:       io.NewIn("ratio", dsp.Float64(0.01)),
		out:         io.NewOut("out"),
		mirror:      io.NewOut("mirror"),
		eoc:         io.NewOut("eoc"),
	}), nil
}

type adsr struct {
	gate, cycle, ratio                           *In
	attack, decay, sustain, sustainHold, release *In
	out, mirror, eoc                             *Out
	state                                        *adsrState
	stateFunc                                    adsrStateFunc
}

func (s *adsr) ProcessSample(i int) {
	s.state.gate = s.gate.Read(i)
	s.state.attack = s.attack.Read(i)
	s.state.decay = s.decay.Read(i)
	s.state.sustain = s.sustain.Read(i)
	s.state.sustainHold = s.sustainHold.Read(i)
	s.state.release = s.release.Read(i)
	s.state.cycle = s.cycle.Read(i)
	s.state.ratio = s.ratio.Read(i)
	s.stateFunc = s.stateFunc(s.state)
	s.state.lastGate = s.state.gate

	s.out.Write(i, s.state.out)
	s.mirror.Write(i, 1-s.state.out)
	s.eoc.Write(i, s.state.eoc)
}

type adsrStateFunc func(*adsrState) adsrStateFunc

type adsrState struct {
	gate, cycle, ratio                           float64
	attack, decay, sustain, sustainHold, release float64
	base, multiplier, sustainDur                 float64
	lastGate                                     float64
	out, eoc                                     float64
}

func adsrIdle(s *adsrState) adsrStateFunc {
	s.out = 0
	s.eoc = -1
	if isTrig(s.lastGate, s.gate) {
		return prepAdsrAttack(s)
	}
	return adsrIdle
}

func adsrAttack(s *adsrState) adsrStateFunc {
	s.out = s.base + s.out*s.multiplier
	s.eoc = -1
	if s.out >= 1 {
		return prepAdsrDecay(s)
	}
	return adsrAttack
}

func adsrDecay(s *adsrState) adsrStateFunc {
	s.out = s.base + s.out*s.multiplier
	if s.out <= s.sustain {
		if s.gate > 0 {
			return adsrHold
		}
		if s.sustainHold > 0 {
			return prepAdsrSustain(s)
		}
		return prepAdsrRelease(s)
	}
	return adsrDecay
}

func adsrHold(s *adsrState) adsrStateFunc {
	if s.gate <= 0 {
		return prepAdsrRelease(s)
	}
	return adsrHold
}

func adsrSustain(s *adsrState) adsrStateFunc {
	s.sustainDur++
	if s.sustainDur >= s.sustainHold {
		return prepAdsrRelease(s)
	}
	return adsrSustain
}

func adsrRelease(s *adsrState) adsrStateFunc {
	if isTrig(s.lastGate, s.gate) {
		return prepAdsrAttack(s)
	}
	s.out = s.base + s.out*s.multiplier
	if s.out < math.SmallestNonzeroFloat64 {
		s.eoc = 1
		s.out = 0
		if s.cycle > 0 {
			return prepAdsrAttack(s)
		}
		return adsrIdle
	}
	return adsrRelease
}

func prepAdsrAttack(s *adsrState) adsrStateFunc {
	s.base, s.multiplier = slopeCoeffs(s.ratio, s.attack, 1, logCurve)
	return adsrAttack
}

func prepAdsrDecay(s *adsrState) adsrStateFunc {
	s.base, s.multiplier = slopeCoeffs(s.ratio, s.decay, s.sustain, expCurve)
	return adsrDecay
}

func prepAdsrSustain(s *adsrState) adsrStateFunc {
	s.sustainDur = 0
	return adsrSustain
}

func prepAdsrRelease(s *adsrState) adsrStateFunc {
	s.base, s.multiplier = slopeCoeffs(s.ratio, s.release, 0, expCurve)
	return adsrRelease
}
