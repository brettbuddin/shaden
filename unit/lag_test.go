package unit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLag(t *testing.T) {
	builder := Builders()["lag"]
	u, err := builder(Config{
		FrameSize:  frameSize,
		SampleRate: sampleRate,
	})
	require.NoError(t, err)

	var (
		in   = u.In["in"]
		rise = u.In["rise"]
		fall = u.In["fall"]
		out  = u.Out["out"].Out()

		samples []float64
	)

	// Impulse that should fade out over 2 samples
	in.Write(0, 1)
	for i := 0; i < frameSize; i++ {
		rise.Write(i, 2)
		fall.Write(i, 2)
		u.ProcessSample(i)
		samples = append(samples, out.Read(i))
	}
	require.Equal(t, []float64{1, 1, 0.5, 0}, samples[:4])
}
