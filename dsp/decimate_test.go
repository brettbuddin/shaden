package dsp

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecimate_Reduction(t *testing.T) {
	// Bit reduction
	decimate := &Decimate{}
	require.Equal(t, 0.49999990734232586, decimate.Tick(0.5, 44100, 24))
	require.Equal(t, 0.24982677324761315, decimate.Tick(0.5, 44100, 2))

	// Rate reduction
	decimate = &Decimate{}
	rate := 2000.0
	times := 1 / (rate / 44100.0)

	require.Equal(t, 0.0, decimate.Tick(0.5, rate, 24))
	for i := 0; i < int(times)-1; i++ {
		require.Equal(t, 0.0, decimate.Tick(0.4, rate, 24), strconv.Itoa(i))
	}
	require.Equal(t, 0.3999999616118673, decimate.Tick(0.4, rate, 24))
}

func TestDecimate_OutOfRange(t *testing.T) {
	decimate := &Decimate{}
	require.Equal(t, 0.49999990734232586, decimate.Tick(0.5, 80000, 24))
	require.Equal(t, 0.24982677324761315, decimate.Tick(0.5, 80000, 2))
	require.Equal(t, 0.5, decimate.Tick(0.5, 80000, -5))
}
