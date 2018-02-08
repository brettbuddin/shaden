package unit

import (
	"math"
	"os"

	"github.com/go-audio/wav"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/errors"
)

func newWAVSample(io *IO, c Config) (*Unit, error) {
	var config struct {
		File string
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	if config.File == "" {
		return nil, errors.New("no WAV file specified")
	}

	f, err := os.Open(config.File)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	w := wav.NewDecoder(f)
	if !w.IsValidFile() {
		return nil, errors.Errorf("%q is not a valid WAV file", config.File)
	}

	buf, err := w.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	var (
		raw   = buf.AsFloat32Buffer().Data
		frame = make([]float64, len(raw))
	)
	for i, s := range raw {
		frame[i] = float64(s)
	}

	return NewUnit(io, &wavSample{
		trigger:     io.NewIn("trigger", dsp.Float64(-1)),
		direction:   io.NewIn("direction", dsp.Float64(1)),
		begin:       io.NewIn("begin", dsp.Float64(0)),
		end:         io.NewIn("end", dsp.Float64(1)),
		cycle:       io.NewIn("cycle", dsp.Float64(0)),
		a:           io.NewOut("a"),
		b:           io.NewOut("b"),
		channels:    buf.Format.NumChannels,
		length:      buf.NumFrames(),
		frame:       frame,
		lastTrigger: -1,
	}), nil
}

type wavSample struct {
	trigger, begin, end, direction, cycle *In
	length, channels                      int
	frame                                 []float64
	current                               int
	a, b                                  *Out
	playing                               bool
	lastTrigger                           float64
}

func (w *wavSample) ProcessSample(i int) {
	var (
		direction = w.direction.Read(i)
		length    = float64(w.length - 1)
		begin     = w.begin.Read(i)
		end       = w.end.Read(i)
		absBegin  = int(math.Min(end, dsp.Clamp(begin, 0, 0.95)) * length)
		absEnd    = int(math.Max(begin, dsp.Clamp(end, 0.05, 1)) * length)
		trigger   = w.trigger.Read(i)
		cycle     = w.cycle.Read(i)
		channels  = w.channels
	)

	// Trigger and reset
	if isTrig(w.lastTrigger, trigger) {
		if direction > 0 {
			w.current = absBegin
		} else {
			w.current = absEnd
		}
		w.playing = true
	}

	// Wrapping
	if direction > 0 && w.current > absEnd {
		w.current = absBegin
		if cycle <= 0 {
			w.playing = false
		}
	}
	if direction <= 0 && w.current < absBegin {
		w.current = absEnd
		if cycle <= 0 {
			w.playing = false
		}
	}

	// Write output
	if w.playing {
		var left, right = w.a, w.b
		if direction <= 0 {
			right, left = left, right
		}

		left.Write(i, w.frame[w.current])

		var incr int
		switch channels {
		case 1:
			right.Write(i, 0)
			incr = 1
		case 2:
			right.Write(i, w.frame[w.current+1])
			incr = 2
		}

		if direction > 0 {
			w.current += incr
		} else if direction <= 0 {
			w.current -= incr
		}
	} else {
		w.a.Write(i, 0)
		w.b.Write(i, 0)
	}

	w.lastTrigger = trigger
}
