package unit

import (
	"math/rand"
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

type genOutput interface {
	CondProcessor
	OutputProcessor
}

var _ = []genOutput{
	&genSine{},
	&genSaw{},
	&genPulse{},
	&genTriangle{},
	&genCluster{},
	&genNoise{},
}

func TestGen_Sine(t *testing.T) {
	rand.Seed(1)

	builder := Builders()["gen"]
	u, err := builder(Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	})
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["sine"].(*genSine)

	freq.Write(0, dsp.Frequency(100, sampleRate).Float64())
	out.ProcessSample(0)
	require.NotEqual(t, 0.0, out.Out().Read(0))

	freq.Write(1, dsp.Frequency(100, sampleRate).Float64())
	out.ProcessSample(1)
	require.NotEqual(t, 0.0, out.Out().Read(1))
}

func TestGen_Saw(t *testing.T) {
	rand.Seed(1)

	builder := Builders()["gen"]
	u, err := builder(Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	})
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["saw"].(*genSaw)

	freq.Write(0, dsp.Frequency(100, sampleRate).Float64())
	out.ProcessSample(0)
	require.NotEqual(t, 0.0, out.Out().Read(0))

	freq.Write(1, dsp.Frequency(100, sampleRate).Float64())
	out.ProcessSample(1)
	require.NotEqual(t, 0.0, out.Out().Read(1))
}

func TestGen_Pulse(t *testing.T) {
	rand.Seed(1)

	builder := Builders()["gen"]
	u, err := builder(Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	})
	require.NoError(t, err)

	freq := u.In["freq"]
	pw := u.In["pulse-width"]
	out := u.Out["pulse"].(*genPulse)

	for i := 0; i < 2; i++ {
		for j := 0; j < frameSize; j++ {
			pw.Write(j, 0.5)
			freq.Write(j, dsp.Frequency(100, sampleRate).Float64())
			out.ProcessSample(j)
		}
	}

	require.NotEqual(t, 0.0, out.Out().Read(0))
	require.NotEqual(t, 0.0, out.Out().Read(170))
}

func TestGen_Triangle(t *testing.T) {
	rand.Seed(1)

	builder := Builders()["gen"]
	u, err := builder(Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	})
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["triangle"].(*genTriangle)

	for i := 0; i < 2; i++ {
		for j := 0; j < frameSize; j++ {
			freq.Write(j, dsp.Frequency(100, sampleRate).Float64())
			out.ProcessSample(j)
		}
	}

	require.NotEqual(t, 0.0, out.Out().Read(0))
	require.NotEqual(t, 0.0, out.Out().Read(170))
}
