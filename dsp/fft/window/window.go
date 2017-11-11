package window

import "math"

// WindowFunc is a windowing function used to prepare frequency domain filters
type WindowFunc func(x float64, n int) float64

// LowPass is a low-pass filter used in the frequency domain
func LowPass(h []float64, wf WindowFunc, cutoff float64) {
	n := len(h)
	for i := range h {
		x := 2 * math.Pi * cutoff
		if i == n/2 {
			h[i] = x
		} else {
			y := float64(i) - float64(n)/2
			h[i] = (math.Sin(x*y) / y) * wf(float64(i), n)
		}
	}
	normalize(h)
}

// HighPass is a high-pass filter used in the frequency domain
func HighPass(h []float64, wf WindowFunc, cutoff float64) {
	LowPass(h, wf, cutoff)
	for i := range h {
		h[i] = -h[i]
	}
}

// BandReject is a band-reject filter used in the frequency domain
func BandReject(h []float64, wf WindowFunc, stop1, stop2 float64) {
	a := make([]float64, len(h))
	b := make([]float64, len(h))
	LowPass(a, wf, stop1)
	HighPass(b, wf, stop2)
	for i := range h {
		h[i] = a[i] + b[i]
	}
}

// BandPass is a band-pass filter used in the frequency domain
func BandPass(h []float64, wf WindowFunc, stop1, stop2 float64) {
	BandReject(h, wf, stop1, stop2)
	for i := range h {
		h[i] = -h[i]
	}
}

func normalize(w []float64) {
	var sum float64
	for i := range w {
		sum += w[i]
	}
	for i := range w {
		w[i] /= sum
	}
}
