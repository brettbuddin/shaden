package unit

import (
	"math/rand"

	"buddin.us/lumen/dsp"
)

const maxRandomSeriesSize = 64

func newRandomSeries(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &randomSeries{
		clock:       io.NewIn("clock", dsp.Float64(8)),
		size:        io.NewIn("size", dsp.Float64(8)),
		trigger:     io.NewIn("trigger", dsp.Float64(-1)),
		min:         io.NewIn("min", dsp.Float64(0)),
		max:         io.NewIn("max", dsp.Float64(1)),
		gate:        io.NewOut("gate"),
		value:       io.NewOut("value"),
		valMemory:   make([]float64, maxRandomSeriesSize),
		gateMemory:  make([]float64, maxRandomSeriesSize),
		lastTrigger: -1,
		lastClock:   -1,
	}), nil
}

type randomSeries struct {
	clock, size, trigger, min, max *In
	gate, value                    *Out
	valMemory, gateMemory          []float64

	idx                    int
	lastTrigger, lastClock float64
}

func (s *randomSeries) ProcessSample(i int) {
	size := dsp.Clamp(s.size.Read(i), 1, maxRandomSeriesSize)
	clock := s.clock.Read(i)
	trigger := s.trigger.Read(i)
	min := s.min.Read(i)
	max := s.max.Read(i)

	if s.lastClock < 0 && clock > 0 {
		s.idx++
		if s.idx >= int(size) {
			s.idx = 0
		}
	}
	if s.lastTrigger < 0 && trigger > 0 {
		for j := 0; j < int(size); j++ {
			s.valMemory[j] = dsp.Lerp(min, max, rand.Float64())
			if rand.Float32() > 0.25 {
				s.gateMemory[j] = 1
			} else {
				s.gateMemory[j] = -1
			}
		}
	}
	s.lastTrigger = trigger
	s.lastClock = clock

	s.value.Write(i, s.valMemory[s.idx])
	if s.gateMemory[s.idx] > 0 {
		s.gate.Write(i, clock)
	} else {
		s.gate.Write(i, -1)
	}
}
