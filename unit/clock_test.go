package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestClock(t *testing.T) {
	builder := Builders()["clock"]
	u, err := builder(nil)
	require.NoError(t, err)

	// High frequency to close the range we have to iterate over
	freq := dsp.Frequency(1000).Float64()

	tempo := u.In["tempo"]
	out := u.Out["out"].Out()

	for i := 0; i < 5; i++ {
		tempo.Write(i, freq)
		u.ProcessSample(i)
		require.Equal(t, 1.0, out.Read(i))
	}

	for i := 5; i < 44; i++ {
		tempo.Write(i, freq)
		u.ProcessSample(i)
		require.Equal(t, -1.0, out.Read(i))
	}

	for i := 44; i < 46; i++ {
		tempo.Write(i, freq)
		u.ProcessSample(i)
		require.Equal(t, 1.0, out.Read(i))
	}
}
