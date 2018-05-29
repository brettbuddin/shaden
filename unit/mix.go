package unit

import (
	"fmt"

	"github.com/brettbuddin/shaden/dsp"
)

const (
	modeSum = iota
	modeAverage
)

func newMix(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 4
	}

	var (
		inputs = make([]*In, config.Size)
		levels = make([]*In, config.Size)
	)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = io.NewIn(fmt.Sprintf("%d/in", i), dsp.Float64(0))
		levels[i] = io.NewIn(fmt.Sprintf("%d/level", i), dsp.Float64(1))
	}

	return NewUnit(io, &mix{
		master: io.NewIn("master", dsp.Float64(1)),
		mode:   io.NewIn("mode", dsp.Float64(0)),
		out:    io.NewOut("out"),
		inputs: inputs,
		levels: levels,
	}), nil
}

type mix struct {
	inputs, levels []*In
	master, mode   *In
	out            *Out
}

func (m *mix) ProcessSample(i int) {
	var (
		master = m.master.ReadSlow(i, clamp(0, 1))
		mode   = m.mode.ReadSlowInt(i, identInt)
		final  float64
	)

	switch mode {
	case modeSum:
		final = m.sum(i)
	case modeAverage:
		var inUse float64
		for _, in := range m.inputs {
			if in.HasSource() {
				inUse++
			}
		}
		final = m.sum(i) / inUse
	}

	m.out.Write(i, final*master)
}

func (m *mix) sum(i int) float64 {
	var final float64
	for j := range m.inputs {
		final += m.inputs[j].Read(i) * m.levels[j].ReadSlow(i, ident)
	}
	return final
}
