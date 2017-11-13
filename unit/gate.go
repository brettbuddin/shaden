package unit

import "buddin.us/shaden/dsp"

func newGate(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &gate{
		filter:     &dsp.SVFilter{Poles: 4},
		in:         io.NewIn("in", dsp.Float64(0)),
		control:    io.NewIn("control", dsp.Float64(1)),
		mode:       io.NewIn("mode", dsp.Float64(1)),
		cutoffhigh: io.NewIn("cutoff-high", dsp.Frequency(20000)),
		cutofflow:  io.NewIn("cutoff-low", dsp.Frequency(0)),
		resonance:  io.NewIn("res", dsp.Float64(1)),
		aux:        io.NewIn("aux", dsp.Float64(0)),
		out:        io.NewOut("out"),
		sum:        io.NewOut("sum"),
	}), nil
}

const (
	gateModeLP int = iota
	gateModeCombo
	gateModeAmp
)

type gate struct {
	in, control, mode, cutoffhigh, cutofflow, resonance, aux *In
	out, sum                                                 *Out
	filter                                                   *dsp.SVFilter
}

func (g *gate) ProcessSample(i int) {
	var (
		control    = g.control.Read(i)
		cutoffhigh = g.cutoffhigh.Read(i)
		cutofflow  = g.cutofflow.Read(i)
		in         = g.in.Read(i)
		mode       = g.mode.Read(i)
		resonance  = g.resonance.Read(i)
		aux        = g.aux.Read(i)

		out float64
	)

	switch int(mode) {
	case gateModeLP:
		out = g.applyFilter(in, control, cutoffhigh, cutofflow, resonance)
	case gateModeCombo:
		out = g.applyFilter(in, control, cutoffhigh, cutofflow, resonance) * control
	case gateModeAmp:
		out = in * control
	default:
		out = g.applyFilter(in, control, cutoffhigh, cutofflow, resonance) * control
	}

	g.out.Write(i, out)
	g.sum.Write(i, out+aux)
}

func (g *gate) applyFilter(in, ctrl, cutoffhigh, cutofflow, res float64) float64 {
	g.filter.Cutoff = cutofflow + (cutoffhigh-cutofflow)*ctrl
	g.filter.Resonance = res
	lp, _, _ := g.filter.Tick(in)
	return lp
}
