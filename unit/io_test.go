package unit

import (
	"testing"

	"buddin.us/lumen/dsp"
	"github.com/stretchr/testify/require"
)

func TestExposeIn(t *testing.T) {
	io := NewIO()
	in := io.NewIn("x", dsp.Float64(1))
	require.Equal(t, in, io.In["x"])
}

func TestExposeOut(t *testing.T) {
	io := NewIO()
	out := io.NewOut("x")
	require.Equal(t, out, io.Out["x"])
}

type output struct {
	out  *Out
	proc func(int)
}

func (o output) Out() *Out {
	return o.out
}

func (o output) ProcessSample(i int) {
	o.proc(i)
}

func TestExposeOutProcessor(t *testing.T) {
	io := NewIO()

	var called bool
	io.ExposeOutProcessor(output{
		out: NewOut("x", make([]float64, dsp.FrameSize)),
		proc: func(n int) {
			called = true
		},
	})
	out := io.Out["x"]
	require.NotNil(t, out)

	out.(SampleProcessor).ProcessSample(1)
	require.True(t, called)
}

func TestExposeProp(t *testing.T) {
	io := NewIO()
	p := io.NewProp("x", 1, nil)
	require.Equal(t, p, io.Prop["x"])
}
