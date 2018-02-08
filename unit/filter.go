package unit

import (
	"buddin.us/shaden/dsp"
)

func newFilter(io *IO, c Config) (*Unit, error) {
	var config struct {
		Poles int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Poles == 0 {
		config.Poles = 4
	}

	return NewUnit(io, &filter{
		filter: &dsp.SVFilter{Poles: config.Poles},
		in:     io.NewIn("in", dsp.Float64(0)),
		cutoff: io.NewIn("cutoff", dsp.Frequency(1000)),
		res:    io.NewIn("res", dsp.Float64(1)),
		lp:     io.NewOut("lp"),
		bp:     io.NewOut("bp"),
		hp:     io.NewOut("hp"),
	}), nil
}

type filter struct {
	in, cutoff, res *In
	lp, bp, hp      *Out
	filter          *dsp.SVFilter
}

func (f *filter) ProcessSample(i int) {
	f.filter.Cutoff = f.cutoff.ReadSlow(i, ident)
	f.filter.Resonance = f.res.ReadSlow(i, ident)
	lp, bp, hp := f.filter.Tick(f.in.Read(i))
	f.lp.Write(i, lp)
	f.bp.Write(i, bp)
	f.hp.Write(i, hp)
}
