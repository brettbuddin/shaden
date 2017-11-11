package dsp

// Lerp performs linear interpolation
func Lerp(x1, x2, m float64) float64 {
	return x1 + (x2-x1)*m
}

// Cubic performs cubic interpolation
func Cubic(y0, y1, y2, y3, m float64) float64 {
	m2 := m * m

	a0 := y3 - y2 - y0 + y1
	a1 := y0 - y1 - a0
	a2 := y2 - y0
	a3 := y1

	return a0*m*m2 + a1*m2 + a2*m + a3
}

// Hermite performs Hermite interpolation
func Hermite(x0, x1, x2, x3, t float64) float64 {
	c0 := x1
	c1 := .5 * (x2 - x0)
	c2 := x0 - (2.5 * x1) + (2 * x2) - (.5 * x3)
	c3 := (0.5 * (x3 - x0)) + (1.5 * (x1 - x2))
	return (((((c3 * t) + c2) * t) + c1) * t) + c0
}

// RollingAverage calculates a rolling average over a specified window size
type RollingAverage struct {
	Window int
	value  float64
}

// Tick advances the filter's state
func (a *RollingAverage) Tick(in float64) float64 {
	a.value -= a.value / float64(a.Window)
	a.value += in / float64(a.Window)
	return a.value
}
