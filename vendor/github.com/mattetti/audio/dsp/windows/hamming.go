package windows

import "math"

// Hamming generates a Hamming window of the requested size
// See https://en.wikipedia.org/wiki/Window_function#Hamming_window
func Hamming(L int) []float64 {
	r := make([]float64, L)
	alpha := 0.54
	beta := 1.0 - alpha
	Lf := float64(L)

	for i := 0; i < L; i++ {
		r[i] = alpha - (beta * math.Cos((twoPi*float64(i))/(Lf-1)))
	}
	return r
}
