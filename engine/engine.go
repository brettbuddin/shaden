// Package engine provides a low-level audio interface.
package engine

import (
	"fmt"
	"time"

	"github.com/brettbuddin/shaden/unit"
)

// Backend is a low-level callback-based engine
type Backend interface {
	Start(func([]float32, [][]float32)) error
	Stop() error
	SampleRate() int
	FrameSize() int
}

// Option is an option for the Engine
type Option func(*Engine)

// WithMessageChannel establishes a MessageChannel to be used for sending and receiving messages within the Engine.
func WithMessageChannel(ch MessageChannel) Option {
	return func(e *Engine) {
		e.messages = ch
	}
}

// WithSingleSampleDisabled disables the single-sample feedback loop behavior.
func WithSingleSampleDisabled() Option {
	return func(e *Engine) {
		e.graph.singleSampleDisabled = true
	}
}

// WithFadeIn fades the engine output in to prevent pops
func WithFadeIn(ms int) Option {
	return func(e *Engine) {
		e.fadeIn = ms
	}
}

// WithGain sets the global gain for all samples written to the output
func WithGain(gain float32) Option {
	return func(e *Engine) {
		e.gain = gain
	}
}

// Engine is the connection of the synthesizer to PortAudio
type Engine struct {
	messages     MessageChannel
	backend      Backend
	graph        *Graph
	errors, stop chan error
	chunks       int
	fadeIn       int
	frameSize    int
	gain         float32
}

// New returns a new Sink
func New(backend Backend, frameSize int, opts ...Option) (*Engine, error) {
	e := &Engine{
		backend:   backend,
		messages:  newMessageChannel(),
		graph:     NewGraph(frameSize),
		errors:    make(chan error),
		stop:      make(chan error),
		chunks:    int(backend.FrameSize() / frameSize),
		frameSize: frameSize,
		gain:      1,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, e.graph.Reset(e.fadeIn, e.frameSize, backend.SampleRate())
}

// SampleRate returns the sample rate
func (e *Engine) SampleRate() int { return e.backend.SampleRate() }

// FrameSize returns the frame size
func (e *Engine) FrameSize() int { return e.frameSize }

// UnitBuilders returns all unit.Builders for Units provided by the Engine.
func (e *Engine) UnitBuilders() map[string]unit.Builder {
	return unit.PrepareBuilders(map[string]unit.IOBuilder{
		"source": e.graph.sourceIOBuilder(),
	})
}

// Reset clears the state of the Engine. This includes clearing the audio graph.
func (e *Engine) Reset() error {
	return e.graph.Reset(e.fadeIn, e.frameSize, e.backend.SampleRate())
}

// SendMessage sends a message to to the engine for it to handle within its goroutine
func (e *Engine) SendMessage(msg *Message) error {
	return e.messages.Send(msg)
}

// Errors returns a channel that expresses any errors during operation of the Engine
func (e *Engine) Errors() <-chan error { return e.errors }

// Run starts the Engine; running the audio stream
func (e *Engine) Run() {
	if err := e.backend.Start(e.callback); err != nil {
		e.errors <- err
	}
	<-e.stop

	if err := e.backend.Stop(); err != nil {
		e.stop <- err
		return
	}
	e.stop <- e.graph.Close()
}

// Stop shuts down the Engine
func (e *Engine) Stop() error {
	e.stop <- nil
	err := <-e.stop
	close(e.errors)
	close(e.stop)
	return err
}

func (e *Engine) call(action interface{}) (interface{}, error) {
	switch fn := action.(type) {
	case func(e *Engine) (interface{}, error):
		return fn(e)
	case func(g *Graph) (interface{}, error):
		return fn(e.graph)
	default:
		return nil, fmt.Errorf("unhandled function type %T", action)
	}
}

func (e *Engine) handle(msg *Message) {
	start := time.Now()
	data, err := e.call(msg.Action)

	if err == nil {
		e.graph.Sort()
	}

	if msg.Reply != nil {
		msg.Reply <- &Reply{
			Duration: time.Since(start),
			Data:     data,
			Error:    err,
		}
	}
}

// callback is the callback function provided to PortAudio; it drives the entire synthesiser.
func (e *Engine) callback(in []float32, out [][]float32) {
	for k := 0; k < e.chunks; k++ {
		if msg := e.messages.Receive(); msg != nil {
			e.handle(msg)
		}

		var (
			frameSize = e.frameSize
			offset    = frameSize * k
			input     = e.graph.in
			leftOut   = e.graph.leftOut
			rightOut  = e.graph.rightOut
			gain      = e.gain
		)
		for i := 0; i < frameSize; i++ {
			input[i] = float64(in[offset+i])
		}
		for _, p := range e.graph.Processors() {
			p.ProcessFrame(frameSize)
		}
		for i := range out {
			for j := 0; j < frameSize; j++ {
				if i%2 == 0 {
					out[i][offset+j] = float32(leftOut[j]) * gain
				} else {
					out[i][offset+j] = float32(rightOut[j]) * gain
				}
			}
		}
	}
}
