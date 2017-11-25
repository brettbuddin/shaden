package unit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGate(t *testing.T) {
	builder := Builders()["gate"]
	u, err := builder(nil)
	require.NoError(t, err)

	var (
		in      = u.In["in"]
		control = u.In["control"]
		mode    = u.In["mode"]
		out     = u.Out["out"].Out()
	)

	// Amplitude mode
	in.Write(0, 1)
	mode.Write(0, float64(gateModeAmp))
	control.Write(0, 0.5)
	u.ProcessSample(0)
	require.Equal(t, 0.5, out.Read(0))

	// Combo mode
	in.Write(0, 1)
	mode.Write(0, float64(gateModeCombo))
	control.Write(0, 0.5)
	u.ProcessSample(0)
	require.Equal(t, 0.32795609653857716, out.Read(0))

	// Low-pass mode
	in.Write(0, 1)
	mode.Write(0, float64(gateModeLP))
	control.Write(0, 0.1)
	u.ProcessSample(0)
	require.Equal(t, 0.9108075595968936, out.Read(0))
}
