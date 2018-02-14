package transforms

import (
	"github.com/mattetti/audio"
	"github.com/mattetti/audio/dsp/analysis"
)

// NominalScale converts the input to a -1.0 / +1.0 scale
func NominalScale(buf *audio.PCMBuffer) error {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}
	buf.SwitchPrimaryType(audio.Float)
	min, max := analysis.MinMaxFloat(buf)
	// check if already in the right scale
	if min >= -1 && max <= 1 {
		return nil
	}

	max = float64(audio.IntMaxSignedValue(buf.Format.BitDepth))
	for i := 0; i < buf.Len(); i++ {
		buf.Floats[i] /= max
	}

	return nil
}
