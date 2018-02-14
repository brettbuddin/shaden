package generator

import "math"

// WaveType is an alias type for the type of waveforms that can be generated
type WaveType uint16

const (
	WaveSine     WaveType = iota // 0
	WaveTriangle                 // 1
	WaveSaw                      // 2
	WaveSqr                      //3
)

const (
	TwoPi = float64(2 * math.Pi)
)

const (
	SineB = 4.0 / math.Pi
	SineC = -4.0 / (math.Pi * math.Pi)
	Q     = 0.775
	SineP = 0.225
)

// Sine takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Sine(x32 float64) float64 {
	x := float64(x32)
	y := SineB*x + SineC*x*(math.Abs(x))
	y = SineP*(y*(math.Abs(y))-y) + y
	return float64(y)
}

const TringleA = 2.0 / math.Pi

// Triangle takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Triangle(x float64) float64 {
	return float64(TringleA*x) - 1.0
}

// Square takes an input value from -Pi to Pi
// and returns -1 or 1
func Square(x float64) float64 {
	if x >= 0.0 {
		return 1
	}
	return -1.0
}

const SawtoothA = 1.0 / math.Pi

// Triangle takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Sawtooth(x float64) float64 {
	return SawtoothA * x
}
