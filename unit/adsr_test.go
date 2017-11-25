package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestADSR(t *testing.T) {
	builder := Builders()["adsr"]
	u, err := builder(nil)
	require.NoError(t, err)

	var (
		gate    = u.In["gate"]
		attack  = u.In["attack"]
		decay   = u.In["decay"]
		release = u.In["release"]
		out     = u.Out["out"].Out()

		samples []float64
	)

	gate.Write(0, 1)
	for i := 0; i < dsp.FrameSize; i++ {
		attack.Write(i, 3)
		decay.Write(i, 5)
		release.Write(i, 10)
		u.ProcessSample(i)
		samples = append(samples, out.Read(i))
	}
	require.Equal(t, []float64{
		0,
		0.7931226244422468,
		0.9634299049219617,
		1,
		0.692631006358899,
		0.5705084798784709,
		0.5219872829376465,
		0.5027090496712591,
		0.495049504950495,
		0.3083477702616633,
		0.19066409695918154,
		0.1164845667847695,
		0.06972699589595815,
		0.040254304032288676,
		0.021676787083078075,
		0.009966823921428891,
		0.002585684793844481,
		0,
	}, samples[:18])
}
