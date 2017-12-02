package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

type genOutput interface {
	CondProcessor
	FrameProcessor
	SampleProcessor
	Output
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
	builder := Builders()["gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["sine"].(*genSine)

	freq.Write(0, dsp.Frequency(100).Float64())
	out.ProcessSample(0)
	require.Equal(t, 0.568479975912154, out.Out().Read(0))

	freq.Write(1, dsp.Frequency(100).Float64())
	out.ProcessSample(1)
	require.Equal(t, 0.5801446938864443, out.Out().Read(1))
}

func TestGen_Saw(t *testing.T) {
	builder := Builders()["gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["saw"].(*genSaw)

	freq.Write(0, dsp.Frequency(100).Float64())
	out.ProcessSample(0)
	require.Equal(t, -0.9042167799057446, out.Out().Read(0))

	freq.Write(1, dsp.Frequency(100).Float64())
	out.ProcessSample(1)
	require.Equal(t, -0.8996816325134543, out.Out().Read(1))
}

func TestGen_Pulse(t *testing.T) {
	builder := Builders()["gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	pw := u.In["pulse-width"]
	out := u.Out["pulse"].(*genPulse)

	for i := 0; i < 2; i++ {
		for j := 0; j < dsp.FrameSize; j++ {
			pw.Write(j, 0.5)
			freq.Write(j, dsp.Frequency(100).Float64())
			out.ProcessSample(j)
		}
	}

	require.Equal(t, -1.0, out.Out().Read(0))
	require.Equal(t, 0.5639447500731988, out.Out().Read(170))
}

func TestGen_Triangle(t *testing.T) {
	builder := Builders()["gen"]
	u, err := builder(nil)
	require.NoError(t, err)

	freq := u.In["freq"]
	out := u.Out["triangle"].(*genTriangle)

	for i := 0; i < 2; i++ {
		for j := 0; j < dsp.FrameSize; j++ {
			freq.Write(j, dsp.Frequency(100).Float64())
			out.ProcessSample(j)
		}
	}

	require.Equal(t, 0.7994370616290737, out.Out().Read(0))
	require.Equal(t, -0.6227234073821785, out.Out().Read(170))
}
