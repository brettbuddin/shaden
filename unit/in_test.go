package unit

import (
	"testing"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/graph"
	"github.com/stretchr/testify/require"
)

func TestIn_BlockRead(t *testing.T) {
	in := NewIn("in", dsp.Float64(0))
	in.frame[0] = 10
	in.frame[5] = 20
	require.Equal(t, 10.0, in.Read(0))
	require.Equal(t, 20.0, in.Read(5))
}

func TestIn_SampleRead(t *testing.T) {
	in := NewIn("in", dsp.Float64(0))
	in.Mode = Sample
	in.frame[0] = 10
	in.frame[5] = 20
	require.Equal(t, 10.0, in.Read(1))
	require.Equal(t, 20.0, in.Read(6))
}

func TestIn_ReadSlow(t *testing.T) {
	in := NewIn("in", dsp.Float64(0))
	in.frame[0] = 20
	in.frame[64] = 30

	add2 := func(v float64) float64 { return v + 2 }

	require.Equal(t, 22.0, in.ReadSlow(0, add2))
	require.Equal(t, 22.0, in.ReadSlow(6, add2))
	require.Equal(t, 32.0, in.ReadSlow(64, add2))
}

func TestIn_Fill(t *testing.T) {
	in := NewIn("in", dsp.Float64(0))
	in.Fill(dsp.Float64(101))
	require.Equal(t, 101.0, in.Read(10))
}

func TestIn_CoupleOutput(t *testing.T) {
	g := graph.New()

	io1 := NewIO()
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	u1 := NewUnit(io1, "example1", nil)
	require.Equal(t, 0, u1.ExternalNeighborCount())

	io2 := NewIO()
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	u2 := NewUnit(io2, "example2", nil)
	require.Equal(t, 0, u2.ExternalNeighborCount())

	io2.In["in"].Fill(dsp.Float64(101))
	io1.Out["out"].Out().frame[10] = 102

	err := u1.Attach(g)
	require.NoError(t, err)
	err = u2.Attach(g)
	require.NoError(t, err)

	err = Patch(g, u1.Out["out"], u2.In["in"])
	require.NoError(t, err)

	require.Equal(t, 102.0, io2.In["in"].Read(10))
}

func TestIn_ReadControlRate(t *testing.T) {
	in := NewIn("in", dsp.Float64(0))
	in.Mode = Sample
	in.Couple(&Out{
		unit:  &Unit{rate: RateControl},
		frame: newFrame(),
	})
	in.frame[0] = 10
	in.frame[5] = 20
	require.Equal(t, 10.0, in.Read(0))
	require.Equal(t, 10.0, in.Read(5))
}
