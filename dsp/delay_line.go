package dsp

import "math"

// DelayLine is a simple delay line
type DelayLine struct {
	buffer []float64
	sizeMS MS
	size   int
	offset int
}

// NewDelayLine returns a new DelayLine of a specific maximum size in milliseconds
func NewDelayLine(size int) *DelayLine {
	return &DelayLine{
		buffer: make([]float64, size),
	}
}

// NewDelayLineMS returns a new DelayLine of a specific maximum size in milliseconds
func NewDelayLineMS(size MS) *DelayLine {
	return &DelayLine{
		buffer: make([]float64, int(size.Float64())),
	}
}

// Size returns the size of the DelayLine buffer
func (d *DelayLine) Size() int {
	return len(d.buffer)
}

// Tick advances the state of the DelayLine using the total delay length
func (d *DelayLine) Tick(v float64) float64 {
	return d.TickAbsolute(v, float64(d.size))
}

// TickAbsolute advances the state of the DelayLine using a specific delay length (less than total) in samples
func (d *DelayLine) TickAbsolute(v, delay float64) float64 {
	d.Write(v)
	return d.ReadAbsolute(delay)
}

// TickRelative advances the state of the DelayLine using a scale between 0 and 1
func (d *DelayLine) TickRelative(v, scale float64) float64 {
	d.Write(v)
	return d.ReadRelative(scale)
}

// Write writes a value to the DelayLine and steps ahead by one sample
func (d *DelayLine) Write(v float64) {
	size := len(d.buffer)
	d.buffer[d.offset] = v
	d.offset = (d.offset - 1 + size) % size
}

// ReadAbsolute reads a specific sample offset from the DelayLine
func (d *DelayLine) ReadAbsolute(pos float64) float64 {
	durationI, durationF := math.Modf(pos)
	size := len(d.buffer)
	offset := d.offset + int(durationI)
	a := d.buffer[(offset+size)%size]
	b := d.buffer[(offset+1+size)%size]
	return Lerp(a, b, durationF)
}

// ReadRelative reads sample from the DelayLine at a scale between 0 and 1
func (d *DelayLine) ReadRelative(scale float64) float64 {
	return d.ReadAbsolute(scale * float64(len(d.buffer)-1))
}
