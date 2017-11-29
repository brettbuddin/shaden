package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

func TestDelay(t *testing.T) {
	builder := Builders()["delay"]
	u, err := builder(nil)
	require.NoError(t, err)

	var (
		in     = u.In["in"]
		mix    = u.In["mix"]
		fbgain = u.In["fb-gain"]
		time   = u.In["time"]
		out    = u.Out["out"].Out()
	)

	// Mix of dry and wet
	in.Write(0, 1)
	u.ProcessSample(0)
	require.Equal(t, 1.0, out.Read(0))

	// Only wet signal
	mix.Write(0, 1)
	in.Write(0, 1)
	u.ProcessSample(0)
	require.Equal(t, -0.0050000000000000044, out.Read(0))

	// Feedback from initial pulse
	in.Write(0, 1)
	for i := 0; i < dsp.FrameSize; i++ {
		time.Write(i, 10)
		fbgain.Write(i, 0.9)
		mix.Write(i, 1)
		u.ProcessSample(i)
	}

	var sample []float64
	for i := 0; i < 20; i++ {
		sample = append(sample, out.Read(i))
	}
	require.Equal(t, []float64{
		-0.004975000000000005,
		-0.004950125000000005,
		-0.004925374375000005,
		-0.004900747503125005,
		-0.0048762437656093794,
		-0.004851862546781332,
		-0.004827603234047425,
		0.9951965347821228,
		0.9902205521082121,
		0.985269449347671,
		-0.019656897899067327,
		-0.01955861340957199,
		-0.01946082034252413,
		-0.01936351624081151,
		-0.019266698659607454,
		-0.019170365166309416,
		-0.01907451334047787,
		0.8810208592262245,
		0.8766157549300934,
		0.8722326761554429,
	}, sample)
}
