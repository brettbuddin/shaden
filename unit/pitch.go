package unit

import (
	"buddin.us/shaden/dsp"
	"buddin.us/musictheory"
)

var pitches = map[int]float64{}

func init() {
	p := musictheory.NewPitch(musictheory.C, musictheory.Natural, 0)
	for i := 0; p.Octaves != 9; i++ {
		pitches[i] = dsp.Frequency(p.Freq()).Float64()
		p = p.Transpose(musictheory.Minor(2)).(musictheory.Pitch)
	}
}

func newPitch(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &pitch{
		class:  io.NewIn("class", dsp.Float64(0)),
		octave: io.NewIn("octave", dsp.Float64(4)),
		out:    io.NewOut("out"),
	}), nil
}

type pitch struct {
	class, octave         *In
	out                   *Out
	lastClass, lastOctave float64
	freq                  float64
}

func (p *pitch) ProcessSample(i int) {
	var (
		class  = dsp.Clamp(p.class.Read(i), 0, 12)
		octave = dsp.Clamp(p.octave.Read(i), 0, 8)
		idx    = int(octave*10 + class)
	)
	p.out.Write(i, pitches[idx])
}
