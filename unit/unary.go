package unit

import (
	"math"

	"buddin.us/lumen/dsp"
)

func newUnary(name string, op unaryOp) BuildFunc {
	return func(Config) (*Unit, error) {
		io := NewIO()
		return NewUnit(io, name, &unary{
			x:   io.NewIn("x", dsp.Float64(0)),
			out: io.NewOut("out"),
			op:  op,
		}), nil
	}
}

type unary struct {
	x   *In
	out *Out
	op  unaryOp
}

func (u *unary) ProcessSample(i int) {
	u.out.Write(i, u.op(u.x.Read(i)))
}

type unaryOp func(x float64) float64

func unaryAbs(x float64) float64   { return math.Abs(x) }
func unaryCeil(x float64) float64  { return math.Ceil(x) }
func unaryFloor(x float64) float64 { return math.Floor(x) }
func unaryInv(x float64) float64   { return -x }
func unaryNoop(x float64) float64  { return x }
func unaryNOT(x float64) float64 {
	if x > 0 {
		return -1
	}
	return 1
}
