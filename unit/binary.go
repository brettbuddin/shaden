package unit

import (
	"math"

	"buddin.us/lumen/dsp"
)

func newBinary(name string, op binaryOp) BuildFunc {
	return func(Config) (*Unit, error) {
		io := NewIO()
		return NewUnit(io, name, &binary{
			x:   io.NewIn("x", dsp.Float64(0)),
			y:   io.NewIn("y", dsp.Float64(0)),
			op:  op,
			out: io.NewOut("out"),
		}), nil
	}
}

type binary struct {
	x, y *In
	out  *Out
	op   binaryOp
}

func (b *binary) ProcessSample(i int) {
	b.out.Write(i, b.op(b.x.Read(i), b.y.Read(i)))
}

type binaryOp func(x, y float64) float64

func binarySum(x, y float64) float64  { return x + y }
func binaryDiff(x, y float64) float64 { return x - y }
func binaryMult(x, y float64) float64 { return x * y }
func binaryDiv(x, y float64) float64  { return x / math.Max(y, 1) }
func binaryMod(x, y float64) float64  { return math.Mod(x, math.Max(y, 1)) }
func binaryGT(x, y float64) float64 {
	if x > y {
		return 1
	}
	return -1
}

func binaryLT(x, y float64) float64 {
	if x < y {
		return 1
	}
	return -1
}

func binaryAND(x, y float64) float64 {
	if x > 0 && y > 0 {
		return 1
	}
	return -1
}

func binaryOR(x, y float64) float64 {
	if x > 0 || y > 0 {
		return 1
	}
	return -1
}

func binaryXOR(x, y float64) float64 {
	if (x > 0 && y <= 0) || (y > 0 && x <= 0) {
		return 1
	}
	return -1
}

func binaryNAND(x, y float64) float64 {
	if x > 0 && y > 0 {
		return -1
	}
	return 1
}

func binaryNOR(x, y float64) float64 {
	if x <= 0 && y <= 0 {
		return 1
	}
	return -1
}

func binaryIMPLY(x, y float64) float64 {
	if x > 0 && y <= 0 {
		return -1
	}
	return 1
}

func binaryXNOR(x, y float64) float64 {
	if binaryXOR(x, y) == 1 {
		return -1
	}
	return 1
}

func binaryMax(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func binaryMin(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
