package dsp

// NewFBComb returns a new FBComb
func NewFBComb(size int) *FBComb {
	return &FBComb{dl: NewDelayLine(size)}
}

// NewFBCombMS returns a new FBComb that's length is represented in milliseconds
func NewFBCombMS(ms MS) *FBComb {
	return &FBComb{dl: NewDelayLineMS(ms)}
}

// FBComb is a feedback comb filter
type FBComb struct {
	dl   *DelayLine
	last float64
}

// Tick advances the filter's state with the default delay
func (c *FBComb) Tick(in, gain float64) float64 {
	return c.TickAbsolute(in, gain, -1)
}

// TickAbsolute advances the filter's state with a specific delay
func (c *FBComb) TickAbsolute(in, gain, delay float64) float64 {
	out := in + c.last
	c.last = gain * tick(c.dl, in, delay)
	return out
}

// NewFFComb returns a new FBComb
func NewFFComb(size int) *FFComb {
	return &FFComb{dl: NewDelayLine(size)}
}

// NewFFCombMS returns a new FFComb that's length is represented in milliseconds
func NewFFCombMS(ms MS) *FFComb {
	return &FFComb{dl: NewDelayLineMS(ms)}
}

// FFComb is a feedforward comb filter
type FFComb struct {
	dl *DelayLine
}

// Tick advances the filter's state with the default delay
func (c *FFComb) Tick(in, gain float64) float64 {
	return c.TickAbsolute(in, gain, -1)
}

// TickAbsolute advances the filter's state with a specific delay
func (c *FFComb) TickAbsolute(in, gain, delay float64) float64 {
	return in + gain*c.dl.TickAbsolute(in, delay)
}

func tick(dl *DelayLine, in, delay float64) float64 {
	if delay < 0 {
		return dl.Tick(in)
	}
	return dl.TickAbsolute(in, delay)
}
