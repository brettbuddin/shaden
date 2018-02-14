package transforms

import (
	"math"

	"github.com/mattetti/audio"
)

var (
	crusherStepSize  = 0.000001
	CrusherMinFactor = 1.0
	CrusherMaxFactor = 2097152.0
)

// BitCrush reduces the resolution of the sample to the target bit depth
// Note that bit crusher effects are usually made of this feature + a decimator
func BitCrush(buf *audio.PCMBuffer, factor float64) {
	buf.SwitchPrimaryType(audio.Float)
	stepSize := crusherStepSize * factor
	for i := 0; i < len(buf.Floats); i++ {
		frac, exp := math.Frexp(buf.Floats[i])
		frac = signum(frac) * math.Floor(math.Abs(frac)/stepSize+0.5) * stepSize
		buf.Floats[i] = math.Ldexp(frac, exp)
	}
}

func signum(v float64) float64 {
	if v >= 0.0 {
		return 1.0
	}
	return -1.0
}
