package engine

import (
	"testing"

	"github.com/brettbuddin/shaden/dsp"
	"github.com/brettbuddin/shaden/unit"
	"github.com/stretchr/testify/require"
)

func TestFillConstant(t *testing.T) {
	var (
		tests = []struct {
			typ    string
			input  interface{}
			output interface{}
		}{
			{"float", 1.0, 1.0},
			{"int", 1, 1.0},
			{"hz", dsp.Frequency(440, sampleRate), 0.009977324263038548},
			{"ms", dsp.Duration(100, sampleRate), 4410.0},
			{"bpm", dsp.BPM(60, sampleRate), 2.2675736961451248e-05},
		}
	)

	for _, test := range tests {
		t.Run(test.typ, func(t *testing.T) {
			var (
				g  = NewGraph(frameSize)
				io = unit.NewIO("dummy", frameSize)
				u  = unit.NewUnit(io, nil)
			)

			err := g.Mount(u)
			require.Nil(t, err)
			io.NewIn("in", dsp.Float64(0))
			fn := PatchInput(u, map[string]interface{}{
				"in": test.input,
			}, false)
			err = fn(g)
			require.Nil(t, err)
			require.False(t, g.graph.HasChanged())
			require.Equal(t, test.output, io.In["in"].Read(0))
			require.Equal(t, test.output, io.In["in"].Read(1))
		})
	}
}

func TestPatch(t *testing.T) {
	g := NewGraph(frameSize)

	io1 := unit.NewIO("dummy1", frameSize)
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	u1 := unit.NewUnit(io1, nil)
	err := g.Mount(u1)
	require.Nil(t, err)

	io2 := unit.NewIO("dummy2", frameSize)
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	u2 := unit.NewUnit(io2, nil)
	err = g.Mount(u2)
	require.Nil(t, err)

	fn := PatchInput(u1, map[string]interface{}{
		"in": unit.OutRef{Unit: u2, Output: "out"},
	}, false)
	err = fn(g)
	require.Nil(t, err)
	require.True(t, g.graph.HasChanged())
	require.True(t, u1.In["in"].HasSource())
	require.Equal(t, 1, u2.Out["out"].Out().DestinationCount())
}

func TestEmitOutputs(t *testing.T) {
	g := NewGraph(frameSize)
	err := g.createSink(100, frameSize, sampleRate)
	require.NoError(t, err)

	io1 := unit.NewIO("dummy1", frameSize)
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	unit1 := unit.NewUnit(io1, nil)
	err = g.Mount(unit1)
	require.NoError(t, err)

	io2 := unit.NewIO("dummy2", frameSize)
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	unit2 := unit.NewUnit(io2, nil)
	err = g.Mount(unit2)
	require.NoError(t, err)

	left := unit.OutRef{Unit: unit1, Output: "out"}
	right := unit.OutRef{Unit: unit2, Output: "out"}

	err = EmitOutputs(left, right)(g)
	require.NoError(t, err)
	require.True(t, g.sink.In["l"].HasSource())
	require.True(t, g.sink.In["r"].HasSource())
	require.Equal(t, 1, unit1.Out["out"].Out().DestinationCount())
	require.Equal(t, 1, unit2.Out["out"].Out().DestinationCount())
}

func TestSwapUnit_ConstantInput(t *testing.T) {
	g := NewGraph(frameSize)
	err := g.createSink(100, frameSize, sampleRate)
	require.NoError(t, err)

	io1 := unit.NewIO("dummy", frameSize)
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	unit1 := unit.NewUnit(io1, nil)
	err = g.Mount(unit1)
	require.NoError(t, err)

	io2 := unit.NewIO("dummy", frameSize)
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	unit2 := unit.NewUnit(io2, nil)
	err = g.Mount(unit2)
	require.NoError(t, err)

	io3 := unit.NewIO("dummy-dest", frameSize)
	io3.NewIn("in", dsp.Float64(0))
	io3.NewOut("out")
	unit3 := unit.NewUnit(io3, nil)
	err = g.Mount(unit3)
	require.NoError(t, err)

	require.NoError(t, g.Patch(3.0, unit1.In["in"]))
	require.NoError(t, g.Patch(unit1.Out["out"], unit3.In["in"]))

	err = SwapUnit(unit1, unit2)(g)
	require.NoError(t, err)
	require.Equal(t, unit2.Out["out"].Out(), unit3.In["in"].Source())
	require.Equal(t, dsp.Float64(3.0), unit1.In["in"].Constant())
	require.Equal(t, 0, unit1.Out["out"].Out().DestinationCount())
}

func TestSwapUnit_SourceInput(t *testing.T) {
	g := NewGraph(frameSize)
	err := g.createSink(100, frameSize, sampleRate)
	require.NoError(t, err)

	io1 := unit.NewIO("dummy", frameSize)
	io1.NewIn("in", dsp.Float64(0))
	io1.NewOut("out")
	unit1 := unit.NewUnit(io1, nil)
	err = g.Mount(unit1)
	require.NoError(t, err)

	io2 := unit.NewIO("dummy", frameSize)
	io2.NewIn("in", dsp.Float64(0))
	io2.NewOut("out")
	unit2 := unit.NewUnit(io2, nil)
	err = g.Mount(unit2)
	require.NoError(t, err)

	io3 := unit.NewIO("dummy", frameSize)
	io3.NewIn("in", dsp.Float64(0))
	io3.NewOut("out")
	unit3 := unit.NewUnit(io3, nil)
	err = g.Mount(unit3)
	require.NoError(t, err)

	io4 := unit.NewIO("dummy-dest", frameSize)
	io4.NewIn("in", dsp.Float64(0))
	io4.NewOut("out")
	unit4 := unit.NewUnit(io4, nil)
	err = g.Mount(unit4)
	require.NoError(t, err)

	require.NoError(t, g.Patch(unit1.Out["out"], unit2.In["in"]))
	require.NoError(t, g.Patch(unit2.Out["out"], unit4.In["in"]))

	err = SwapUnit(unit2, unit3)(g)
	require.NoError(t, err)
	require.Equal(t, unit3.Out["out"].Out(), unit4.In["in"].Source())
	require.Equal(t, unit1.Out["out"].Out(), unit3.In["in"].Source())
	require.False(t, unit2.In["in"].HasSource())
	require.Equal(t, 0, unit2.Out["out"].Out().DestinationCount())
}
