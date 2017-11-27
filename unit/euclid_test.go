package unit

import (
	"bytes"
	"fmt"
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestEuclidUnit(t *testing.T) {
	builder := Builders()["euclid"]
	u, err := builder(nil)
	require.NoError(t, err)

	clock := u.In["clock"]
	fill := u.In["fill"]
	span := u.In["span"]
	out := u.Out["out"].Out()

	var (
		clockv = -1.0
		gates  []float64
	)
	for i := 0; i < dsp.FrameSize; i++ {
		clock.Write(i, clockv)
		fill.Write(i, 2)
		span.Write(i, 4)
		u.ProcessSample(i)
		if clockv > 0 {
			clockv = -1
		} else {
			clockv = 1
		}
		gates = append(gates, out.Read(i))
	}

	var count int
	for _, v := range gates {
		if v > 0 {
			count++
		}
	}
	require.Equal(t, 64, count)
}

func TestEuclidPatternCreation(t *testing.T) {
	var tests = []struct {
		pattern    string
		span, fill int
	}{
		{"x_x_x", 5, 3},
		{"xxx_", 4, 3},
		{"__x__x__x", 9, 3},
		{"_x_x_x_x_", 9, 4},
		{"x_x_x_x_x", 9, 5},
		{"_________x", 10, 1},
		{"____x____x", 10, 2},
		{"_x_x_x_x_x", 10, 5},
	}

	for _, test := range tests {
		t.Run(test.pattern, func(t *testing.T) {
			var (
				pattern    = make([]bool, 32)
				counts     = make([]int, 32)
				remainders = make([]int, 32)
			)
			euclidean(pattern, counts, remainders, test.span, test.fill)

			buf := bytes.NewBuffer(nil)
			for _, v := range pattern[:test.span] {
				if v {
					fmt.Fprint(buf, "x")
				} else {
					fmt.Fprint(buf, "_")
				}
			}
			require.Equal(t, test.pattern, buf.String())
		})
	}
}
