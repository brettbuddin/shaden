package fft

import "fmt"

// NewConvolution returns a new Convolution
func NewConvolution(fft FFT, window []float64) (*Convolution, error) {
	if len(window) != fft.Size/2 {
		return nil, fmt.Errorf("window size mismatch: %d != %d", len(window), fft.Size/2)
	}

	var (
		input   = make([]complex128, fft.Size)
		process = make([]complex128, fft.Size)
		output  = make([]complex128, fft.Size)
		winf    = make([]complex128, fft.Size)
		wint    = make([]complex128, fft.Size)
	)

	for i := 0; i < fft.Size/2; i++ {
		wint[i] = complex(window[i], 0)
	}
	if err := fft.Transform(winf, wint); err != nil {
		return nil, err
	}

	return &Convolution{
		fft:       fft,
		input:     &input,
		process:   &process,
		output:    &output,
		frequency: make([]complex128, fft.Size),
		window:    winf,
	}, nil
}

// Convolution is a FFT convolution implementation
type Convolution struct {
	fft                    FFT
	input, process, output *[]complex128
	frequency, window      []complex128
}

// Convolve performs convolution using overlap-add method
func (c *Convolution) Convolve(w []float64, r []float64) error {
	size := c.fft.Size
	if len(w) != size/2 {
		return fmt.Errorf("destination size mismatch: %d != %d", len(w), size/2)
	}
	if len(r) != size/2 {
		return fmt.Errorf("source size mismatch: %d != %d", len(r), size/2)
	}

	// Read in the current input and write out the last output we produced
	for i := 0; i < size; i++ {
		if i < size/2 {
			(*c.input)[i] = complex(r[i], 0)
			w[i] = real((*c.output)[i])
		} else {
			(*c.input)[i] = 0
		}
	}

	// Shift the roles of the buffers
	temp := c.process
	c.process = c.input
	c.input = c.output
	c.output = temp

	// FFT -> apply windowed sinc function -> IFFT
	if err := c.fft.Transform(c.frequency, (*c.process)); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		c.frequency[i] *= c.window[i]
	}
	if err := c.fft.Inverse((*c.process), c.frequency); err != nil {
		return err
	}

	// Overlap-add the real part of the numbers
	for i := 0; i < size/2; i++ {
		o := real((*c.output)[i+size/2])
		(*c.process)[i] += complex(o, 0)
	}

	return nil
}
