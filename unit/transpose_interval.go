package unit

import (
	"github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/shaden/dsp"
)

type quality int

const (
	qualityPerfect quality = iota
	qualityMinor
	qualityMajor
	qualityDiminished
	qualityAugmented
)

var perfectFirst = musictheory.Perfect(1)

func newTransposeInterval(io *IO, _ Config) (*Unit, error) {
	return NewUnit(io, &transposeInterval{
		in:       io.NewIn("in", dsp.Float64(0)),
		quality:  io.NewIn("quality", dsp.Float64(0)),
		step:     io.NewIn("step", dsp.Float64(0)),
		out:      io.NewOut("out"),
		interval: perfectFirst,
	}), nil
}

type transposeInterval struct {
	in, quality, step *In
	interval          musictheory.Interval
	out               *Out
}

func (t *transposeInterval) ProcessSample(i int) {
	var (
		in   = t.in.Read(i)
		qual = int(t.quality.Read(i))
		step = int(t.step.Read(i))
	)
	t.out.Write(i, t.calc(in, qual, step))
}

func (t *transposeInterval) calc(in float64, qual, step int) float64 {
	var interval musictheory.Interval

	switch quality(qual) {
	case qualityPerfect:
		interval = musictheory.Perfect(step)
	case qualityMinor:
		interval = musictheory.Minor(step)
	case qualityMajor:
		interval = musictheory.Major(step)
	case qualityDiminished:
		interval = musictheory.Diminished(step)
	case qualityAugmented:
		interval = musictheory.Augmented(step)
	}

	return in * interval.Ratio()
}
