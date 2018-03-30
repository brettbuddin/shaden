package unit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlope(t *testing.T) {
	builder := Builders()["slope"]
	u, err := builder(Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	})
	require.NoError(t, err)

	var (
		trigger = u.In["trigger"]
		rise    = u.In["rise"]
		fall    = u.In["fall"]
		out     = u.Out["out"].Out()

		samples []float64
	)

	trigger.Write(0, 1)
	for i := 0; i < frameSize; i++ {
		rise.Write(i, 3)
		fall.Write(i, 5)
		u.ProcessSample(i)
		samples = append(samples, out.Read(i))
	}
	require.Equal(t, []float64{
		0,
		0.7931226244422468,
		0.9634299049219617,
		1,
		0.3912888557303687,
		0.149438362112266,
		0.05334736424906474,
		0.01516890229014065,
	}, samples[:8])
}
