package midi

import "buddin.us/shaden/unit"

type midiOutput interface {
	unit.Output
	unit.CondProcessor
	unit.FrameProcessor
	unit.SampleProcessor
}

// Ensure midi-in outputs conform to the interfaces necessary for processing.
var _ = []midiOutput{
	&pitch{},
	&pitchRaw{},
	&gate{},
	&bend{},
	&cc{},
}
