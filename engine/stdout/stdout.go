package stdout

import (
	"encoding/binary"
	"io"
	"math"
	"sync"

	"github.com/brettbuddin/shaden/engine"
)

var _ engine.Backend = &Stdout{}

func New(out io.Writer, frameSize, sampleRate int) *Stdout {
	return &Stdout{
		out:        out,
		frameSize:  frameSize,
		sampleRate: sampleRate,
		running:    true,
	}
}

type Stdout struct {
	out                   io.Writer
	frameSize, sampleRate int

	mutex   sync.Mutex
	running bool
}

func (s *Stdout) Start(callback func([]float32, [][]float32)) error {
	var (
		in  = make([]float32, s.frameSize)
		out = [][]float32{
			make([]float32, s.frameSize),
			make([]float32, s.frameSize),
		}
	)

	for {
		s.mutex.Lock()
		if !s.running {
			return nil
		}
		s.mutex.Unlock()
		callback(in, out)
		for i := 0; i < s.frameSize; i++ {
			binary.Write(s.out, binary.LittleEndian, toInt16(out[0][i]))
			binary.Write(s.out, binary.LittleEndian, toInt16(out[1][i]))
		}
	}

	return nil
}

func (s *Stdout) Stop() error {
	s.mutex.Lock()
	s.running = false
	s.mutex.Unlock()
	return nil
}
func (s *Stdout) SampleRate() int { return s.sampleRate }
func (s *Stdout) FrameSize() int  { return s.frameSize }

func toInt16(v float32) int16 {
	return int16(v * float32(math.MaxInt16))
}
