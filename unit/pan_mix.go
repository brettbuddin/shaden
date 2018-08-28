package unit

import (
	"fmt"
	"math"

	"github.com/brettbuddin/shaden/dsp"
)

func newPanMix(io *IO, c Config) (*Unit, error) {
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
		pans   = make([]*In, config.Size)
	)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = io.NewIn(fmt.Sprintf("%d/in", i), dsp.Float64(0))
		levels[i] = io.NewIn(fmt.Sprintf("%d/level", i), dsp.Float64(1))
		pans[i] = io.NewIn(fmt.Sprintf("%d/pan", i), dsp.Float64(0))
	}

	return NewUnit(io, &panMix{
		master: io.NewIn("master", dsp.Float64(1)),
		mode:   io.NewIn("mode", dsp.Float64(0)),
		a:      io.NewOut("a"),
		b:      io.NewOut("b"),
		inputs: inputs,
		levels: levels,
		pans:   pans,
	}), nil
}

type panMix struct {
	inputs, levels, pans []*In
	master, mode         *In
	a, b                 *Out
}

func (m *panMix) ProcessSample(i int) {
	var (
		master = dsp.Clamp(m.master.Read(i), 0, 1)
		mode   = dsp.Clamp(m.mode.Read(i), 0, 1)

		a, b, aInUse, bInUse float64
	)

	if mode > modeSum {
		mode = modeSum
	} else if mode < modeAverage {
		mode = modeAverage
	}

	for j := 0; j < len(m.inputs); j++ {
		var (
			in         = m.inputs[j].Read(i) * m.levels[j].ReadSlow(i, ident)
			aPan, bPan = dsp.PanMix(m.pans[j].ReadSlow(i, ident), in, in)
		)

		a += aPan
		b += bPan

		if mode == modeAverage {
			if aPan != 0 {
				aInUse++
			}
			if bPan != 0 {
				bInUse++
			}
		}
	}

	a *= master
	b *= master

	switch mode {
	case modeAverage:
		m.a.Write(i, a/math.Max(aInUse, 1))
		m.b.Write(i, b/math.Max(bInUse, 1))
	case modeSum:
		m.a.Write(i, a)
		m.b.Write(i, b)
	}
}
