package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	builder := Builders()["filter"]
	u, err := builder(Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	})
	require.NoError(t, err)

	var (
		in     = u.In["in"]
		cutoff = u.In["cutoff"]
		lp     = u.Out["lp"].Out()
		bp     = u.Out["bp"].Out()
		hp     = u.Out["hp"].Out()
	)

	cutoff.Write(0, dsp.Frequency(100, sampleRate).Float64())
	in.Write(0, 1)
	u.ProcessSample(0)
	require.Equal(t, 0.00012781289351266066, lp.Read(0))
	require.Equal(t, 0.0157444826565947, bp.Read(0))
	require.Equal(t, 0.9841277044498926, hp.Read(0))
	require.Equal(t, 1.0, lp.Read(0)+bp.Read(0)+hp.Read(0))
}
