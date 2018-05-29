package unit

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/brettbuddin/shaden/dsp"
)

var A4 = dsp.Frequency(440, 44100.0).Float64()

func TestAllUnits(t *testing.T) {
	var tests = []struct {
		unit         string
		configValues map[string]interface{}
		scenario     []scenario
	}{
		{
			unit: "adjust",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":   []float64{1, 1, 1, 1},
						"mult": []float64{1, 3, 2, 4},
						"add":  []float64{1, -1, 1, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{2, 2, 3, 5},
					},
				},
			},
		},
		{
			unit: "abs",
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
			unit: "sum",
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
			unit: "ceil",
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
			unit: "floor",
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
			unit: "invert",
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
			unit: "noop",
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
			unit: "not",
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
			unit: "val-gate",
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
			unit: "diff",
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
			unit: "mult",
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
			unit: "div",
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
			unit: "mod",
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
			unit: "gt",
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
			unit: "lt",
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
			unit: "and",
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
			unit: "or",
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
			unit: "xor",
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
			unit: "nand",
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
			unit: "nor",
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
			unit: "imply",
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
			unit: "xnor",
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
			unit: "max",
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
			unit: "min",
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
			unit: "clip",
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
			unit: "overload",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":   repeatControl([]float64{3, 1, 3, -3}),
						"gain": padControl([]float64{1, 1, 10, 1}),
					},
					outputs: map[string][]float64{
						"out": repeatControl([]float64{0.950212931632136, 0.6321205588285577, 0.9999999999999064, -0.950212931632136}),
					},
				},
			},
		},
		{
			unit: "clock-mult",
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
			unit: "clock-div",
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
			unit: "cond",
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
			unit: "count",
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
			unit: "xfade",
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
			unit: "xfeed",
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
			unit: "pan",
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
			unit: "fold",
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
			unit:         "mux",
			configValues: map[string]interface{}{"size": 4},
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
			unit: "mix",
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
			unit: "panmix",
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
			unit: "switch",
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
			unit:         "demux",
			configValues: map[string]interface{}{"size": 4},
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
			unit:         "latch",
			configValues: map[string]interface{}{"size": 4},
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
			unit: "toggle",
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
			unit: "transpose",
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
			unit: "transpose-interval",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":      []float64{A4, A4, A4, A4, A4},
						"quality": []float64{0, 1, 2, 3, 4},
						"step":    []float64{1, 2, 3, 5, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{0.009977324263038548, 0.010570606837144897, 0.01257064086062912, 0.01411006728898326, 0.01411006728898326},
					},
				},
			},
		},
		{
			unit: "chebyshev",
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
			unit: "cluster",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"freq": []float64{A4, A4},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 0.6008394124819831},
					},
				},
			},
		},
		{
			unit: "decimate",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":   repeatControl([]float64{1, 0.5}),
						"bits": padControl([]float64{24, 2}),
					},
					outputs: map[string][]float64{
						"out": repeatControl([]float64{0.9999999933746846, 0.24982677324761315}),
					},
				},
			},
		},
		{
			unit: "midi-hz",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in": []float64{60, 24},
					},
					outputs: map[string][]float64{
						"out": []float64{0.005932552501147361, 0.0007415690626434202},
					},
				},
			},
		},
		{
			unit: "pitch",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"class":  []float64{0, 7},
						"octave": []float64{0, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{0.0006606629273215556, A4},
					},
				},
			},
		},
		{
			unit: "smooth",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":   []float64{1, 2, 2, 2},
						"time": []float64{1, 3, 3, 3},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 1.3333333333333335, 1.5555555555555558, 1.7037037037037037},
					},
				},
			},
		},
		{
			unit: "gate-mix",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"0": []float64{-1, -1, 1, 1},
						"1": []float64{-1, 1, 1, -1},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, 1, 1, 1},
					},
				},
			},
		},
		{
			unit: "lerp",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":    []float64{0, 1, 0.5, 0.25},
						"min":   []float64{0, 1, 0, 0},
						"max":   []float64{0, 1, 2, 2},
						"scale": []float64{0, 1, 1, 4},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 1, 1, 2},
					},
				},
			},
		},
		{
			// TODO: This just checks for explosions. Find a better way to test this monster.
			unit: "reverb",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"a": []float64{0, 1},
						"b": []float64{0, 1},
					},
					outputs: map[string][]float64{
						"a": []float64{0, 1},
						"b": []float64{0, 1},
					},
				},
			},
		},
		{
			unit: "bipolar",
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
			unit: "unipolar",
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
		{
			unit: "logic",
			scenario: []scenario{
				{
					description: "OR/AND",
					inputs: map[string][]float64{
						"x":    []float64{1, 0, 1, 1},
						"y":    []float64{0, 1, 1, 0},
						"mode": []float64{0, 0, 1, 1},
					},
					outputs: map[string][]float64{
						"out": []float64{1, 1, 1, -1},
					},
				},
				{
					description: "XOR/NOR",
					inputs: map[string][]float64{
						"x":    []float64{1, 1, 1, 0},
						"y":    []float64{0, 1, 0, 0},
						"mode": []float64{2, 2, 3, 3},
					},
					outputs: map[string][]float64{
						"out": []float64{1, -1, -1, 1},
					},
				},
				{
					description: "NAND/XNOR",
					inputs: map[string][]float64{
						"x":    []float64{1, 0, 1, 0},
						"y":    []float64{1, 0, 0, 0},
						"mode": []float64{4, 4, 5, 5},
					},
					outputs: map[string][]float64{
						"out": []float64{-1, 1, -1, 1},
					},
				},
			},
		},
		{
			unit: "center",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in": []float64{0, 1, -1},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 1, -1.005},
					},
				},
			},
		},
		{
			unit:         "gate-series",
			configValues: map[string]interface{}{"size": 2},
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"clock":   []float64{-1, 1, -1, 1},
						"advance": []float64{-1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"0": []float64{-1, 1, -1, -1},
						"1": []float64{-1, -1, -1, 1},
					},
				},
			},
		},
		{
			unit: "dynamics",
			scenario: []scenario{
				{
					inputs: map[string][]float64{
						"in":      []float64{1, 1, 1, 1},
						"clamp":   []float64{1, 1, 1, 1},
						"relax":   []float64{1, 1, 1, 1},
						"control": []float64{0, 0, 0, 0},
					},
					outputs: map[string][]float64{
						"out": []float64{0, 0.00390625, 0.0077777099609375, 0.011614613437652587},
					},
				},
			},
		},
		{
			unit: "random-series",
			scenario: []scenario{
				{
					description: "unlocked",
					inputs: map[string][]float64{
						"clock": []float64{
							-1, 1, -1, 1,
							-1, 1, -1, 1,
						},
						"length": []float64{
							2, 2, 2, 2,
							2, 2, 2, 2,
						},
					},
					outputs: map[string][]float64{
						"gate": []float64{
							-1, -1, -1, 1,
							1, 1, 1, -1,
						},
						"value": []float64{
							0, 0, 0, 0.9405090880450124,
							0.9405090880450124, 0.4377141871869802, 0.4377141871869802, 0.6868230728671094,
						},
					},
				},
				{
					description: "partially locked",
					inputs: map[string][]float64{
						"clock": []float64{
							-1, 1, -1, 1,
							-1, 1, -1, 1,
						},
						"length": []float64{
							2, 2, 2, 2,
							2, 2, 2, 2,
						},
						"lock": []float64{
							0.5, 0.5, 0.5, 0.5,
							0.5, 0.5, 0.5, 0.5,
						},
					},
					outputs: map[string][]float64{
						"gate": []float64{
							-1, -1, -1, 1,
							1, 1, 1, 1,
						},
						"value": []float64{
							0, 0, 0, 0.9405090880450124,
							0.9405090880450124, 0.4377141871869802, 0.4377141871869802, 0.9405090880450124,
						},
					},
				},
			},
		},
		{
			unit:         "stages",
			configValues: map[string]interface{}{"size": 3},
			scenario: []scenario{
				{
					description: "first gate mode + data",
					inputs: map[string][]float64{
						"0/pulses": []float64{2, 2, 2, 2, 2, 2},
						"0/mode":   []float64{1, 1, 1, 1, 1, 1},
						"0/data":   []float64{100, 100, 100, 100, 100, 100},
						"1/pulses": []float64{1, 1, 1, 1, 1, 1},
						"1/mode":   []float64{1, 1, 1, 1, 1, 1},
						"1/data":   []float64{200, 200, 200, 200, 200, 200},
						"clock":    []float64{-1, 1, -1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"gate": []float64{-1, 1, -1, -1, -1, -1},
						"data": []float64{100, 100, 100, 100, 100, 200},
					},
				},
				{
					description: "first gate mode + frequency",
					inputs: map[string][]float64{
						"0/pulses": []float64{2, 2, 2, 2, 2, 2},
						"0/mode":   []float64{1, 1, 1, 1, 1, 1},
						"0/freq":   []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
						"1/pulses": []float64{1, 1, 1, 1, 1, 1},
						"1/mode":   []float64{1, 1, 1, 1, 1, 1},
						"1/freq":   []float64{0.2, 0.2, 0.2, 0.2, 0.2, 0.2},
						"clock":    []float64{-1, 1, -1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"gate": []float64{-1, 1, -1, -1, -1, -1},
						"freq": []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.2},
					},
				},
				{
					description: "all gate mode",
					inputs: map[string][]float64{
						"0/pulses": []float64{2, 2, 2, 2, 2, 2},
						"0/mode":   []float64{3, 3, 3, 3, 3, 3},
						"0/data":   []float64{100, 100, 100, 100, 100, 100},
						"1/pulses": []float64{1, 1, 1, 1, 1, 1},
						"1/mode":   []float64{1, 1, 1, 1, 1, 1},
						"1/data":   []float64{200, 200, 200, 200, 200, 200},
						"clock":    []float64{-1, 1, -1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"gate": []float64{-1, 1, -1, 1, -1, -1},
						"data": []float64{100, 100, 100, 100, 100, 200},
					},
				},
				{
					description: "last gate mode",
					inputs: map[string][]float64{
						"0/pulses": []float64{2, 2, 2, 2, 2, 2},
						"0/mode":   []float64{2, 2, 2, 2, 2, 2},
						"0/data":   []float64{100, 100, 100, 100, 100, 100},
						"1/pulses": []float64{1, 1, 1, 1, 1, 1},
						"1/mode":   []float64{1, 1, 1, 1, 1, 1},
						"1/data":   []float64{200, 200, 200, 200, 200, 200},
						"clock":    []float64{-1, 1, -1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"gate": []float64{-1, -1, -1, 1, -1, -1},
						"data": []float64{100, 100, 100, 100, 100, 200},
					},
				},
				{
					description: "hold gate mode",
					inputs: map[string][]float64{
						"0/pulses": []float64{2, 2, 2, 2, 2, 2},
						"0/mode":   []float64{4, 4, 4, 4, 4, 4},
						"0/data":   []float64{100, 100, 100, 100, 100, 100},
						"1/pulses": []float64{1, 1, 1, 1, 1, 1},
						"1/mode":   []float64{1, 1, 1, 1, 1, 1},
						"1/data":   []float64{200, 200, 200, 200, 200, 200},
						"clock":    []float64{-1, 1, -1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"gate": []float64{-1, 1, 1, 1, 1, -1},
						"data": []float64{100, 100, 100, 100, 100, 200},
					},
				},
				{
					description: "all gate mode + reverse",
					inputs: map[string][]float64{
						"mode":     []float64{1, 1, 1, 1, 1, 1},
						"0/pulses": []float64{2, 2, 2, 2, 2, 2},
						"0/mode":   []float64{3, 3, 3, 3, 3, 3},
						"0/data":   []float64{100, 100, 100, 100, 100, 100},
						"1/pulses": []float64{1, 1, 1, 1, 1, 1},
						"1/mode":   []float64{1, 1, 1, 1, 1, 1},
						"1/data":   []float64{200, 200, 200, 200, 200, 200},
						"2/pulses": []float64{1, 1, 1, 1, 1, 1},
						"2/mode":   []float64{1, 1, 1, 1, 1, 1},
						"2/data":   []float64{300, 300, 300, 300, 300, 300},
						"clock":    []float64{-1, 1, -1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"gate": []float64{-1, 1, -1, 1, -1, -1},
						"data": []float64{100, 100, 100, 100, 100, 300},
					},
				},
				{
					description: "ping pong",
					inputs: map[string][]float64{
						"mode":     []float64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
						"0/pulses": []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
						"0/mode":   []float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
						"0/data":   []float64{100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
						"1/pulses": []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
						"1/mode":   []float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
						"1/data":   []float64{200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200},
						"2/pulses": []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
						"2/mode":   []float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
						"2/data":   []float64{300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300},
						"clock":    []float64{-1, 1, -1, 1, -1, 1, -1, 1, -1, 1, -1, 1},
					},
					outputs: map[string][]float64{
						"gate": []float64{-1, 1, -1, -1, 1, -1, 1, -1, 1, -1, 1, -1},
						"data": []float64{100, 100, 100, 200, 200, 300, 300, 200, 200, 100, 100, 200},
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
				rand.Seed(1)
				builder := builders[test.unit]
				u, err := builder(Config{
					Values:     test.configValues,
					SampleRate: sampleRate,
					FrameSize:  frameSize,
				})
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
			require.Equal(t, v, u.Out[name].Out().Read(i), fmt.Sprintf("scenario %d -> output %q -> sample %d", index, name, i))
		}
	}
}

func repeatControl(s []float64) []float64 {
	var ss []float64
	for _, v := range s {
		for i := 0; i < controlPeriod; i++ {
			ss = append(ss, v)
		}
	}
	return ss
}

func padControl(s []float64) []float64 {
	var ss []float64
	for _, v := range s {
		ss = append(ss, v)
		for i := 0; i < controlPeriod-1; i++ {
			ss = append(ss, 0)
		}
	}
	return ss
}
