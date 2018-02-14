package windows

import "math"

var (
	twoPi  = math.Pi * 2
	fourPi = math.Pi * 4
	sixPi  = math.Pi * 6
)

// Nuttall generates a Blackman-Nutall window
// See https://en.wikipedia.org/wiki/Window_function#Nuttall_window.2C_continuous_first_derivative
func Nuttall(L int) []float64 {
	r := make([]float64, L)
	LF := float64(L)
	for i := 0; i < L; i++ {
		iF := float64(i)
		r[i] = 0.355768 - 0.487396*math.Cos((twoPi*iF)/(LF-1)) + 0.144232*math.Cos((fourPi*iF)/(LF-1)) - 0.012604*math.Cos((sixPi*iF)/(LF-1))
	}

	return r
}