package unit

import "buddin.us/shaden/dsp"

const clusterSize = 12

func newCluster(io *IO, c Config) (*Unit, error) {
	return NewUnit(io, &cluster{
		freq:     io.NewIn("freq", dsp.Frequency(440, c.SampleRate)),
		interval: io.NewIn("interval", dsp.Float64(1.1)),
		phases:   make([]float64, clusterSize),
		out:      io.NewOut("out"),
	}), nil
}

type cluster struct {
	freq, interval *In
	out            *Out
	phases         []float64
}

func (s *cluster) ProcessSample(i int) {
	var (
		freq          = s.freq.Read(i)
		interval      = s.interval.Read(i)
		intervalAccum = 1.0
		out           float64
	)
	for i, phase := range s.phases {
		idx := float64(i)
		out += dsp.Sin(phase) * (clusterSize - idx) / clusterSize
		s.phases[i] += (freq * intervalAccum) * twoPi
		if s.phases[i] >= twoPi {
			s.phases[i] -= twoPi
		}
		intervalAccum *= interval
	}
	s.out.Write(i, out)
}
