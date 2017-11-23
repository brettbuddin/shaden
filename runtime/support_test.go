package runtime

import (
	"errors"
	"sync"
	"time"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/engine"
)

func newBackend(calls int) *backend {
	written := make([][]float32, 2)
	for i := range written {
		written[i] = make([]float32, dsp.FrameSize)
	}
	return &backend{
		calls:   calls,
		written: written,
	}
}

type backend struct {
	sync.Mutex
	calls   int
	written [][]float32
}

func (b *backend) read(i, j int) float32 {
	b.Lock()
	defer b.Unlock()
	return b.written[i][j]
}

func (b *backend) Start(cb func([]float32, [][]float32)) error {
	b.Lock()
	defer b.Unlock()
	var (
		in  = make([]float32, dsp.FrameSize)
		out = [][]float32{
			make([]float32, dsp.FrameSize),
			make([]float32, dsp.FrameSize),
		}
	)
	for i := 0; i < b.calls; i++ {
		cb(in, out)
	}
	copy(b.written, out)
	return nil
}
func (backend) Stop() error    { return nil }
func (backend) FrameSize() int { return dsp.FrameSize }

type messageChannel struct {
	messages chan *engine.Message
}

func (c messageChannel) Send(msg *engine.Message) error {
	select {
	case c.messages <- msg:
	case <-time.After(10 * time.Second):
		return errors.New("timeout sending message")
	}
	return nil
}

func (c messageChannel) Receive() *engine.Message {
	return <-c.messages
}

func (c messageChannel) Close() { close(c.messages) }
