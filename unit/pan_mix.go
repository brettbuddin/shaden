package unit

import (
	"fmt"

	"buddin.us/lumen/dsp"
)

func newPanMix(name string, c Config) (*Unit, error) {
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
		io     = NewIO()
		inputs = make([]*In, config.Size)
		levels = make([]*In, config.Size)
		pans   = make([]*In, config.Size)
	)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = io.NewIn(fmt.Sprintf("%d/in", i+1), dsp.Float64(0))
		levels[i] = io.NewIn(fmt.Sprintf("%d/level", i+1), dsp.Float64(1))
		pans[i] = io.NewIn(fmt.Sprintf("%d/pan", i+1), dsp.Float64(0))
	}

	return NewUnit(io, name, &panMix{
		master: io.NewIn("master", dsp.Float64(1)),
		a:      io.NewOut("a"),
		b:      io.NewOut("b"),
		inputs: inputs,
		levels: levels,
		pans:   pans,
	}), nil
}

type panMix struct {
	inputs, levels, pans []*In
	master               *In
	a, b                 *Out

	size int
}

func (m *panMix) ProcessSample(i int) {
	master := dsp.Clamp(m.master.Read(i), 0, 1)
	var a, b float64
	for j := 0; j < len(m.inputs); j++ {
		in := m.inputs[j].Read(i) * m.levels[j].ReadSlow(i, ident)
		aPan, bPan := dsp.PanMix(m.pans[j].ReadSlow(i, ident), in, in)
		a += aPan
		b += bPan
	}
	m.a.Write(i, a*master)
	m.b.Write(i, b*master)
}
