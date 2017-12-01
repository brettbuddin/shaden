package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

func newTranspose(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &transpose{
		in:        io.NewIn("in", dsp.Float64(0)),
		semitones: io.NewIn("semitones", dsp.Float64(0)),
		out:       io.NewOut("out"),
	}), nil
}

type transpose struct {
	in, semitones *In
	out           *Out
}

func (t *transpose) ProcessSample(i int) {
	var (
		in        = t.in.Read(i)
		semitones = t.semitones.Read(i)
	)
	t.out.Write(i, in*math.Pow(2, semitones/12))
}
