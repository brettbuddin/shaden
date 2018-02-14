package filters

import (
	"math"

	"github.com/mattetti/audio/dsp/windows"
)

// Sinc represents a sinc function
// The sinc function also called the "sampling function," is a function that
// arises frequently in signal processing and the theory of Fourier transforms.
// The full name of the function is "sine cardinal," but it is commonly referred to by
// its abbreviation, "sinc."
// http://mathworld.wolfram.com/SincFunction.html
type Sinc struct {
	CutOffFreq   float64
	SamplingFreq int
	// Taps are the numbers of samples we go back in time when processing the sync function.
	// The tap numbers will affect the shape of the filter. The more taps, the more
	// shape but the more delays being injected.
	Taps           int
	Window         windows.Function
	_lowPassCoefs  []float64
	_highPassCoefs []float64
}

// LowPassCoefs returns the coeficients to create a low pass filter
func (s *Sinc) LowPassCoefs() []float64 {
	if s == nil {
		return nil
	}
	if s._lowPassCoefs != nil && len(s._lowPassCoefs) > 0 {
		return s._lowPassCoefs
	}
	size := s.Taps + 1
	// sample rate is 2 pi radians per second.
	// we get the cutt off frequency in radians per second
	b := (2 * math.Pi) * s.TransitionFreq()
	s._lowPassCoefs = make([]float64, size)
	// we use a window of size taps + 1
	winData := s.Window(size)

	// we only do half the taps because the coefs are symmetric
	// but we fill up all the coefs
	for i := 0; i < (s.Taps / 2); i++ {
		c := float64(i) - float64(s.Taps)/2
		y := math.Sin(c*b) / (math.Pi * c)
		s._lowPassCoefs[i] = y * winData[i]
		s._lowPassCoefs[size-1-i] = s._lowPassCoefs[i]
	}

	// then we do the ones we missed in case we have an odd number of taps
	s._lowPassCoefs[s.Taps/2] = 2 * s.TransitionFreq() * winData[s.Taps/2]
	return s._lowPassCoefs
}

// HighPassCoefs returns the coeficients to create a high pass filter
func (s *Sinc) HighPassCoefs() []float64 {
	if s == nil {
		return nil
	}
	if s._highPassCoefs != nil && len(s._highPassCoefs) > 0 {
		return s._highPassCoefs
	}

	// we take the low pass coesf and invert them
	size := s.Taps + 1
	s._highPassCoefs = make([]float64, size)
	lowPassCoefs := s.LowPassCoefs()
	winData := s.Window(size)

	for i := 0; i < (s.Taps / 2); i++ {
		s._highPassCoefs[i] = -lowPassCoefs[i]
		s._highPassCoefs[size-1-i] = s._highPassCoefs[i]
	}
	s._highPassCoefs[s.Taps/2] = (1 - 2*s.TransitionFreq()) * winData[s.Taps/2]
	return s._highPassCoefs
}

// TransitionFreq returns a ratio of the cutoff frequency and the sample rate.
func (s *Sinc) TransitionFreq() float64 {
	if s == nil {
		return 0
	}
	return s.CutOffFreq / float64(s.SamplingFreq)
}
