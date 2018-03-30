package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSVFilter(t *testing.T) {
	filter := &SVFilter{
		Poles:     4,
		Cutoff:    Frequency(100, sampleRate).Float64(),
		Resonance: 1,
	}

	in := 1.0
	hp, bp, lp := filter.Tick(in)
	require.Equal(t, 0.00012781289351266066, hp)
	require.Equal(t, 0.0157444826565947, bp)
	require.Equal(t, 0.9841277044498926, lp)
	require.Equal(t, in, hp+bp+lp)
}

func TestSimpleLowPass(t *testing.T) {
	filter := NewFilter(LowPass, 4)
	filter.Cutoff = Frequency(50, sampleRate).Float64()
	var v float64
	for i := 0; i < frameSize*10; i++ {
		v = filter.Tick(1.0)
	}
	require.Equal(t, 0.9999865165554183, v)
}

func TestSimpleHighPass(t *testing.T) {
	filter := NewFilter(HighPass, 4)
	filter.Cutoff = Frequency(50, sampleRate).Float64()
	var v float64
	for i := 0; i < frameSize*10; i++ {
		v = filter.Tick(1.0)
	}
	require.Equal(t, 1.348344458174111e-05, v)
}
