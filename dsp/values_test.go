package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrequency(t *testing.T) {
	freq := Frequency(440, sampleRate)
	require.Equal(t, 0.009977324263038548, freq.Float64())
	require.Equal(t, "440.00Hz", freq.String())
}

func TestPitch(t *testing.T) {
	pitch, err := ParsePitch("A4", sampleRate)
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, pitch.Float64())
	require.Equal(t, "A4", pitch.String())
}

func TestMS(t *testing.T) {
	ms := Duration(1, sampleRate)
	require.Equal(t, 44.1, ms.Float64())
	require.Equal(t, "1.00ms", ms.String())
}

func TestBPM(t *testing.T) {
	tempo := BPM(60, sampleRate)
	require.Equal(t, 2.2675736961451248e-05, tempo.Float64())
	require.Equal(t, "60.00BPM", tempo.String())
}
