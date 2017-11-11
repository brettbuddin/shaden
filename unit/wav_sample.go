package unit

import (
	"buddin.us/lumen/dsp"
	"buddin.us/lumen/wav"
)

func newWAVSample(name string, c Config) (*Unit, error) {
	var config struct {
		Files []string
	}
	if err := c.Decode(&config); err != nil {
		return nil, err
	}

	frames := make([][]float64, len(config.Files))
	for i, filename := range config.Files {
		w, err := wav.Open(filename)
		if err != nil {
			return nil, err
		}
		defer w.Close()
		frames[i], err = loadSamples(w)
		if err != nil {
			return nil, err
		}
	}

	io := NewIO()
	return NewUnit(io, name, &wavSample{
		selection: io.NewIn("select", dsp.Float64(0)),
		reset:     io.NewIn("reset", dsp.Float64(0)),
		offset:    io.NewIn("offset", dsp.Float64(0)),
		direction: io.NewIn("direction", dsp.Float64(1)),
		out:       io.NewOut("out"),
		frames:    frames,
		lastReset: -1,
	}), nil
}

type wavSample struct {
	selection, reset, offset, direction *In
	frames                              [][]float64
	current                             int
	out                                 *Out
	lastReset                           float64
}

func (w *wavSample) ProcessSample(i int) {
	direction := w.direction.Read(i)
	sel := dsp.Clamp(w.selection.Read(i), 0, float64(len(w.frames)-1))
	offset := dsp.Clamp(w.offset.Read(i), 0, 0.95)

	reset := w.reset.Read(i)
	frame := w.frames[int(sel)]

	if w.lastReset < 0 && reset > 0 {
		if direction > 0 {
			w.current = int(float64(len(frame)) * offset)
		} else {
			w.current = int((float64(len(frame)) - 1) - float64(len(frame))*offset)
		}
	}
	if direction > 0 && w.current > len(frame)-1 {
		w.current = int(float64(len(frame)) * offset)
	}
	if direction <= 0 && w.current < 0 {
		w.current = int((float64(len(frame)) - 1) - float64(len(frame))*offset)
	}
	w.out.Write(i, frame[w.current])
	if direction > 0 {
		w.current++
	} else if direction < 0 {
		w.current--
	}
	w.lastReset = reset
}

func loadSamples(w *wav.Wav) ([]float64, error) {
	samples, err := w.ReadAll()
	if err != nil {
		return nil, err
	}
	ratio := int(dsp.SampleRate / float64(w.SampleRate))
	size := len(samples) * ratio
	frame := make([]float64, size)

	var i int
	for _, s := range samples {
		for j := 0; j < ratio; j++ {
			frame[i+j] = float64(s)
		}
		i += ratio
	}
	return frame, nil
}
