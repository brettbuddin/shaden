package windows

import "math"

// Blackman generates a Blackman window of the requested size
// See https://en.wikipedia.org/wiki/Window_function#Blackman_windows
func Blackman(L int) []float64 {
	r := make([]float64, L)
	LF := float64(L)
	alpha := 0.16
	a0 := (1 - alpha) / 2.0
	a1 := 0.5
	a2 := alpha / 2.0

	for i := 0; i < L; i++ {
		iF := float64(i)
		r[i] = a0 - (a1 * math.Cos((twoPi*iF)/(LF-1))) + (a2 * math.Cos((fourPi*iF)/(LF-1)))
	}
	return r
}
