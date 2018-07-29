package stdout

import (
	"encoding/binary"
	"math"
	"os"
)

func New(frameSize, sampleRate int) *Stdout {
	return &Stdout{
		frameSize:  frameSize,
		sampleRate: sampleRate,
	}
}

type Stdout struct {
	frameSize, sampleRate int
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
		callback(in, out)
		for i := 0; i < s.frameSize; i++ {
			binary.Write(os.Stdout, binary.LittleEndian, toInt16(out[0][i]))
			binary.Write(os.Stdout, binary.LittleEndian, toInt16(out[1][i]))
		}
	}

	return nil
}

func (s *Stdout) Stop() error     { return nil }
func (s *Stdout) SampleRate() int { return s.sampleRate }
func (s *Stdout) FrameSize() int  { return s.frameSize }

func toInt16(v float32) int16 {
	return int16(v * float32(math.MaxInt16))
}
