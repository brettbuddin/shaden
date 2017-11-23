package dsp

import (
	"math"
	"math/rand"
)

// RandRange returns random values between a specified range
func RandRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// ExpRatio produces an (inverse-)exponential curve that's inflection can be controlled by a specific ratio
func ExpRatio(ratio, speed float64) float64 {
	return math.Exp(-math.Log(float64((1+ratio)/ratio)) / float64(speed))
}

// SoftClamp limits a value to a specific range, but compresses the value as it goes beyond the threshold
func SoftClamp(s, threshold float64) float64 {
	abs := math.Abs(s)
	if abs <= 0.5 {
		return s
	}
	return (abs - 0.25*(1-threshold)) / s
}

// Overload is a sigmoid function that simulates soft clip overloading
func Overload(x float64) float64 {
	return math.Copysign(1, x) * (1 - math.Exp(-math.Abs(x)))
}

// Clamp limits a value to a specific range
func Clamp(s, min, max float64) float64 {
	if s > max {
		return max
	} else if s < min {
		return min
	}
	return s
}

// Mix sums two panned inputs
func Mix(mix, a, b float64) float64 {
	aOut, bOut := PanMix(mix, a, b)
	return aOut + bOut
}

// PanMix pans two inputs between two outputs
func PanMix(pan, a, b float64) (float64, float64) {
	pan = Clamp(pan, -1, 1)
	if pan > 0 {
		return (1 - pan) * a, b
	} else if pan < 0 {
		return a, (1 + pan) * b
	}
	return a, b
}

// IsPowerOfTwo determines whether or not an integer is a power of two
func IsPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}

// Fold reflects a value exceeding minimum/maximum thresholds back over those thresholds
func Fold(s, min, max float64) float64 {
	bottomdiff := s - min

	if s >= max {
		s = max + max - s
		if s >= min {
			return s
		}
	} else if s < min {
		s = min + min - s
		if s < max {
			return s
		}
	} else {
		return s
	}

	if max == min {
		return min
	}

	diff := max - min
	diff2 := diff + diff
	s = bottomdiff - diff2*math.Floor(bottomdiff/diff2)
	if s >= diff {
		s += -diff2
	}
	return s + min
}
