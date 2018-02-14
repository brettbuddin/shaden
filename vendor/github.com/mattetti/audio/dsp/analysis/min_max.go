package analysis

import "github.com/mattetti/audio"

// MinMaxFloat returns the smallest and biggest samples in the buffer
func MinMaxFloat(buf *audio.PCMBuffer) (min, max float64) {
	if buf == nil || buf.Len() == 0 {
		return 0, 0
	}
	buf.SwitchPrimaryType(audio.Float)
	min = buf.Floats[0]

	for _, v := range buf.Floats {
		if v > max {
			max = v
		} else if v < min {
			min = v
		}
	}

	return min, max
}
