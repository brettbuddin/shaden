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
			_, err = fn(g)
			require.Nil(t, err)
			require.False(t, g.HasChanged())
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
	_, err = fn(g)
	require.Nil(t, err)
	require.True(t, g.HasChanged())
	require.True(t, u1.In["in"].HasSource())
	require.Equal(t, 1, u2.Out["out"].Out().DestinationCount())
}
