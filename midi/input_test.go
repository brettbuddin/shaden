package midi

import "buddin.us/shaden/unit"

type midiOutput interface {
	Out() *unit.Out
	ExternalNeighborCount() int
	unit.FrameProcessor
	unit.SampleProcessor
}

// Ensure midi-in outputs conform to the interfaces necessary for processing.
// TODO: Improve the behavior of thes OutputProcessors so this type of test isn't as necessary.
var _ = []midiOutput{
	&pitch{},
	&pitchRaw{},
	&gate{},
	&bend{},
	&cc{},
}
