package unit

import (
	"math"
	"os"

	"github.com/go-audio/wav"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/errors"
)

func newWAVSample(name string, c Config) (*Unit, error) {
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
	if w.IsValidFile() {
		return nil, errors.Errorf("%q is not a valid WAV file", config.File)
	}

	buf, err := w.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	var (
		raw   = buf.AsFloatBuffer().Data
		frame = make([]float64, len(raw))
		max   float64
	)
	for _, s := range raw {
		if s > max {
			max = s
		}
	}
	for i, s := range raw {
		frame[i] = s / max
	}

	io := NewIO()
	return NewUnit(io, name, &wavSample{
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
		begin     = dsp.Clamp(w.begin.Read(i), 0, 0.95)
		end       = math.Max(begin, dsp.Clamp(w.end.Read(i), 0, 1))
		trigger   = w.trigger.Read(i)
		cycle     = w.cycle.Read(i)
		channels  = w.channels
		length    = int(math.Min(float64(w.length), float64(w.length)*end))
	)

	// Trigger and reset
	if isTrig(w.lastTrigger, trigger) {
		if direction > 0 {
			w.current = int(float64(length) * begin)
		} else {
			w.current = int((float64(length) - 1) - float64(length)*begin)
		}
		w.playing = true
	}

	// Wrapping
	if direction > 0 && w.current > length-1 {
		w.current = int(float64(length) * begin)
		if cycle <= 0 {
			w.playing = false
		}
	}
	if direction <= 0 && w.current < 0 {
		w.current = int((float64(length) - 1) - float64(length)*begin)
		if cycle <= 0 {
			w.playing = false
		}
	}

	// Write output
	if w.playing {
		w.a.Write(i, w.frame[w.current])

		var incr int
		switch channels {
		case 1:
			w.b.Write(i, 0)
			incr = 1
		case 2:
			w.b.Write(i, w.frame[w.current+1])
			incr = 2
		}

		if direction > 0 {
			w.current += incr
		} else if direction < 0 {
			w.current -= incr
		}
	} else {
		w.a.Write(i, 0)
		w.b.Write(i, 0)
	}

	w.lastTrigger = trigger
}
