package analysis

// FFTShiftF shifts a buffer of floats. The passed buffer is modified.
// See http://www.mathworks.com/help/matlab/ref/fftshift.html
func FFTShiftF(buffer []float64) []float64 {
	var tmp float64
	size := len(buffer) / 2
	for i := 0; i < size; i++ {
		tmp = buffer[i]
		buffer[i] = buffer[size+i]
		buffer[size+i] = tmp
	}
	return buffer
}
