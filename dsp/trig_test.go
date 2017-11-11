package dsp

import (
	"math"
	"testing"
)

const tolerance = 0.00001

func BenchmarkSinTableLookup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sin(0.3)
	}
}

func BenchmarkSinStdLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		math.Sin(0.3)
	}
}

func TestSin(t *testing.T) {
	tests := []float64{
		0.5,
		-0.5,
		math.Pi,
		-math.Pi,
		10 * math.Pi,
		-10 * math.Pi,
	}

	for _, f := range tests {
		expected, actual := math.Sin(f), Sin(f)
		diff := math.Abs(expected - float64(actual))
		if diff > tolerance {
			t.Errorf("%v not within %v tolerance: input=%v expected=%v actual=%v", diff, tolerance, f, expected, actual)
		}
	}
}

func TestTan(t *testing.T) {
	tests := []float64{
		0.5,
		-0.5,
		math.Pi,
		-math.Pi,
		10 * math.Pi,
		-10 * math.Pi,
	}

	for _, f := range tests {
		expected, actual := math.Tan(f), Tan(f)
		diff := math.Abs(expected - float64(actual))
		if diff > tolerance {
			t.Errorf("%v not within %v tolerance: input=%v expected=%v actual=%v", diff, tolerance, f, expected, actual)
		}
	}
}
