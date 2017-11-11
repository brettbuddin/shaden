package unit

import (
	"math"

	"buddin.us/lumen/dsp"
)

func newMIDIToHz(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &midiToHz{
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
	m.out.Write(i, 440*math.Pow(2, (in-69)/12))
}
