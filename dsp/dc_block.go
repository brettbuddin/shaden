package dsp

// DCBlock keeps a signal centered around zero
type DCBlock struct {
	lastIn, lastOut float64
}

// Tick advances the state
func (dc *DCBlock) Tick(in float64) float64 {
	out := in - dc.lastIn + dc.lastOut*0.995
	dc.lastIn, dc.lastOut = in, out
	return out
}
