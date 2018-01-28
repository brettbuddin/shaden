package unit

import (
	"fmt"
	"math"
	"math/rand"

	"buddin.us/shaden/dsp"
)

const (
	pulseModeRest int = iota
	pulseModeFirst
	pulseModeLast
	pulseModeAll
	pulseModeHold
)

const (
	patternModeForward int = iota
	patternModeReverse
	patternModePingPong
	patternModeRandom
)

func newStages(name string, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 5
	}

	var (
		io          = NewIO()
		stageInputs = make([]*pulseSequencerStage, config.Size)
	)
	for i := range stageInputs {
		stageInputs[i] = &pulseSequencerStage{
			freq:   io.NewIn(fmt.Sprintf("%d/freq", i), dsp.Float64(0)),
			pulses: io.NewIn(fmt.Sprintf("%d/pulses", i), dsp.Float64(1)),
			mode:   io.NewIn(fmt.Sprintf("%d/mode", i), dsp.Float64(pulseModeFirst)),
			glide:  io.NewIn(fmt.Sprintf("%d/glide", i), dsp.Float64(0)),
			data:   io.NewIn(fmt.Sprintf("%d/data", i), dsp.Float64(0)),
		}
	}

	return NewUnit(io, name, &pulseSequencer{
		clock:       io.NewIn("clock", dsp.Float64(-1)),
		mode:        io.NewIn("mode", dsp.Float64(patternModeForward)),
		reset:       io.NewIn("reset", dsp.Float64(-1)),
		totalStages: io.NewIn("stages", dsp.Float64(config.Size)),
		glidetime:   io.NewIn("glide-time", dsp.Float64(0)),
		out:         io.NewOut("freq"),
		gate:        io.NewOut("gate"),
		data:        io.NewOut("data"),
		eos:         io.NewOut("eos"),
		slew:        newSlew(),
		stageInputs: stageInputs,
		pulse:       -1,
		lastClock:   -1,
		lastReset:   -1,
	}), nil
}

type pulseSequencerStage struct {
	freq, pulses, mode, glide, data *In
	values                          pulseSequencerValues
}

type pulseSequencerValues struct {
	freq, pulses, mode, glide, data float64
}

type pulseSequencer struct {
	clock, reset, mode, glidetime, totalStages *In
	stageInputs                                []*pulseSequencerStage
	out, gate, data, eos                       *Out

	slew         *slew
	pong         bool
	stage, pulse int

	stageOnset           bool
	lastStage            int
	lastClock, lastReset float64
}

func (s *pulseSequencer) ProcessSample(i int) {
	for j, stg := range s.stageInputs {
		s.stageInputs[j].values.freq = stg.freq.ReadSlow(i, ident)
		s.stageInputs[j].values.pulses = stg.pulses.ReadSlow(i, minZero)
		s.stageInputs[j].values.mode = stg.mode.ReadSlow(i, ident)
		s.stageInputs[j].values.glide = stg.glide.ReadSlow(i, ident)
		s.stageInputs[j].values.data = stg.data.Read(i)
	}

	var (
		actualStageCount = float64(len(s.stageInputs))
		totalStages      = int(math.Max(math.Min(actualStageCount, s.totalStages.Read(i)), 1))
		glideTime        = s.glidetime.ReadSlow(i, minZero)
		clock            = s.clock.Read(i)
		reset            = s.reset.Read(i)
		mode             = int(s.mode.Read(i))
	)

	if isTrig(s.lastClock, clock) {
		s.advance(totalStages, mode)
	} else if isTrig(s.lastReset, reset) {
		s.stage = 0
		s.pulse = 0
	}

	s.fillGate(i, clock)
	s.fillFreq(i, glideTime)
	s.fillData(i)
	s.fillEOS(i)

	s.lastClock = clock
	s.lastReset = reset
	s.lastStage = s.stage
}

func (s *pulseSequencer) advance(totalStages, mode int) {
	pulses := int(s.stageInputs[s.stage].values.pulses)

	if s.pulse < 0 {
		s.pulse++
		return
	}

	s.pulse = (s.pulse + 1) % pulses
	if s.lastStage < 0 || s.pulse != 0 {
		return // Keep counting pulses
	}

	s.pulse = 0

	switch mode {
	case patternModeForward:
		s.stage = (s.stage + 1) % totalStages
		s.pong = false
	case patternModeReverse:
		s.stage -= 1
		if s.stage < 0 {
			s.stage = totalStages - 1
		}
		s.pong = false
	case patternModePingPong:
		var inc = 1
		if s.pong {
			inc = -1
		}
		s.stage += inc

		if s.stage > totalStages-1 {
			s.stage = totalStages - 1
			s.pong = true
		} else if s.stage < 0 {
			s.stage = 0
			s.pong = false
		}
	case patternModeRandom:
		s.stage = rand.Intn(totalStages)
		s.pong = false
	}
}

func (s *pulseSequencer) fillGate(i int, clock float64) {
	var (
		stage     = s.stageInputs[s.stage]
		mode      = int(stage.values.mode)
		lastPulse = int(stage.values.pulses) - 1
	)

	if s.lastStage != s.stage {
		s.gate.Write(i, -1)
		s.stageOnset = true
	} else {
		switch mode {
		case pulseModeHold:
			s.gate.Write(i, 1)
		case pulseModeAll:
			if s.stageOnset || isHigh(clock) {
				s.gate.Write(i, 1)
			} else {
				s.gate.Write(i, -1)
			}
		case pulseModeFirst:
			if s.stageOnset || (s.pulse == 0 && isHigh(clock)) {
				s.gate.Write(i, 1)
			} else {
				s.gate.Write(i, -1)
			}
		case pulseModeLast:
			if s.stageOnset || (s.pulse == lastPulse && isHigh(clock)) {
				s.gate.Write(i, 1)
			} else {
				s.gate.Write(i, -1)
			}
		case pulseModeRest:
			s.gate.Write(i, -1)
		}

		s.stageOnset = false
	}
}

func (s *pulseSequencer) fillFreq(i int, glidetime float64) {
	var (
		stage = s.stageInputs[s.stage]
		freq  = stage.values.freq
		glide = stage.values.glide
	)
	if glide == 0 {
		glidetime = 0
	}
	s.out.Write(i, s.slew.Tick(freq, glidetime, glidetime))
}

func (s *pulseSequencer) fillData(i int) {
	s.data.Write(i, s.stageInputs[s.stage].values.data)
}

func (s *pulseSequencer) fillEOS(i int) {
	if s.lastStage != s.stage {
		s.eos.Write(i, 1)
	} else {
		s.eos.Write(i, -1)
	}
}
