package portaudio_test

import (
	"buddin.us/shaden/engine"
	"buddin.us/shaden/engine/portaudio"
)

var _ engine.Backend = &portaudio.PortAudio{}
