package unit

import (
	"fmt"
	"math"

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

		sum, inUse float64
	)

	if mode > modeAverage {
		mode = modeAverage
	} else if mode < modeSum {
		mode = modeSum
	}

	for j := range m.inputs {
		var (
			in  = m.inputs[j].Read(i)
			lvl = m.levels[j].ReadSlow(i, ident)
		)

		if mode == modeAverage && in != 0 {
			inUse++
		}

		sum += in * lvl
	}

	sum *= master

	switch mode {
	case modeAverage:
		m.out.Write(i, sum/math.Max(1, inUse))
	case modeSum:
		m.out.Write(i, sum)
	}
}
