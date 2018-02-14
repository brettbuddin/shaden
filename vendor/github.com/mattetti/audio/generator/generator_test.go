package generator

import (
	"math"
	"testing"
)

func TestSine(t *testing.T) {
	testCases := []struct {
		in  float64
		out float64
	}{
		0: {float64(-math.Pi), 0},
		1: {0.007, 0.006909727339533104},
		2: {-0.5, -0.47932893655759223},
		3: {0.1, 0.09895415534087945},
		4: {1.5862234, 0.9998818440160414},
		5: {2.0, 0.909795856141705},
		6: {3.0, 0.14008939955174454},
		7: {math.Pi, 0},
	}

	for i, tc := range testCases {
		if out := Sine(tc.in); !nearlyEqual(out, tc.out, 0.0001) {
			t.Logf("[%d] sine(%f) => %.7f != %.7f", i, tc.in, out, tc.out)
			t.Fail()
		}
	}
}

func nearlyEqual(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	absA := math.Abs(float64(a))
	absB := math.Abs(float64(b))
	diff := math.Abs(float64(a) - float64(b))

	if a == 0 || b == 0 || diff < 0.0000001 {
		return diff < (float64(epsilon))
	}
	return diff/(absA+absB) < float64(epsilon)

}
