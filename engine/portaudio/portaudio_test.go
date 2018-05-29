package portaudio_test

import (
	"github.com/brettbuddin/shaden/engine"
	"github.com/brettbuddin/shaden/engine/portaudio"
)

var _ engine.Backend = &portaudio.PortAudio{}
