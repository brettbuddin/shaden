package analysis

import (
	"go-dsp/fft"
	"math"
	"math/cmplx"
)

// DFT is the Discrete Fourier Transform representation of a signal
// https://en.wikipedia.org/wiki/Discrete_Fourier_transform
type DFT struct {
	// in audio, we only get real numbers
	// Coefs are the amount of signal energy at those frequency,
	// the amplitude is relative but can be compared as absolute values
	// between buffers.
	Coefs      []complex128
	SampleRate int

	_binWidth int
}

// NewDFT returns the FFT result wrapped in a DFT struct
func NewDFT(sr int, x []float64) *DFT {
	return &DFT{
		SampleRate: sr,
		Coefs:      fft.FFTReal(x),
	}
}

// IFFT runs an inverse fast fourrier transform and returns the float values
func (d *DFT) IFFT() []float64 {
	sndDataCmplx := fft.IFFT(d.Coefs)
	sndData := make([]float64, len(sndDataCmplx))
	for i, cpx := range sndDataCmplx {
		sndData[i] = cmplx.Abs(cpx)
	}
	return sndData
}

// BinWidth is the width of a bin (in frequency).
// It is calculate by using the Nyquist frequency (sample rate/2) divided by
// the DFT size.
func (d *DFT) BinWidth() int {
	if d == nil {
		return 0
	}
	if d._binWidth > 0 {
		return d._binWidth
	}
	d._binWidth = (d.SampleRate / 2) / len(d.Coefs)
	return d._binWidth
}

// ToFreqRange returns a map with the frequency and their values (normalized)
func (d *DFT) ToFreqRange() map[int]float64 {
	if d == nil {
		return nil
	}
	output := make(map[int]float64, len(d.Coefs)/2)
	for i := 0; i < len(d.Coefs)/2; i++ {
		f := (i * d.SampleRate) / (len(d.Coefs))
		// calculate the magnitude
		output[f] = math.Log10(math.Sqrt(math.Pow(real(d.Coefs[i]), 2) + math.Pow(imag(d.Coefs[i]), 2)))
	}
	return output
}
