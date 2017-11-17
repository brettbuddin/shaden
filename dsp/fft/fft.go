package fft

import (
	"fmt"
	"math"
	"math/cmplx"
)

// New returns a new FFT
func New(size int) (FFT, error) {
	if !isPowerOfTwo(size) {
		return FFT{}, fmt.Errorf("size not a power of two")
	}

	roots := map[int][]complex128{}
	calcRoots(roots, size)

	return FFT{
		Size:    size,
		roots:   roots,
		scratch: make([]complex128, size),
	}, nil
}

// FFT is a Fast-Fourier Transform implementation
type FFT struct {
	Size    int
	roots   map[int][]complex128
	scratch []complex128
}

// Transform performs a forward transform
func (f FFT) Transform(w, r []complex128) error {
	if len(w) != f.Size {
		return fmt.Errorf("destination size mismatch: %d != %d", len(w), f.Size)
	}
	if len(r) != f.Size {
		return fmt.Errorf("source size mismatch: %d != %d", len(r), f.Size)
	}
	f.fft(w, r, f.Size, 1)
	return nil
}

// fft is the recursive element of the FFT
func (f FFT) fft(w, r []complex128, n, step int) {
	if n == 1 {
		w[0] = r[0]
		return
	}

	f.fft(w, r, n/2, 2*step)
	f.fft(w[n/2:], r[step:], n/2, 2*step)

	roots := f.roots[n]
	for k := 0; k < n/2; k++ {
		t := roots[k] * w[k+n/2]
		w[k], w[k+n/2] = w[k]+t, w[k]-t
	}
}

// Inverse performs an inverse transform
func (f FFT) Inverse(w, r []complex128) error {
	if len(w) != f.Size {
		return fmt.Errorf("destination size mismatch: %d != %d", len(w), f.Size)
	}
	if len(r) != f.Size {
		return fmt.Errorf("source size mismatch: %d != %d", len(r), f.Size)
	}
	for i := 0; i < len(r); i++ {
		f.scratch[i] = cmplx.Conj(r[i])
	}
	if err := f.Transform(w, f.scratch); err != nil {
		return err
	}
	for i := 0; i < len(r); i++ {
		w[i] = cmplx.Conj(w[i]) / complex(float64(len(r)), 0)
	}
	return nil
}

// calcRoots pre-calculates coefficients to speed up the FFT
func calcRoots(out map[int][]complex128, size int) {
	if size == 1 {
		return
	}
	out[size] = make([]complex128, int(size/2))
	for k := 0; k < size/2; k++ {
		phase := -2.0 * math.Pi * float64(k) / float64(size)
		out[size][k] = cmplx.Rect(1, phase)
	}
	calcRoots(out, size/2)
}

func isPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}
