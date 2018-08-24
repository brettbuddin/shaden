package unit

import (
	"fmt"

	"github.com/brettbuddin/shaden/dsp"
)

func newFilterBank(io *IO, c Config) (*Unit, error) {
	var config struct {
		Size int
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.Size == 0 {
		config.Size = 5
	}

	var (
		filters    = make([]*dsp.SVFilter, config.Size)
		cutoffs    = make([]*In, config.Size)
		resonances = make([]*In, config.Size)
		freq       = 300.0
	)

	for i := range filters {
		var (
			cutoff = dsp.Frequency(freq, c.SampleRate)
			res    = 1.0
		)
		filters[i] = &dsp.SVFilter{
			Cutoff:    cutoff.Float64(),
			Poles:     4,
			Resonance: res,
		}
		cutoffs[i] = io.NewIn(fmt.Sprintf("%d/cutoff", i), cutoff)
		resonances[i] = io.NewIn(fmt.Sprintf("%d/res", i), dsp.Float64(res))
		freq += 600.0
	}

	return NewUnit(io, &filterBank{
		filters:    filters,
		in:         io.NewIn("in", dsp.Float64(0)),
		cutoffs:    cutoffs,
		resonances: resonances,
		out:        io.NewOut("out"),
	}), nil
}

type filterBank struct {
	in                  *In
	cutoffs, resonances []*In
	out                 *Out

	filters []*dsp.SVFilter
}

func (f *filterBank) ProcessSample(i int) {
	var (
		in   = f.in.Read(i)
		size = len(f.filters)

		out, sum float64
	)

	for j, filter := range f.filters {
		filter.Cutoff = f.cutoffs[j].Read(i)
		filter.Resonance = f.resonances[j].Read(i)

		if j == 0 {
			out, _, _ = filter.Tick(in)
		} else if j < size-1 {
			_, out, _ = filter.Tick(in)
		} else {
			_, _, out = filter.Tick(in)
		}
		sum += out
	}

	f.out.Write(i, sum/float64(len(f.filters)))
}
