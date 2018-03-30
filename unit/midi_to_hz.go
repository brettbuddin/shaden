package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newMIDIToHz(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &midiToHz{
		in:         io.NewIn("in", dsp.Float64(0)),
		out:        io.NewOut("out"),
		sampleRate: float64(c.SampleRate),
	}), nil
}

type midiToHz struct {
	in         *In
	out        *Out
	sampleRate float64
}

func (m *midiToHz) ProcessSample(i int) {
	in := m.in.Read(i)
	m.out.Write(i, 440*math.Pow(2, (in-69)/12)/m.sampleRate)
}
