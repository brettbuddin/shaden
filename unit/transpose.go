package unit

import (
	"math"

	"github.com/brettbuddin/shaden/dsp"
)

func newTranspose(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &transpose{
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
