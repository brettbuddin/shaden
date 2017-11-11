package engine

import (
	"buddin.us/lumen/dsp"
	"buddin.us/lumen/unit"
)

func unitBuilders(e *Engine) map[string]unit.BuildFunc {
	return map[string]unit.BuildFunc{
		"source": newSource(e),
	}
}

func newSink() *unit.Unit {
	io := unit.NewIO()
	io.NewIn("l", dsp.Float64(0))
	io.NewIn("r", dsp.Float64(0))
	return unit.NewUnit(io, "Sink", nil)
}

func newSource(e *Engine) unit.BuildFunc {
	return func(unit.Config) (*unit.Unit, error) {
		io := unit.NewIO()
		io.NewOutWithFrame("output", e.input)
		return unit.NewUnit(io, "source", nil), nil
	}
}
