package dsp

// NewAllPassMS returns a new AllPass
func NewAllPassMS(ms MS) *AllPass {
	return &AllPass{dl: NewDelayLineMS(ms)}
}

// NewAllPass returns a new AllPass
func NewAllPass(size int) *AllPass {
	return &AllPass{dl: NewDelayLine(size)}
}

// AllPass is an allpass filter
type AllPass struct {
	dl   *DelayLine
	last float64
}

// Tick advances the allpass filter state using the full length of the delay line
func (a *AllPass) Tick(in, gain float64) float64 {
	return a.TickAbsolute(in, gain, -1)
}

// TickAbsolute advances the allpass filter state with a specified absolute delay length (less than maximum)
func (a *AllPass) TickAbsolute(in, gain, delay float64) float64 {
	before := in + -gain*a.last
	a.last = tick(a.dl, before, delay)
	return a.last + gain*before
}

// TickRelative advances the allpass filter state with a delay length relative to the length of the delay line
func (a *AllPass) TickRelative(in, gain, scale float64) float64 {
	return a.TickAbsolute(in, gain, float64(len(a.dl.buffer))*scale)
}
