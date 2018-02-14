package transforms

import "github.com/mattetti/audio"

// PCMScale converts a buffer with audio content from -1 to 1 into
// the PCM scale based on the buffer's bitdepth.
func PCMScale(buf *audio.PCMBuffer) error {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}
	buf.SwitchPrimaryType(audio.Float)
	factor := float64(audio.IntMaxSignedValue(buf.Format.BitDepth))
	for i := 0; i < buf.Len(); i++ {
		buf.Floats[i] *= factor
	}

	return nil
}
