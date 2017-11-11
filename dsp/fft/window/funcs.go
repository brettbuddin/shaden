package window

import "math"

// Blackman is a Blackman window function
func Blackman(x float64, n int) float64 {
	return 0.42 - (0.5 * math.Cos((2*math.Pi*x)/float64(n))) + (0.08 * math.Cos((4*math.Pi*x)/float64(n)))
}

// Hann is a Hann window function
func Hann(x float64, n int) float64 {
	return 0.5 - 0.5*math.Cos(2*x*math.Pi/float64(n))
}

// Lanczos is a Lanczos window function
func Lanczos(x float64, n int) float64 {
	return sinc((2 * x / float64(n)) - 1)
}

func sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Sin(math.Pi*x) / (math.Pi * x)
}
