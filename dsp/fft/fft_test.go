package fft

import (
	"math"
	"math/cmplx"
	"math/rand"
	"testing"
)

func TestFFT(t *testing.T) {
	const size = 512

	var (
		in   = make([]complex128, size)
		out  = make([]complex128, size)
		freq = make([]complex128, size)
	)

	for i := range in {
		in[i] = complex(rand.Float64(), 0)
	}
	fft, err := New(len(in))
	if err != nil {
		t.Error(err)
	}

	fft.Transform(freq, in)
	fft.Inverse(out, freq)

	if !areRealsSimilar(in, out) {
		t.Errorf("input and output frames are not similar")
	}
}

func areRealsSimilar(a, b []complex128) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		rA, _ := cmplx.Polar(a[i])
		rB, _ := cmplx.Polar(b[i])
		if math.Abs(rB-rA) > 1e-15 {
			return false
		}
	}
	return true
}

func BenchmarkFFT(b *testing.B) {
	const size = 8192

	var (
		in   = make([]complex128, size)
		freq = make([]complex128, size)
	)
	fft, err := New(len(in))
	if err != nil {
		b.Error(err)
	}

	for i := range in {
		in[i] = complex(rand.Float64(), 0)
	}

	for i := 0; i < b.N; i++ {
		fft.Transform(in, freq)
	}
}

func BenchmarkIFFT(b *testing.B) {
	const size = 8192

	var (
		in   = make([]complex128, size)
		out  = make([]complex128, size)
		freq = make([]complex128, size)
	)
	fft, err := New(len(in))
	if err != nil {
		b.Error(err)
	}

	for i := range in {
		in[i] = complex(rand.Float64(), 0)
	}
	if err := fft.Transform(in, freq); err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		fft.Inverse(out, freq)
	}
}
