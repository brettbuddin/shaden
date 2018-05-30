package unit

import (
	"github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/shaden/dsp"
)

func newPitch(io *IO, c Config) (*Unit, error) {
	var (
		pitches = map[int]float64{}
		p       = musictheory.NewPitch(musictheory.C, musictheory.Natural, 0)
	)

	for i := 0; p.Octaves != 9; i++ {
		pitches[i] = dsp.Frequency(p.Freq(), c.SampleRate).Float64()
		p = p.Transpose(musictheory.Minor(2))
	}

	return NewUnit(io, &pitch{
		pitches: pitches,
		class:   io.NewIn("class", dsp.Float64(0)),
		octave:  io.NewIn("octave", dsp.Float64(4)),
		out:     io.NewOut("out"),
	}), nil
}

type pitch struct {
	class, octave *In
	out           *Out
	pitches       map[int]float64
}

func (p *pitch) ProcessSample(i int) {
	var (
		class  = dsp.Clamp(p.class.Read(i), 0, 12)
		octave = dsp.Clamp(p.octave.Read(i), 0, 8) + 1
		idx    = int(octave*10 + class)
	)
	p.out.Write(i, p.pitches[idx])
}
