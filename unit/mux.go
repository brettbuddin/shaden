package unit

import (
	"fmt"

	"buddin.us/shaden/dsp"
)

func newMux(name string, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 2
	}

	io := NewIO()

	inputs := make([]*In, config.Size)
	for i := range inputs {
		inputs[i] = io.NewIn(fmt.Sprintf("%d", i+1), dsp.Float64(0))
	}

	return NewUnit(io, name, &mux{
		selection: io.NewIn("select", dsp.Float64(1)),
		out:       io.NewOut("out"),
		inputs:    inputs,
	}), nil
}

type mux struct {
	inputs    []*In
	selection *In
	out       *Out
}

func (m *mux) ProcessAudio(n int) {
	for i := 0; i < n; i++ {
		m.ProcessSample(i)
	}
}

func (m *mux) ProcessSample(i int) {
	max := float64(len(m.inputs) - 1)
	s := int(dsp.Clamp(m.selection.Read(i), 0, max))
	m.out.Write(i, m.inputs[s].Read(i))
}

func (m *mux) ProcessControl() {
	var (
		max = float64(len(m.inputs) - 1)
		s   = int(dsp.Clamp(m.selection.Read(0), 0, max))
		in  = m.inputs[s].Read(0)
	)
	m.out.Write(0, in)
}
