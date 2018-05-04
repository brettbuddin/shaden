package dsp

import "math"

const (
	sineLength = 1024
	sineStep   = sineLength / (2 * math.Pi)
)

var (
	sineTable = [sineLength]float64{}
	sineDiff  = [sineLength]float64{}
)

func init() {
	for i := 0; i < sineLength; i++ {
		sineTable[i] = math.Sin(float64(i) * (1 / sineStep))
	}
	for i := 0; i < sineLength; i++ {
		next := sineTable[(i+1)%sineLength]
		sineDiff[i] = next - sineTable[i]
	}
}

// Sin is a lookup table version of math.Sin
func Sin(x float64) float64 {
	step := x * sineStep
	if x < 0 {
		step = -step
	}

	var (
		trunc = int(step)
		i     = trunc % sineLength
		out   = sineTable[i] + sineDiff[i]*(step-float64(trunc))
	)

	if x < 0 {
		return -out
	}
	return out
}

// Tan is a lookup table version of math.Tan
func Tan(x float64) float64 {
	return Sin(x) / Cos(x)
}

// Cos is a lookup table version of math.Cos
func Cos(x float64) float64 {
	return Sin(x + 0.5*math.Pi)
}
