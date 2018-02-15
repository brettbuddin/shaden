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

func newStages(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 5
	}

	stageInputs := make([]*stage, config.Size)
	for i := range stageInputs {
		stageInputs[i] = &stage{
			freq:   io.NewIn(fmt.Sprintf("%d/freq", i), dsp.Float64(0)),
			pulses: io.NewIn(fmt.Sprintf("%d/pulses", i), dsp.Float64(1)),
			mode:   io.NewIn(fmt.Sprintf("%d/mode", i), dsp.Float64(pulseModeFirst)),
			glide:  io.NewIn(fmt.Sprintf("%d/glide", i), dsp.Float64(0)),
			data:   io.NewIn(fmt.Sprintf("%d/data", i), dsp.Float64(0)),
		}
	}

	return NewUnit(io, &stages{
		clock:       io.NewIn("clock", dsp.Float64(-1)),
		mode:        io.NewIn("mode", dsp.Float64(patternModeForward)),
		reset:       io.NewIn("reset", dsp.Float64(-1)),
		totalStages: io.NewIn("stages", dsp.Float64(config.Size)),
		glideTime:   io.NewIn("glide-time", dsp.Float64(0)),
		out:         io.NewOut("freq"),
		gate:        io.NewOut("gate"),
		data:        io.NewOut("data"),
		eos:         io.NewOut("eos"),
		slew:        newSlew(),
		stageInputs: stageInputs,
		pulse:       -1,
		lastStage:   -1,
		lastClock:   -1,
		lastReset:   -1,
	}), nil
}

type stage struct {
	freq, pulses, mode, glide, data *In
	values                          stageValues
}

func (s *stage) read(i int) {
	s.values.data = s.data.Read(i)
	s.values.freq = s.freq.Read(i)
	s.values.glide = s.glide.Read(i)
	s.values.mode = int(s.mode.Read(i))
	s.values.pulses = int(math.Max(s.pulses.Read(i), 0))
}

type stageValues struct {
	freq, glide, data float64
	mode, pulses      int
}

type stages struct {
	clock, reset, mode, glideTime, totalStages *In
	stageInputs                                []*stage
	out, gate, data, eos                       *Out

	slew                    *slew
	pong, firstPulse        bool
	stage, pulse, lastStage int
	lastClock, lastReset    float64
}

func (s *stages) ProcessSample(i int) {
	var (
		clock            = s.clock.Read(i)
		reset            = s.reset.Read(i)
		glideTime        = s.glideTime.Read(i)
		mode             = s.mode.ReadSlowInt(i, identInt)
		actualStageCount = float64(len(s.stageInputs))
		totalStages      = s.totalStages.ReadSlowInt(i, clampInt(1, actualStageCount))
	)

	stage := s.stageInputs[s.stage]
	stage.read(i)

	if isTrig(s.lastClock, clock) {
		s.advance(stage, totalStages, mode)
	} else if isTrig(s.lastReset, reset) {
		s.stage = 0
		s.pulse = 0
	}

	if s.lastStage != s.stage {
		stage = s.stageInputs[s.stage]
		stage.read(i)
	}

	s.fillGate(i, stage, clock)
	s.fillFreq(i, stage, glideTime)
	s.fillData(i, stage)
	s.fillEOS(i)

	s.lastClock = clock
	s.lastReset = reset
	s.lastStage = s.stage
}

func (s *stages) advance(stage *stage, totalStages, mode int) {
	pulses := stage.values.pulses

	if pulses != 0 {
		if s.pulse < 0 {
			s.pulse++
			return // we just started the sequencer
		}
		s.pulse = (s.pulse + 1) % pulses
		if s.lastStage < 0 || s.pulse != 0 {
			return // keep counting pulses
		}
	}

	s.advanceStage(totalStages, mode)
}

func (s *stages) advanceStage(totalStages, mode int) {
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

		if s.stage >= totalStages-1 {
			inc = -1
			s.pong = true
		} else if s.stage <= 0 {
			inc = 1
			s.pong = false
		}

		s.stage += inc
		if s.stage < 0 {
			s.stage = 0
		}
	case patternModeRandom:
		s.stage = rand.Intn(totalStages)
		s.pong = false
	}
}

func (s *stages) fillGate(i int, stage *stage, clock float64) {
	if s.lastStage != s.stage {
		s.gate.Write(i, -1)
		s.firstPulse = true
	} else {
		switch stage.values.mode {
		case pulseModeHold:
			s.gate.Write(i, 1)
		case pulseModeAll:
			if s.firstPulse || isHigh(clock) {
				s.gate.Write(i, 1)
			} else {
				s.gate.Write(i, -1)
			}
		case pulseModeFirst:
			if s.firstPulse || (s.pulse == 0 && isHigh(clock)) {
				s.gate.Write(i, 1)
			} else {
				s.gate.Write(i, -1)
			}
		case pulseModeLast:
			lastPulse := stage.values.pulses - 1
			if s.pulse == lastPulse && isHigh(clock) {
				s.gate.Write(i, 1)
			} else {
				s.gate.Write(i, -1)
			}
		case pulseModeRest:
			s.gate.Write(i, -1)
		}

		s.firstPulse = false
	}
}

func (s *stages) fillFreq(i int, stage *stage, glidetime float64) {
	var (
		freq  = stage.values.freq
		glide = stage.values.glide
	)
	if glide == 0 {
		glidetime = 0
	}
	s.out.Write(i, s.slew.Tick(freq, glidetime, glidetime))
}

func (s *stages) fillData(i int, stage *stage) {
	s.data.Write(i, stage.values.data)
}

func (s *stages) fillEOS(i int) {
	if s.lastStage != s.stage {
		s.eos.Write(i, 1)
	} else {
		s.eos.Write(i, -1)
	}
}
