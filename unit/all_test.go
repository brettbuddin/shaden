package unit

import (
	"fmt"
	"testing"

	"buddin.us/shaden/dsp"
	"github.com/stretchr/testify/require"
)

var A4 = dsp.Frequency(440).Float64()

func TestAllUnits(t *testing.T) {
	var tests = []struct {
		unit     string
		config   Config
		scenario []scenario
	}{
		{
			unit:   "adjust",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":     []float64{1, 1, 1, 1},
						"gain":   []float64{1, 3, 2, 4},
						"offset": []float64{1, -1, 1, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{2, 2, 3, 5},
					},
				},
			},
		},
		{
			unit:   "abs",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{1, -2, 1, -1},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 2, 1, 1},
					},
				},
			},
		},
		{
			unit:   "sum",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{1, 1, 1, -1},
						"y": []float64{1, 3, 2, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{2, 4, 3, 3},
					},
				},
			},
		},
		{
			unit:   "ceil",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0.1, 2.7},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 3},
					},
				},
			},
		},
		{
			unit:   "floor",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0.1, 2.7},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 2},
					},
				},
			},
		},
		{
			unit:   "invert",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{-0.1, 2.7},
					},
					outputs: map[string][]float64{
						"out": []float64{0.1, -2.7},
					},
				},
			},
		},
		{
			unit:   "noop",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{-0.1, 2.7},
					},
					outputs: map[string][]float64{
						"out": []float64{-0.1, 2.7},
					},
				},
			},
		},
		{
			unit:   "not",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{-0.1, 2.7},
					},
					outputs: map[string][]float64{
						"out": []float64{1, -1},
					},
				},
			},
		},
		{
			unit:   "val-gate",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in": []float64{-0.1, 2.7},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, 1},
					},
				},
			},
		},
		{
			unit:   "diff",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{1, 1, 1, -1},
						"y": []float64{1, 3, 2, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{0, -2, -1, -5},
					},
				},
			},
		},
		{
			unit:   "mult",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{1, 1, 1, -1},
						"y": []float64{1, 3, 2, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 3, 2, -4},
					},
				},
			},
		},
		{
			unit:   "div",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{1, 1, 1, -1},
						"y": []float64{0, 3, 2, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 1.0 / 3.0, 0.5, -1.0 / 4},
					},
				},
			},
		},
		{
			unit:   "mod",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{1, 1, 10, -4},
						"y": []float64{0, 3, 2, 2},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 1, 0, 0},
					},
				},
			},
		},
		{
			unit:   "gt",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{1, 1, 10, -4},
						"y": []float64{0, 3, 2, 2},
					},
					outputs: map[string][]float64{
						"out": []float64{1, -1, 1, -1},
					},
				},
			},
		},
		{
			unit:   "lt",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, 3, 2, 2},
						"y": []float64{1, 1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, -1, 1, -1},
					},
				},
			},
		},
		{
			unit:   "and",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, 3, 2, 2},
						"y": []float64{1, 1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, 1, 1, -1},
					},
				},
			},
		},
		{
			unit:   "or",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, -1, 1, 1},
					},
				},
			},
		},
		{
			unit:   "xor",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, -1, -1, 1},
					},
				},
			},
		},
		{
			unit:   "nand",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 1, -1, 1},
					},
				},
			},
		},
		{
			unit:   "nor",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, 1, -1, -1},
					},
				},
			},
		},
		{
			unit:   "imply",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 1, 1, -1},
					},
				},
			},
		},
		{
			unit:   "xnor",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, 1, 1, -1},
					},
				},
			},
		},
		{
			unit:   "max",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{1, -1, 10, 2},
					},
				},
			},
		},
		{
			unit:   "min",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, -3, 2, 2},
						"y": []float64{1, -1, 10, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{0, -3, 2, -4},
					},
				},
			},
		},
		{
			unit:   "clip",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":    []float64{3, 1, 1, -3},
						"level": []float64{2, 2, 2, 2},
						"soft":  []float64{1, 0, 0, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{1.083333333333333333, 1, 1, -1.083333333333333333},
					},
				},
			},
		},
		{
			unit:   "clock-mult",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":   []float64{-1, -1, 1, -1, -1, 1, -1, -1, 1, -1, -1},
						"mult": []float64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, -1, 1, -1},
					},
				},
			},
		},
		{
			unit:   "clock-div",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":  []float64{-1, 1, -1, 1, -1, 1, -1, 1},
						"div": []float64{2, 2, 2, 2, 2, 2, 2, 2},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, -1, -1, 1, -1, -1, -1},
					},
				},
			},
		},
		{
			unit:   "cond",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"cond": []float64{1, -1, 0},
						"x":    []float64{1, 3, 4},
						"y":    []float64{2, 4, 5},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 4, 5},
					},
				},
			},
		},
		{
			unit:   "count",
			config: nil,
			scenario: []scenario{
				{
					description: "basic counting",
					inputs: map[string][]float64{
						"trigger": []float64{-1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 1, 1, 2},
					},
				},
				{
					description: "basic counting by interval",
					inputs: map[string][]float64{
						"trigger": []float64{-1, 1, -1, 1},
						"step":    []float64{3, 3, 3, 3},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 3, 3, 6},
					},
				},
				{
					description: "manual reset",
					inputs: map[string][]float64{
						"trigger": []float64{-1, 1, -1, 1},
						"reset":   []float64{-1, -1, 1, -1},
					},
					outputs: map[string][]float64{
						"out":   []float64{0, 1, 0, 1},
						"reset": []float64{-1, -1, 1, -1},
					},
				},
				{
					description: "reset gate on wrap",
					inputs: map[string][]float64{
						"trigger": []float64{-1, 1, -1, 1},
						"limit":   []float64{2, 2, 2, 2},
					},
					outputs: map[string][]float64{
						"out":   []float64{0, 1, 1, 0},
						"reset": []float64{-1, -1, -1, 1},
					},
				},
			},
		},
		{
			unit:   "xfade",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"a":   []float64{1, 1, 1, 1},
						"b":   []float64{1, 3, 3, 3},
						"mix": []float64{0, -1, 1, 0.5},
					},
					outputs: map[string][]float64{
						"out": []float64{2, 1, 3, 3.5},
					},
				},
			},
		},
		{
			unit:   "xfeed",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"a":      []float64{1, 1, 1},
						"b":      []float64{3, 3, 3},
						"amount": []float64{0, 0.5, 1},
					},
					outputs: map[string][]float64{
						"a": []float64{1, 2.5, 4},
						"b": []float64{3, 3.5, 4},
					},
				},
			},
		},
		{
			unit:   "pan",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":  []float64{1, 1, 1, 1},
						"pan": []float64{-1, 0, 0.5, 1},
					},
					outputs: map[string][]float64{
						"a": []float64{1, 1, 0.5, 0},
						"b": []float64{0, 1, 1, 1},
					},
				},
			},
		},
		{
			unit:   "fold",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in": []float64{1, 3, -4},
					},
					outputs: map[string][]float64{
						"out": []float64{0.6000000000000001, -0.2000000000000004, -0.8},
					},
				},
			},
		},
		{
			unit:   "mux",
			config: Config{"size": 4},
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"0":      []float64{1, 1, 1, 1},
						"1":      []float64{2, 2, 2, 2},
						"2":      []float64{3, 3, 3, 3},
						"3":      []float64{4, 4, 4, 4},
						"select": []float64{0, 1, 2, 3},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 2, 3, 4},
					},
				},
			},
		},
		{
			unit:   "mix",
			config: nil,
			scenario: []scenario{
				{
					description: "all inputs",
					inputs: map[string][]float64{
						"0/in": []float64{1, 1, 1, 1},
						"1/in": []float64{2, 2, 2, 2},
						"2/in": []float64{3, 3, 3, 3},
						"3/in": []float64{4, 4, 4, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{10, 10, 10, 10},
					},
				},
				{
					description: "one input attenuated",
					inputs: map[string][]float64{
						"0/in":    []float64{1, 1, 1, 1},
						"0/level": []float64{0.1, 0.1, 0.1, 0.1},
						"1/in":    []float64{2, 2, 2, 2},
						"2/in":    []float64{3, 3, 3, 3},
						"3/in":    []float64{4, 4, 4, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{9.1, 9.1, 9.1, 9.1},
					},
				},
			},
		},
		{
			unit:   "panmix",
			config: nil,
			scenario: []scenario{
				{
					description: "all inputs",
					inputs: map[string][]float64{
						"0/in": []float64{1, 1, 1, 1},
						"1/in": []float64{2, 2, 2, 2},
						"2/in": []float64{3, 3, 3, 3},
						"3/in": []float64{4, 4, 4, 4},
					},
					outputs: map[string][]float64{
						"a": []float64{10, 10, 10, 10},
						"b": []float64{10, 10, 10, 10},
					},
				},
				{
					description: "one input attenuated with pan",
					inputs: map[string][]float64{
						"0/in":    []float64{1, 1, 1, 1},
						"0/level": []float64{0.1, 0.1, 0.1, 0.1},
						"0/pan":   []float64{0.1, 0.1, 0.1, 0.1},
						"1/in":    []float64{2, 2, 2, 2},
						"2/in":    []float64{3, 3, 3, 3},
						"3/in":    []float64{4, 4, 4, 4},
					},
					outputs: map[string][]float64{
						"a": []float64{9.09, 9.09, 9.09, 9.09},
						"b": []float64{9.1, 9.1, 9.1, 9.1},
					},
				},
			},
		},
		{
			unit:   "switch",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"0":       []float64{1, 1, 1, 1},
						"1":       []float64{2, 2, 2, 2},
						"2":       []float64{3, 3, 3, 3},
						"3":       []float64{4, 4, 4, 4},
						"trigger": []float64{-1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 2, 2, 3},
					},
				},
			},
		},
		{
			unit:   "demux",
			config: Config{"size": 4},
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":     []float64{1, 1, 1, 1},
						"select": []float64{0, 1, 2, 3},
					},
					outputs: map[string][]float64{
						"0": []float64{1, 0, 0, 0},
						"1": []float64{0, 1, 0, 0},
						"2": []float64{0, 0, 1, 0},
						"3": []float64{0, 0, 0, 1},
					},
				},
			},
		},
		{
			unit:   "latch",
			config: Config{"size": 4},
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":      []float64{1, 2, 3, 4},
						"trigger": []float64{-1, -1, 1, -1},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 1, 3, 3},
					},
				},
			},
		},
		{
			unit:   "toggle",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"trigger": []float64{-1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 1, 1, -1},
					},
				},
			},
		},
		{
			unit:   "transpose",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":        []float64{A4, A4},
						"semitones": []float64{0, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{0.009977324263038548, 0.010570606837144897},
					},
				},
			},
		},
		{
			unit:   "transpose-interval",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":      []float64{A4, A4},
						"quality": []float64{0, 1},
						"step":    []float64{1, 2},
					},
					outputs: map[string][]float64{
						"out": []float64{0.009977324263038548, 0.010570606837144897},
					},
				},
			},
		},
		{
			unit:   "chebyshev",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in": []float64{1, 1},
						"a":  []float64{0.5, 1},
						"b":  []float64{0.1, 1},
						"c":  []float64{0.1, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{0.7, 3},
					},
				},
			},
		},
		{
			unit:   "bipolar",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{0, 0.25, 0.50, 0.75, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, -0.5, 0, 0.5, 1},
					},
				},
			},
		},
		{
			unit:   "unipolar",
			config: nil,
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"x": []float64{-1, -0.5, 0, 0.5, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 0.25, 0.50, 0.75, 1},
					},
				},
			},
		},
	}

	builders := Builders()
	for _, test := range tests {
		for i, s := range test.scenario {
			name := test.unit
			if s.description != "" {
				name += "_" + s.description
			}
			t.Run(name, func(t *testing.T) {
				builder := builders[test.unit]
				u, err := builder(test.config)
				require.NoError(t, err)
				s.TestUnit(t, i, u)
			})
		}
	}
}

type scenario struct {
	description string
	inputs      map[string][]float64
	outputs     map[string][]float64
}

func (s scenario) TestUnit(t *testing.T, index int, u *Unit) {
	var max int
	for _, values := range s.inputs {
		if max < len(values) {
			max = len(values)
		}
	}

	for name, values := range s.inputs {
		for i, v := range values {
			u.In[name].Write(i, v)
		}
	}

	for i := 0; i < max; i++ {
		u.ProcessSample(i)
	}

	for name, values := range s.outputs {
		for i, v := range values {
			require.Equal(t, v, u.Out[name].Out().Read(i), fmt.Sprintf("output %q -> scenario %d -> sample %d", name, index, i))
		}
	}
}
