package transforms

import "github.com/mattetti/audio"

// MonoDownmix converts the buffer to a mono buffer
// by downmixing the channels together.
func MonoDownmix(buf *audio.PCMBuffer) error {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}
	nChans := buf.Format.NumChannels
	if nChans < 2 {
		return nil
	}
	nChansF := float64(nChans)

	frameCount := buf.Size()
	newData := make([]float64, frameCount)
	buf.SwitchPrimaryType(audio.Float)
	for i := 0; i < frameCount; i++ {
		newData[i] = 0
		for j := 0; j < nChans; j++ {
			newData[i] += buf.Floats[i*nChans+j]
		}
		newData[i] /= nChansF
	}
	buf.Floats = newData
	buf.Format.NumChannels = 1

	return nil
}
