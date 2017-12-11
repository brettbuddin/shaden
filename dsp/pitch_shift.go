package dsp

import "math"

const shiftSize = 2048

// NewPitchShift returns a new PitchShift
func NewPitchShift() *PitchShift {
	return &PitchShift{
		delay: NewDelayLine(shiftSize),
	}
}

// PitchShift shifts a pitch
type PitchShift struct {
	delay *DelayLine
	phase float64
}

// TickSemitones shifts a pitch by a specific number of semitones in equal temperament.
func (s *PitchShift) TickSemitones(in, semitones float64) float64 {
	if semitones == 0 {
		return in
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
	coeff := s.delay.ReadAbsolute(phase) * apply

	midway := phase + size*0.5
	if midway >= size {
		midway -= size
	}
	return coeff + s.delay.ReadAbsolute(midway)*(1.0-apply)
}
