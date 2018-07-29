package stdout

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdout(t *testing.T) {
	const frameSize = 256

	var (
		r, w   = io.Pipe()
		stdout = New(w, frameSize, 44100)
		msg    = make(chan struct{})
	)

	go func() {
		defer close(msg)
		err := stdout.Start(func(_ []float32, out [][]float32) {
			for i := 0; i < frameSize; i++ {
				out[0][i] = 1
				out[1][i] = 1
			}
		})
		assert.NoError(t, err)
	}()

	var (
		expected = bytes.NewBuffer(nil)
		actual   = make([]byte, 2)
	)

	binary.Write(expected, binary.LittleEndian, int16(math.MaxInt16))

	n, err := r.Read(actual)
	assert.Equal(t, 2, n)
	assert.NoError(t, err)
	assert.Equal(t, expected.Bytes(), actual)

	n, err = r.Read(actual)
	assert.Equal(t, 2, n)
	assert.NoError(t, err)
	assert.Equal(t, expected.Bytes(), actual)

	assert.NoError(t, stdout.Stop())
}
