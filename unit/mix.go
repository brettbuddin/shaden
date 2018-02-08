package unit

import (
	"fmt"

	"buddin.us/shaden/dsp"
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
		master:      io.NewIn("master", dsp.Float64(1)),
		out:         io.NewOut("out"),
		inputs:      inputs,
		levels:      levels,
		levelValues: make([]float64, len(levels)),
	}), nil
}

type mix struct {
	inputs, levels []*In
	levelValues    []float64
	master         *In
	out            *Out
}

func (m *mix) ProcessSample(i int) {
	var sum float64
	for j := 0; j < len(m.inputs); j++ {
		sum += m.inputs[j].Read(i) * m.levels[j].ReadSlow(i, ident)
	}
	master := dsp.Clamp(m.master.Read(i), 0, 1)
	m.out.Write(i, sum*master)
}
