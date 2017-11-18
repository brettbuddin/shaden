package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllPass_Tick(t *testing.T) {
	ap := NewAllPass(10)
	out := []float64{}
	for i := 0; i < 10; i++ {
		out = append(out, ap.Tick(float64(i), 0.5))
	}
	require.Equal(t, []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5}, out)
}

func TestAllPass_TickAbsolute(t *testing.T) {
	ap := NewAllPass(10)
	out := []float64{}
	for i := 0; i < 10; i++ {
		out = append(out, ap.TickAbsolute(float64(i), 0.5, 1))
	}
	require.Equal(t, []float64{
		0,
		1.5,
		2.25,
		3.375,
		4.3125,
		5.34375,
		6.328125,
		7.3359375,
		8.33203125,
		9.333984375,
	}, out)
}

func TestAllPass_TickRelative(t *testing.T) {
	ap := NewAllPass(10)
	out := []float64{}
	for i := 0; i < 10; i++ {
		out = append(out, ap.TickRelative(float64(i), 0.5, 0.5))
	}
	require.Equal(t, []float64{0, 0.5, 1, 1.5, 2, 3.5, 4.75, 6, 7.25, 8.5}, out)
}
