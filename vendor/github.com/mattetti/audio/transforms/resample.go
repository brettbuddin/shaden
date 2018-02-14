package transforms

import (
	"math"

	"github.com/mattetti/audio"
	"github.com/mattetti/audio/transforms/filters"
)

// Resample down or up samples the buffer. Note that the amplitude
// will be affected by upsampling.
func Resample(buf *audio.PCMBuffer, fs float64) error {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}
	if buf.Format.SampleRate == int(fs) {
		return nil
	}
	buf.SwitchPrimaryType(audio.Float)

	// downsample
	if fs < float64(buf.Format.SampleRate) {
		factor := float64(buf.Format.SampleRate) / fs

		// apply a low pass filter at the nyquist frequency to avoid
		// aliasing.
		if err := filters.LowPass(buf, float64(buf.Format.SampleRate)/2); err != nil {
			return err
		}

		// drop samples to match the decimation factor
		newLength := int(math.Floor(float64(len(buf.Floats)) / factor))
		var targetI int
		for i := 0; i < newLength; i++ {
			targetI = int(math.Floor(float64(i) * factor))
			buf.Floats[i] = buf.Floats[targetI]
		}
		buf.Floats = buf.Floats[:newLength]
		buf.Format.SampleRate = int(fs)
		return nil
	}

	// oversample
	// Note: oversampling reduces the amplitude
	factor := fs / float64(buf.Format.SampleRate)
	newLength := int(math.Ceil(float64(len(buf.Floats)) * factor))
	newFloats := make([]float64, newLength)
	padding := int(
		math.Ceil(
			float64(newLength) /
				float64(len(buf.Floats))))
	var idx int
	for i := 0; i < len(buf.Floats); i++ {
		idx = i * padding
		if idx >= len(newFloats) {
			break
		}
		newFloats[idx] = buf.Floats[i]
	}
	buf.Floats = newFloats
	buf.Format.SampleRate = int(fs)
	return nil
}
