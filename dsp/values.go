// Package dsp provides common digital signal processing operations.
package dsp

import (
	"fmt"

	"buddin.us/musictheory"
)

// Valuer is the wrapper interface around the Value method; which is used in obtaining the constant value
type Valuer interface {
	Float64() float64
}

// Float64 is a wrapper for float64 that implements Valuer
type Float64 float64

// Float64 returns the constant value
func (v Float64) Float64() float64 { return float64(v) }
func (v Float64) String() string   { return fmt.Sprintf("%.2f", v) }

// Hz represents cycles-per-second
type Hz struct {
	Valuer
	Raw float64
}

// Frequency returns a scalar value in Hz
func Frequency(v float64, sampleRate int) Hz {
	return Hz{Raw: v, Valuer: Float64(v / float64(sampleRate))}
}

// Float64 returns the constant value
func (hz Hz) Float64() float64 {
	if hz.Valuer == nil {
		return 0
	}
	return hz.Valuer.Float64()
}
func (hz Hz) String() string { return fmt.Sprintf("%.2fHz", hz.Raw) }

// ParsePitch parses the scientific notation of a pitch
func ParsePitch(v string, sampleRate int) (Pitch, error) {
	p, err := musictheory.ParsePitch(v)
	if err != nil {
		return Pitch{}, err
	}
	return Pitch{
		Valuer: Frequency(p.Freq(), sampleRate),
		Raw:    v,
	}, nil
}

// Pitch is a pitch that has been expressed in scientific notation
type Pitch struct {
	Valuer
	Raw string
}

// Float64 returns the constant value
func (p Pitch) Float64() float64 {
	if p.Valuer == nil {
		return 0
	}
	return p.Valuer.Float64()
}
func (p Pitch) String() string { return p.Raw }

// MS is a value representation of milliseconds
type MS struct {
	Valuer
	Raw float64
}

// DurationInt returns a scalar value (int) in MS
func DurationInt(v, sampleRate int) MS { return Duration(float64(v), sampleRate) }

// Duration returns a scalar value (float64) in MS
func Duration(v float64, sampleRate int) MS {
	return MS{
		Valuer: Float64(v * float64(sampleRate) * 0.001),
		Raw:    v,
	}
}

// Float64 returns the constant value
func (ms MS) Float64() float64 {
	if ms.Valuer == nil {
		return 0
	}
	return ms.Valuer.Float64()
}
func (ms MS) String() string { return fmt.Sprintf("%.2fms", ms.Raw) }

// BeatsPerMin represents beats-per-minute
type BeatsPerMin struct {
	Valuer
	Raw float64
}

// BPM returns a scalar value in beats-per-minute
func BPM(v float64, sampleRate int) BeatsPerMin {
	return BeatsPerMin{
		Valuer: Float64(v / 60 / float64(sampleRate)),
		Raw:    v,
	}
}

// Float64 returns the constant value
func (bpm BeatsPerMin) Float64() float64 {
	if bpm.Valuer == nil {
		return 0
	}
	return bpm.Valuer.Float64()
}
func (bpm BeatsPerMin) String() string { return fmt.Sprintf("%.2fBPM", bpm.Raw) }
