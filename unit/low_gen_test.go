package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

var _ = []genOutput{
	&lowGenSine{},
	&lowGenSaw{},
	&lowGenPulse{},
	&lowGenTriangle{},
}

func TestLowGen_Sine(t *testing.T) {
	builder := Builders()["low-gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["sine"].(*lowGenSine)

	// Only processes the first sample in frame
	freqv := dsp.Frequency(100).Float64()
	freq.Write(0, freqv)
	freq.Write(1, freqv)
	out.ProcessFrame(dsp.FrameSize)
	require.Equal(t, 0.7613121374252574, out.Out().Read(0))
	require.Equal(t, 0.0, out.Out().Read(1))
}

func TestLowGen_Saw(t *testing.T) {
	builder := Builders()["low-gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["saw"].(*lowGenSaw)

	// Only processes the first sample in frame
	freqv := dsp.Frequency(100).Float64()
	freq.Write(0, freqv)
	freq.Write(1, freqv)
	out.ProcessFrame(dsp.FrameSize)
	require.Equal(t, -0.974709126186973, out.Out().Read(0))
	require.Equal(t, 0.0, out.Out().Read(1))
}

func TestLowGen_Pulse(t *testing.T) {
	builder := Builders()["low-gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["pulse"].(*lowGenPulse)

	// Only processes the first sample in frame
	freqv := dsp.Frequency(100).Float64()
	freq.Write(0, freqv)
	freq.Write(1, freqv)
	out.ProcessFrame(dsp.FrameSize)
	require.Equal(t, 1.0, out.Out().Read(0))
	require.Equal(t, 0.0, out.Out().Read(1))
}

func TestLowGen_Triangle(t *testing.T) {
	builder := Builders()["low-gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["triangle"].(*lowGenTriangle)

	// Only processes the first sample in frame
	freqv := dsp.Frequency(100).Float64()
	freq.Write(0, freqv)
	freq.Write(1, freqv)
	out.ProcessFrame(dsp.FrameSize)
	require.Equal(t, -0.6555251334762077, out.Out().Read(0))
	require.Equal(t, 0.0, out.Out().Read(1))
}
