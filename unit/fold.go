package unit

import (
	"buddin.us/shaden/dsp"
)

func newFold(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &fold{
		in:    io.NewIn("in", dsp.Float64(0)),
		level: io.NewIn("level", dsp.Float64(0)),
		gain:  io.NewIn("gain", dsp.Float64(1)),
		out:   io.NewOut("out"),
	}), nil
}

type fold struct {
	in, level, gain, stages *In
	out                     *Out
}

func (f *fold) ProcessSample(i int) {
	var (
		in   = f.in.Read(i)
		lvl  = f.level.Read(i)
		gain = f.gain.Read(i)
		out  = dsp.Fold(in, -lvl, lvl)
	)
	f.out.Write(i, out*gain)
}
