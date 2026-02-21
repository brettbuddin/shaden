package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandRange(t *testing.T) {
	for i := 0; i < 1000; i++ {
		v := RandRange(0, 20)
		require.True(t, v >= 0)
		require.True(t, v <= 20.0)
	}
}

func TestExpRatio(t *testing.T) {
	tests := []struct {
		ratio, speed, expected float64
	}{
		{0.1, 44.1, 0.9470777929275811},
		{0.01, 44.1, 0.9006385575191372},
		{0.001, 44.1, 0.8549937618901365},
		{0.0001, 44.1, 0.8115140950389653},
	}

	for _, test := range tests {
		v := ExpRatio(test.ratio, test.speed)
		require.InEpsilon(t, test.expected, v, 1e-15)
	}
}

func TestSoftClamp(t *testing.T) {
	tests := []struct {
		input, threshold, expected float64
	}{
		{0.5, 1.0, 0.5},
		{0.5, 0.5, 0.5},
		{0.7, 0.5, 0.8214285714285714},
		{0.9, 0.5, 0.8611111111111112},
	}

	for _, test := range tests {
		v := SoftClamp(test.input, test.threshold)
		require.Equal(t, test.expected, v)
	}
}

func TestOverload(t *testing.T) {
	tests := []struct {
		input, expected float64
	}{
		{0.1, 0.09516258196404048},
		{0.5, 0.3934693402873666},
		{1.0, 0.6321205588285577},
		{5.0, 0.9932620530009145},
	}

	for _, test := range tests {
		v := Overload(test.input)
		require.Equal(t, test.expected, v)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		input, min, max, expected float64
	}{
		{0.1, 0.5, 0.9, 0.5},
		{1.0, 0.5, 0.9, 0.9},
		{0.5, 0.5, 0.9, 0.5},
		{0.6, 0.5, 0.9, 0.6},
		{0.9, 0.5, 0.9, 0.9},
	}

	for _, test := range tests {
		v := Clamp(test.input, test.min, test.max)
		require.Equal(t, test.expected, v)
	}
}

func TestMix(t *testing.T) {
	tests := []struct {
		mix, a, b, expected float64
	}{
		{0, 0.5, 0.5, 1.0},
		{0, 0.4, 0.5, 0.9},
		{-1, 0.4, 0.5, 0.4},
		{1, 0.4, 0.5, 0.5},
		{0.5, 0.4, 0.5, 0.7},
	}

	for _, test := range tests {
		v := Mix(test.mix, test.a, test.b)
		require.Equal(t, test.expected, v)
	}
}

func TestPanMix(t *testing.T) {
	tests := []struct {
		pan, a, b, expectedA, expectedB float64
	}{
		{0, 0.5, 0.5, 0.5, 0.5},
		{1, 0.5, 0.5, 0, 0.5},
		{-1, 0.5, 0.5, 0.5, 0},
	}

	for _, test := range tests {
		a, b := PanMix(test.pan, test.a, test.b)
		require.Equal(t, test.expectedA, a)
		require.Equal(t, test.expectedB, b)
	}
}

func TestIsPowerTwo(t *testing.T) {
	require.True(t, IsPowerOfTwo(2))
	require.True(t, IsPowerOfTwo(4))
	require.False(t, IsPowerOfTwo(5))
}

func TestFold(t *testing.T) {
	tests := []struct {
		input, min, max, expected float64
	}{
		{1.0, 0.1, 0.5, 0.19999999999999998},
		{0.6, 0.1, 0.5, 0.4},
		{0.05, 0.1, 0.5, 0.15000000000000002},
	}

	for _, test := range tests {
		v := Fold(test.input, test.min, test.max)
		require.Equal(t, test.expected, v)
	}
}

func TestChebyshev(t *testing.T) {
	tests := []struct {
		order           int
		input, expected float64
	}{
		{0, 1.0, 1.0},
		{1, 2.0, 2.0},
		{2, 2.0, 7.0},
		{3, 1.5, 9.0},
		{3, 3.0, 99.0},
	}

	for _, test := range tests {
		v := Chebyshev(test.order, test.input)
		require.Equal(t, test.expected, v)
	}
}
