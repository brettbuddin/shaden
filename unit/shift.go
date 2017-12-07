package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

const shiftSize = 2048

func newShift(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &shift{
		in:        io.NewIn("in", dsp.Float64(0)),
		semitones: io.NewIn("semitones", dsp.Float64(0)),
		out:       io.NewOut("out"),
		delay:     dsp.NewDelayLine(shiftSize),
	}), nil
}

type shift struct {
	in, semitones *In
	out           *Out
	delay         *dsp.DelayLine
	phase         float64
}

func (s *shift) ProcessSample(i int) {
	in := s.in.Read(i)
	semitones := s.semitones.Read(i)

	if semitones == 0 {
		s.out.Write(i, in)
		return
	}

	ratio := math.Pow(2, semitones/12)
	size := float64(shiftSize)

	s.phase += (1.0 - ratio) / size
	if s.phase >= 1.0 {
		s.phase -= 1.0
	} else if s.phase <= 0.0 {
		s.phase += 1.0
	}
	phase := s.phase * size

	apply := s.phase
	if s.phase >= 0.5 {
		apply = 1.0 - s.phase
	}
	apply *= 2.0

	s.delay.Write(in)
	accum := s.delay.ReadAbsolute(phase) * apply

	midway := phase + size*0.5
	if midway >= size {
		midway -= size
	}
	accum += s.delay.ReadAbsolute(midway) * (1.0 - apply)
	s.out.Write(i, accum)
}
