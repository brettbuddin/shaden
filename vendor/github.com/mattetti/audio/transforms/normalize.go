package transforms

import (
	"math"

	"github.com/mattetti/audio"
)

// NormalizeMax sets the max value to 1 and normalize the rest of the data.
func NormalizeMax(buf *audio.PCMBuffer) {
	if buf == nil {
		return
	}
	buf.SwitchPrimaryType(audio.Float)
	max := 0.0

	for i := 0; i < buf.Len(); i++ {
		if math.Abs(buf.Floats[i]) > max {
			max = math.Abs(buf.Floats[i])
		}
	}

	if max != 0.0 {
		for i := 0; i < buf.Len(); i++ {
			buf.Floats[i] /= max
		}
	}
}
