package transforms

import (
	"math"

	"github.com/mattetti/audio"
)

// FullWaveRectifier to make all signal positive
// See https://en.wikipedia.org/wiki/Rectifier#Full-wave_rectification
func FullWaveRectifier(buf *audio.PCMBuffer) error {
	if buf == nil {
		return audio.ErrInvalidBuffer
	}
	buf.SwitchPrimaryType(audio.Float)
	for i := 0; i < buf.Len(); i++ {
		buf.Floats[i] = math.Abs(buf.Floats[i])
	}

	return nil
}

// MonoRMS converts the buffer to mono and apply an RMS treatment.
// rms = sqrt ( (1/n) * (x12 + x22 + … + xn2) )
// multiplying by 1/n effectively assigns equal weights to all the terms, making it a rectangular window.
// Other window equations can be used instead which would favor terms in the middle of the window.
// This results in even greater accuracy of the RMS value since brand new samples (or old ones at
// the end of the window) have less influence over the signal’s power.)
// TODO: use a sine wave at amplitude of 1: rectication + average = 1.4 (1/square root of 2)
func MonoRMS(b *audio.PCMBuffer, windowSize int) error {
	if b == nil {
		return audio.ErrInvalidBuffer
	}
	if b.Len() == 0 {
		return nil
	}
	b.SwitchPrimaryType(audio.Float)
	out := []float64{}
	winBuf := make([]float64, windowSize)
	windowSizeF := float64(windowSize)

	processWindow := func(idx int) {
		total := 0.0
		for i := 0; i < len(winBuf); i++ {
			total += winBuf[idx] * winBuf[idx]
		}
		v := math.Sqrt((1.0 / windowSizeF) * total)
		out = append(out, v)
	}

	nbrChans := b.Format.NumChannels
	samples := b.AsFloat64s()

	var windowIDX int
	// process each frame, convert it to mono and them RMS it
	for i := 0; i < len(samples); i++ {
		v := samples[i]
		if nbrChans > 1 {
			for j := 1; j < nbrChans; j++ {
				i++
				v += samples[i]
			}
			v /= float64(nbrChans)
		}
		winBuf[windowIDX] = v
		windowIDX++
		if windowIDX == windowSize || i == (len(samples)-1) {
			windowIDX = 0
			processWindow(windowIDX)
		}
	}

	b.Format.NumChannels = 1
	b.Format.SampleRate /= windowSize
	b.Floats = out
	return nil
}
