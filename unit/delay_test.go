package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestDelay(t *testing.T) {
	builder := Builders()["delay"]
	u, err := builder(nil)
	require.NoError(t, err)

	var (
		in     = u.In["in"]
		mix    = u.In["mix"]
		fbgain = u.In["fb-gain"]
		time   = u.In["time"]
		out    = u.Out["out"].Out()
	)

	// Mix of dry and wet
	in.Write(0, 1)
	u.ProcessSample(0)
	require.Equal(t, 1.0, out.Read(0))

	// Only wet signal
	mix.Write(0, 1)
	in.Write(0, 1)
	u.ProcessSample(0)
	require.Equal(t, 0.0, out.Read(0))

	// Feedback from initial pulse
	in.Write(0, 1)
	for i := 0; i < dsp.FrameSize; i++ {
		time.Write(i, 10)
		fbgain.Write(i, 0.9)
		mix.Write(i, 1)
		u.ProcessSample(i)
	}

	var sample []float64
	for i := 0; i < 20; i++ {
		sample = append(sample, out.Read(i))
	}
	require.Equal(t, []float64{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0.9, 0.9, 0.9}, sample)
}
