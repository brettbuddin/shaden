package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newMIDIToHz(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &midiToHz{
		in:  io.NewIn("in", dsp.Float64(0)),
		out: io.NewOut("out"),
	}), nil
}

type midiToHz struct {
	in  *In
	out *Out
}

func (m *midiToHz) ProcessSample(i int) {
	in := m.in.Read(i)
	m.out.Write(i, 440*math.Pow(2, (in-69)/12)/dsp.SampleRate)
}
