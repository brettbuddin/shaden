package unit

import (
	"buddin.us/shaden/dsp"
)

func newShift(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &shift{
		in:        io.NewIn("in", dsp.Float64(0)),
		semitones: io.NewIn("semitones", dsp.Float64(0)),
		out:       io.NewOut("out"),
		shift:     dsp.NewPitchShift(),
	}), nil
}

type shift struct {
	in, semitones *In
	out           *Out
	shift         *dsp.PitchShift
}

func (s *shift) ProcessSample(i int) {
	in := s.in.Read(i)
	semitones := s.semitones.Read(i)
	s.out.Write(i, s.shift.TickSemitones(in, semitones))
}
