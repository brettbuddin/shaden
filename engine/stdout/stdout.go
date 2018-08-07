package stdout

import (
	"encoding/binary"
	"io"
	"math"
	"runtime"
	"sync"
)

// New returns a new Stdout
func New(out io.Writer, frameSize, sampleRate int) *Stdout {
	return &Stdout{
		out:        out,
		frameSize:  frameSize,
		sampleRate: sampleRate,
		running:    true,
	}
}

// Stdout is an engine backend that writes little-endian int16s to an output
// stream (stdout).
type Stdout struct {
	out                   io.Writer
	frameSize, sampleRate int

	mutex   sync.Mutex
	running bool
}

// Start starts the backend.
func (s *Stdout) Start(callback func([]float32, [][]float32)) error {
	var (
		in  = make([]float32, s.frameSize)
		out = [][]float32{
			make([]float32, s.frameSize),
			make([]float32, s.frameSize),
		}
	)

	go func() {
		for {
			s.mutex.Lock()
			if !s.running {
				return
			}
			s.mutex.Unlock()
			callback(in, out)
			for i := 0; i < s.frameSize; i++ {
				binary.Write(s.out, binary.LittleEndian, toInt16(out[0][i]))
				binary.Write(s.out, binary.LittleEndian, toInt16(out[1][i]))
			}
			runtime.Gosched()
		}
	}()

	return nil
}

// Stop stops the backend.
func (s *Stdout) Stop() error {
	s.mutex.Lock()
	s.running = false
	s.mutex.Unlock()
	return nil
}

// SampleRate returns the sample rate of the backend.
func (s *Stdout) SampleRate() int { return s.sampleRate }

// FrameSize returns the frame size of the backend.
func (s *Stdout) FrameSize() int { return s.frameSize }

func toInt16(v float32) int16 {
	return int16(v * float32(math.MaxInt16))
}
