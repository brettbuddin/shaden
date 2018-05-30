package musictheory

import "math"

func normalizeChromatic(v int) int {
	return int(mod(float64(v), 12))
}

func normalizeDiatonic(v int) int {
	return int(mod(float64(v), 7))
}

func diatonicOctaves(v int) int {
	return v / 7
}

func chromaticOctaves(v int) int {
	return v / 12
}

func inverseChromatic(v int) int {
	return 12 - v
}

func inverseDiatonic(v int) int {
	return 7 - v
}

func mod(n, m float64) float64 {
	out := math.Mod(n, m)
	if out < 0 {
		out += m
	}
	return out
}
