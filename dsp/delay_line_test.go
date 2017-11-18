package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDelayLine_Tick(t *testing.T) {
	dl := NewDelayLine(10)
	for i := 0; i < 10; i++ {
		dl.TickAbsolute(float64(i), 1)
	}
	require.Equal(t, []float64{0, 9, 8, 7, 6, 5, 4, 3, 2, 1}, dl.buffer)
	require.Equal(t, 10, dl.Size())
}

func TestDelayLine_ReadAbsoluteIntepolation(t *testing.T) {
	dl := NewDelayLine(10)
	for i := 0; i < 10; i++ {
		dl.TickAbsolute(float64(i), 1)
	}
	// 0, 9, 8, 7, 6, 5, 4, 3, 2, 1
	require.Equal(t, 0.0, dl.ReadAbsolute(0))
	require.Equal(t, 4.5, dl.ReadAbsolute(0.5))
	require.Equal(t, 9.0, dl.ReadAbsolute(1))
	require.Equal(t, 8.0, dl.ReadAbsolute(2))
	require.Equal(t, 7.5, dl.ReadAbsolute(2.5))
}

func TestDelayLine_ReadRelativeIntepolation(t *testing.T) {
	dl := NewDelayLine(10)
	for i := 0; i < 10; i++ {
		dl.TickAbsolute(float64(i), 1)
	}
	// 0, 9, 8, 7, 6, 5, 4, 3, 2, 1
	require.Equal(t, 0.0, dl.ReadRelative(0))
	require.Equal(t, 8.1, dl.ReadRelative(0.1))
	require.Equal(t, 5.5, dl.ReadRelative(0.5))
	require.Equal(t, 3.7, dl.ReadRelative(0.7))
}
