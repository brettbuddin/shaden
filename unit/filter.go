package unit

import (
	"github.com/brettbuddin/shaden/dsp"
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
	} else if config.Poles > 8 {
		config.Poles = 8
	}

	return NewUnit(io, &filter{
		filter: &dsp.SVFilter{Poles: config.Poles},
		in:     io.NewIn("in", dsp.Float64(0)),
		cutoff: io.NewIn("cutoff", dsp.Frequency(1000, c.SampleRate)),
		res:    io.NewIn("res", dsp.Float64(1)),
		poles:  io.NewIn("poles", dsp.Float64(config.Poles)),
		lp:     io.NewOut("lp"),
		bp:     io.NewOut("bp"),
		hp:     io.NewOut("hp"),
	}), nil
}

type filter struct {
	in, cutoff, res, poles *In
	lp, bp, hp             *Out
	filter                 *dsp.SVFilter
}

func (f *filter) ProcessSample(i int) {
	f.filter.Poles = f.poles.ReadSlowInt(i, clampInt(1, 8))
	f.filter.Cutoff = f.cutoff.ReadSlow(i, ident)
	f.filter.Resonance = f.res.ReadSlow(i, ident)
	lp, bp, hp := f.filter.Tick(f.in.Read(i))
	f.lp.Write(i, lp)
	f.bp.Write(i, bp)
	f.hp.Write(i, hp)
}
