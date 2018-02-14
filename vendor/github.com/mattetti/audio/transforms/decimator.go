package transforms

import (
	"errors"

	"github.com/mattetti/audio"
	"github.com/mattetti/audio/transforms/filters"
)

// Decimate drops samples to switch to a lower sample rate.
// Factor is the decimation factor, for instance a factor of 2 of a 44100Hz buffer
// will convert the buffer in a 22500 buffer.
func Decimate(buf *audio.PCMBuffer, factor int) (err error) {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}

	if factor < 0 {
		return errors.New("can't use a negative factor")
	}

	// apply a low pass filter at the nyquist frequency to avoid
	// aliasing.
	if err := filters.LowPass(buf, float64(buf.Format.SampleRate)/2); err != nil {
		return err
	}

	// drop samples to match the decimation factor
	newLength := len(buf.Floats) / factor
	for i := 0; i < newLength; i++ {
		buf.Floats[i] = buf.Floats[i*factor]
	}
	buf.Floats = buf.Floats[:newLength]
	buf.Format.SampleRate /= factor

	return nil
}
