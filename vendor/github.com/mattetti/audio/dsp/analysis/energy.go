package analysis

import "github.com/mattetti/audio"

// TotalEnergy is the the sum of squared moduli
// See https://www.dsprelated.com/freebooks/mdft/Signal_Metrics.html
func TotalEnergy(buf *audio.PCMBuffer) float64 {
	var e float64
	for _, v := range buf.AsFloat64s() {
		e += v * v
	}
	return e
}
