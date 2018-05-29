package unit

import "github.com/brettbuddin/shaden/dsp"

const (
	logicOR logicMode = iota
	logicAND
	logicXOR
	logicNOR
	logicNAND
	logicXNOR
)

type logicMode int

func newLogic(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &logic{
		x:    io.NewIn("x", dsp.Float64(0)),
		y:    io.NewIn("y", dsp.Float64(0)),
		mode: io.NewIn("mode", dsp.Float64(logicOR)),
		out:  io.NewOut("out"),
	}), nil
}

type logic struct {
	x, y, mode *In
	out        *Out
}

func (l *logic) ProcessSample(i int) {
	x, y := l.x.Read(i), l.y.Read(i)
	var out float64
	switch logicMode(l.mode.Read(i)) {
	case logicOR:
		out = binaryOR(x, y)
	case logicAND:
		out = binaryAND(x, y)
	case logicXOR:
		out = binaryXOR(x, y)
	case logicNOR:
		out = binaryNOR(x, y)
	case logicNAND:
		out = binaryNAND(x, y)
	case logicXNOR:
		out = binaryXNOR(x, y)
	}
	l.out.Write(i, out)
}
