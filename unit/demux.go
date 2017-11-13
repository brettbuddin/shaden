package unit

import (
	"fmt"

	"buddin.us/lumen/dsp"
)

func newDemux(name string, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 2
	}

	var (
		io   = NewIO()
		outs = make([]*Out, config.Size)
	)
	for i := 0; i < config.Size; i++ {
		outs[i] = io.NewOut(fmt.Sprintf("%d", i))
	}

	return NewUnit(io, name, &demux{
		in:        io.NewIn("in", dsp.Float64(0)),
		selection: io.NewIn("select", dsp.Float64(1)),
		outs:      outs,
	}), nil
}

type demux struct {
	in, selection *In
	outs          []*Out
}

func (d *demux) ProcessSample(i int) {
	for j := 0; j < len(d.outs); j++ {
		var (
			in  = d.in.Read(i)
			max = float64(len(d.outs) - 1)
			s   = int(dsp.Clamp(d.selection.Read(i), 0, max))
		)
		if j == s {
			d.outs[j].Write(i, in)
		} else {
			d.outs[j].Write(i, 0)
		}
	}
}
