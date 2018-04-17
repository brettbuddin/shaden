package unit

import (
	"fmt"
	"math"

	"buddin.us/shaden/dsp"
)

func newMorph(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 3
	}

	inputs := make([]*In, config.Size)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = io.NewIn(fmt.Sprintf("%d", i), dsp.Float64(0))
	}

	return NewUnit(io, &morph{
		morph:  io.NewIn("morph", dsp.Float64(0)),
		out:    io.NewOut("out"),
		inputs: inputs,
	}), nil
}

type morph struct {
	morph  *In
	inputs []*In
	out    *Out
}

func (m *morph) ProcessSample(i int) {
	var (
		n       = float64(len(m.inputs))
		morph   = dsp.Clamp(m.morph.Read(i), 0, 1)
		pos     = morph * (n - 1)
		begin   = int(math.Floor(pos))
		end     = int(math.Ceil(pos))
		beginIn = m.inputs[begin].Read(i)
		endIn   = m.inputs[end].Read(i)
		_, frac = math.Modf(pos)
	)

	m.out.Write(i, dsp.Lerp(beginIn, endIn, frac))
}
