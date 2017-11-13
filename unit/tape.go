package unit

import (
	"buddin.us/shaden/dsp"
)

func newTape(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &tape{
		in:      io.NewIn("in", dsp.Float64(0)),
		mix:     io.NewIn("mix", dsp.Float64(-1)),
		speed:   io.NewIn("speed", dsp.Float64(1)),
		record:  io.NewIn("record", dsp.Float64(-1)),
		play:    io.NewIn("play", dsp.Float64(1)),
		out:     io.NewOut("out"),
		stateFn: tapePlay,
		state: &tapeState{
			memory: make([]float64, 5*dsp.SampleRate),
			clean:  true,
		},
	}), nil
}

type tape struct {
	in, mix, speed, record, play, pos *In
	state                             *tapeState
	stateFn                           tapeStateFunc
	out                               *Out
}

func (t *tape) ProcessSample(i int) {
	t.state.in = t.in.Read(i)
	t.state.play = t.play.Read(i)
	t.state.mix = dsp.Clamp(t.mix.Read(i), -1, 1)
	t.state.record = dsp.Clamp(t.record.Read(i), -1, 1)
	t.state.speed = dsp.Clamp(t.speed.Read(i), -1, 1)

	t.stateFn = t.stateFn(t.state)
	t.out.Write(i, t.state.out)

	t.state.last.record = t.state.record
}

type tapeState struct {
	in, mix, speed, record, play, out float64
	offset, recordEnd                 int
	partialOffset                     float64
	memory                            []float64
	clean                             bool
	last                              lastTapeState
}

type lastTapeState struct {
	mix, record, play float64
}

type tapeStateFunc func(*tapeState) tapeStateFunc

func tapePlay(s *tapeState) tapeStateFunc {
	if isLow(s.play) {
		return tapePlay
	}
	if isTrig(s.last.record, s.record) {
		return tapeRecord
	}
	playback(s)
	advance(s)
	return tapePlay
}

func tapeRecord(s *tapeState) tapeStateFunc {
	record(s)
	if isTrig(s.last.record, s.record) {
		if s.clean {
			s.recordEnd = s.offset
		}
		return tapePlay
	}
	playback(s)
	advance(s)
	return tapeRecord
}

func playback(s *tapeState) {
	s.out = dsp.Mix(s.mix, s.in, s.memory[s.offset])
}

func record(s *tapeState) {
	s.memory[s.offset] = dsp.Mix(s.mix, s.in, s.memory[s.offset])
}

func advance(s *tapeState) {
	if s.speed < 1 || s.speed > -1 {
		s.partialOffset += s.speed
		if s.partialOffset > 1 {
			s.partialOffset--
			s.offset++
		} else if s.partialOffset < 0 {
			s.partialOffset++
			s.offset--
		}
	} else if s.offset >= 1 {
		s.offset++
	} else if s.offset <= -1 {
		s.offset--
	}
	s.offset = (s.offset + 1) % len(s.memory)

	if s.recordEnd > 0 {
		if s.offset >= s.recordEnd {
			s.offset = 0
		} else if s.offset < 0 {
			s.offset = s.recordEnd - 1
		}
	} else {
		if s.offset >= len(s.memory) {
			s.offset = 0
		} else if s.offset < 0 {
			s.offset = len(s.memory) - 1
		}
	}
}
