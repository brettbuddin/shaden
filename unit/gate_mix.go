package unit

import (
	"fmt"

	"github.com/brettbuddin/shaden/dsp"
)

func newGateMix(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 4
	}

	inputs := make([]*In, config.Size)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = io.NewIn(fmt.Sprintf("%d", i), dsp.Float64(-1))
	}

	return NewUnit(io, &gateMix{
		out:    io.NewOut("out"),
		inputs: inputs,
	}), nil
}

type gateMix struct {
	inputs []*In
	out    *Out
}

func (m *gateMix) ProcessSample(i int) {
	var out float64 = -1
	for j := 0; j < len(m.inputs); j++ {
		if m.inputs[j].Read(i) > 0 {
			out = 1
			break
		}
	}
	m.out.Write(i, out)
}
