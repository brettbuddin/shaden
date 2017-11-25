package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestRandemSeries(t *testing.T) {
	builder := Builders()["random-series"]
	u, err := builder(nil)
	require.NoError(t, err)

	var (
		clock   = u.In["clock"]
		size    = u.In["size"]
		trigger = u.In["trigger"]
		gate    = u.Out["gate"].Out()
		value   = u.Out["value"].Out()
	)

	var (
		clockv        = 1.0
		gates, values []float64
		triggerAt     = 0
	)
	for i := 0; i < dsp.FrameSize; i++ {
		clock.Write(i, clockv)
		if i == triggerAt {
			trigger.Write(i, 1)
		} else {
			trigger.Write(i, -1)
		}
		size.Write(i, 5)
		u.ProcessSample(i)

		if clockv > 0 {
			clockv = -1
		} else {
			clockv = 1
		}
		gates = append(gates, gate.Read(i))
		values = append(values, value.Read(i))
	}

	require.Equal(t, gates[1:11], gates[11:21])
	require.Equal(t, values[1:11], values[11:21])
}
