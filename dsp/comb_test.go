package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFBComb_Tick(t *testing.T) {
	c := NewFBComb(10)
	for i := 0; i < 10; i++ {
		c.Tick(float64(i), 0.5)
	}
	out := []float64{}
	for i := 0; i < 10; i++ {
		out = append(out, c.Tick(float64(i), 0.5))
	}
	require.Equal(t, []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 10.5, 12, 13.5}, out)
}

func TestFBComb_TickAbsolute(t *testing.T) {
	c := NewFBComb(10)
	out := []float64{}
	for i := 0; i < 10; i++ {
		out = append(out, c.TickAbsolute(float64(i), 0.5, 1))
	}
	require.Equal(t, []float64{0, 1, 2.5, 4, 5.5, 7, 8.5, 10, 11.5, 13}, out)
}

func TestFFComb_Tick(t *testing.T) {
	c := NewFFComb(10)
	for i := 0; i < 10; i++ {
		c.Tick(float64(i), 0.5)
	}
	out := []float64{}
	for i := 0; i < 10; i++ {
		out = append(out, c.Tick(float64(i), 0.5))
	}
	require.Equal(t, []float64{1, 2.5, 4, 5.5, 7, 8.5, 10, 11.5, 8, 9.5}, out)
}

func TestFFComb_TickAbsolute(t *testing.T) {
	c := NewFFComb(10)
	out := []float64{}
	for i := 0; i < 10; i++ {
		out = append(out, c.TickAbsolute(float64(i), 0.5, 1))
	}
	require.Equal(t, []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 10.5, 12, 13.5}, out)
}
