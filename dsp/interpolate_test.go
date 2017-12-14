package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCubic(t *testing.T) {
	require.Equal(t, 1.0, Cubic(0.0, 1.0, 3.0, 4.0, 0.0))
	require.Equal(t, 1.272, Cubic(0.0, 1.0, 3.0, 4.0, 0.1))
	require.Equal(t, 1.496, Cubic(0.0, 1.0, 3.0, 4.0, 0.2))
	require.Equal(t, 2.728, Cubic(0.0, 1.0, 3.0, 4.0, 0.9))
	require.Equal(t, 3.0, Cubic(0.0, 1.0, 3.0, 4.0, 1.0))
}

func TestHermite(t *testing.T) {
	require.Equal(t, 1.0, Hermite(0.0, 1.0, 3.0, 4.0, 0.0))
	require.Equal(t, 1.164, Hermite(0.0, 1.0, 3.0, 4.0, 0.1))
	require.Equal(t, 1.352, Hermite(0.0, 1.0, 3.0, 4.0, 0.2))
	require.Equal(t, 2.8360000000000003, Hermite(0.0, 1.0, 3.0, 4.0, 0.9))
	require.Equal(t, 3.0, Hermite(0.0, 1.0, 3.0, 4.0, 1.0))
}

func TestRollingAverage(t *testing.T) {
	avg := &RollingAverage{Window: 3}

	require.Equal(t, 0.0, avg.Tick(0.0))
	require.Equal(t, 0.3333333333333333, avg.Tick(1.0))
	require.Equal(t, 0.5555555555555556, avg.Tick(1.0))
	require.Equal(t, 0.7037037037037037, avg.Tick(1.0))
	require.Equal(t, 0.8024691358024691, avg.Tick(1.0))
	require.Equal(t, 0.5349794238683128, avg.Tick(0.0))
}
