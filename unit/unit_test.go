package unit

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/graph"
)

const sampleRate = 44100.0
const frameSize = 256

var namePattern = regexp.MustCompile(`\w+`)

func TestRegisteredNames(t *testing.T) {
	for k := range Builders() {
		require.True(t, namePattern.Match([]byte(k)))
	}
}

func TestUnit_GraphAttachment(t *testing.T) {
	io := NewIO("example", frameSize)
	io.NewIn("in", dsp.Float64(0))
	io.NewOut("out")

	u := NewUnit(io, nil)
	g := graph.New()
	require.Equal(t, 0, g.Size())

	err := u.Attach(g)
	require.NoError(t, err)
	require.Equal(t, 3, g.Size())
	require.Len(t, u.node.InNeighbors(), 1)
	require.Len(t, u.node.OutNeighbors(), 1)

	err = u.Detach(g)
	require.NoError(t, err)
	require.Equal(t, 0, g.Size())
}

func TestUnit_Processable(t *testing.T) {
	var (
		io = NewIO("example", frameSize)
		u  = NewUnit(io, noopSampleProc{})
		g  = graph.New()
	)
	require.False(t, u.IsProcessable())
	err := u.Attach(g)
	require.NoError(t, err)
	require.True(t, u.IsProcessable())
}

func TestUnit_NotProcessable(t *testing.T) {
	var (
		io = NewIO("example", frameSize)
		u  = NewUnit(io, nil)
		g  = graph.New()
	)
	require.False(t, u.IsProcessable())
	err := u.Attach(g)
	require.NoError(t, err)
	require.False(t, u.IsProcessable())
}

func TestUnit_ProxyProcessFrame(t *testing.T) {
	var (
		recorded int

		io        = NewIO("example", frameSize)
		processor = frameProcessor{frame: func(n int) { recorded = n }}
		u         = NewUnit(io, processor)
		g         = graph.New()
	)
	err := u.Attach(g)
	require.NoError(t, err)
	u.ProcessFrame(10)
	require.Equal(t, 10, recorded)
}

func TestUnit_ProcessFrameHonorRateSetting(t *testing.T) {
	var (
		sampleProcCalled bool

		io        = NewIO("example", frameSize)
		processor = sampleProcessor{fn: func(i int) {
			require.Equal(t, 0, i)
		}}
		u = NewUnit(io, processor)
		g = graph.New()
	)
	u.rate = RateControl
	err := u.Attach(g)
	require.NoError(t, err)
	u.ProcessFrame(10)
	require.False(t, sampleProcCalled)
}

func TestUnit_ExternalNeighborCount(t *testing.T) {
	g := graph.New()

	io1 := NewIO("example1", frameSize)
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	u1 := NewUnit(io1, nil)
	require.Equal(t, 0, u1.ExternalNeighborCount())

	io2 := NewIO("example2", frameSize)
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	u2 := NewUnit(io2, nil)
	require.Equal(t, 0, u2.ExternalNeighborCount())

	require.NoError(t, u1.Attach(g))
	require.NoError(t, u2.Attach(g))

	require.NoError(t, Patch(g, u1.Out["out"], u2.In["in"]))

	require.Equal(t, 1, u1.ExternalNeighborCount())
	require.Equal(t, 1, u2.ExternalNeighborCount())

	require.NoError(t, Unpatch(g, u2.In["in"]))
	require.Equal(t, 0, u1.ExternalNeighborCount())
	require.Equal(t, 0, u2.ExternalNeighborCount())
}

func TestUnit_DetachInboundConnectionRemoval(t *testing.T) {
	g := graph.New()

	io1 := NewIO("example1", frameSize)
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	u1 := NewUnit(io1, nil)
	require.Equal(t, 0, u1.ExternalNeighborCount())

	io2 := NewIO("example2", frameSize)
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	u2 := NewUnit(io2, nil)
	require.Equal(t, 0, u2.ExternalNeighborCount())

	require.NoError(t, u1.Attach(g))
	require.NoError(t, u2.Attach(g))

	require.NoError(t, Patch(g, u1.Out["out"], u2.In["in"]))

	require.Equal(t, 1, u1.ExternalNeighborCount())
	require.Equal(t, 1, u2.ExternalNeighborCount())

	require.NoError(t, u2.Detach(g))

	require.Equal(t, 0, u1.ExternalNeighborCount())
}

func TestUnit_DetachOutboundConnectionRemoval(t *testing.T) {
	g := graph.New()

	io1 := NewIO("example1", frameSize)
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	u1 := NewUnit(io1, nil)
	require.Equal(t, 0, u1.ExternalNeighborCount())

	io2 := NewIO("example2", frameSize)
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	u2 := NewUnit(io2, nil)
	require.Equal(t, 0, u2.ExternalNeighborCount())

	require.NoError(t, u1.Attach(g))
	require.NoError(t, u2.Attach(g))

	require.NoError(t, Patch(g, u1.Out["out"], u2.In["in"]))

	require.Equal(t, 1, u1.ExternalNeighborCount())
	require.Equal(t, 1, u2.ExternalNeighborCount())

	require.NoError(t, u1.Detach(g))

	require.Equal(t, 0, u2.ExternalNeighborCount())
}

func TestUnit_Close(t *testing.T) {
	var closeCalled int

	io := NewIO("example1", frameSize)
	io.NewIn("in", dsp.Float64(0))
	out := NewOut("out", make([]float64, frameSize))
	io.ExposeOutputProcessor(outProcessorCloser{
		closer: closer{
			fn: func() error {
				closeCalled++
				return nil
			},
		},
		outProcessor: outProcessor{out: out},
	})

	processor := sampleProcessorCloser{
		closer: closer{fn: func() error {
			closeCalled++
			return nil
		}},
	}
	u := NewUnit(io, processor)

	err := u.Close()
	require.NoError(t, err)
	require.Equal(t, 2, closeCalled)
}
