package unit

import (
	"testing"

	"buddin.us/musictheory"
	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestQuantize(t *testing.T) {
	builder := Builders()["quantize"]
	u, err := builder(Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	})
	require.NoError(t, err)

	u.Prop["intervals"].SetValue([]interface{}{
		musictheory.Perfect(1),
		musictheory.Perfect(5),
		musictheory.Minor(7),
	})

	u.In["tonic"].Write(0, dsp.Frequency(400, sampleRate).Float64())
	u.In["in"].Write(0, 0)
	u.ProcessSample(0)
	require.Equal(t, 0.009070294784580499, u.Out["out"].Out().Read(0))

	u.In["in"].Write(0, 0.3)
	u.ProcessSample(0)
	require.Equal(t, 0.013590086865094617, u.Out["out"].Out().Read(0))

	u.In["in"].Write(0, 1)
	u.ProcessSample(0)
	require.Equal(t, 0.016161427993475544, u.Out["out"].Out().Read(0))
}
