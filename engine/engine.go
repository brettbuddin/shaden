// Package engine provides a low-level audio interface.
package engine

import (
	"fmt"
	"io"
	"time"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/graph"
	"buddin.us/shaden/unit"
)

// Engine is the connection of the synthesizer to PortAudio
type Engine struct {
	messages             MessageChannel
	backend              Backend
	graph                *graph.Graph
	unit                 *unit.Unit
	processors           []unit.FrameProcessor
	errors, stop         chan error
	input                []float64
	lout, rout           []float64
	chunks               int
	singleSampleDisabled bool
	fadeIn               bool
}

// Backend is a low-level callback-based engine
type Backend interface {
	Start(func([]float32, [][]float32)) error
	Stop() error
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
		e.singleSampleDisabled = true
	}
}

// WithFadeIn fades the engine output in to prevent pops
func WithFadeIn() Option {
	return func(e *Engine) {
		e.fadeIn = true
	}
}

// New returns a new Sink
func New(backend Backend, opts ...Option) (*Engine, error) {
	e := &Engine{
		backend:  backend,
		messages: newMessageChannel(),
		graph:    graph.New(),
		errors:   make(chan error),
		stop:     make(chan error),
		input:    make([]float64, dsp.FrameSize),
		chunks:   int(dsp.Float64(backend.FrameSize()) / dsp.FrameSize),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, e.createSink()
}

// UnitBuilders returns all unit.Builders for Units provided by the Engine.
func (e *Engine) UnitBuilders() map[string]unit.Builder {
	return unit.PrepareBuilders(map[string]unit.IOBuilder{
		"source": newSource(e),
	})
}

func (e *Engine) closeProcessors() error {
	for _, p := range e.processors {
		if closer, ok := p.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Reset clears the state of the Engine. This includes clearing the audio graph.
func (e *Engine) Reset() error {
	if err := e.closeProcessors(); err != nil {
		return err
	}
	e.graph = graph.New()

	if err := e.createSink(); err != nil {
		return err
	}
	e.sort()

	return nil
}

func (e *Engine) createSink() error {
	sinkUnit, sink := newSink(e.fadeIn)
	if err := sinkUnit.Attach(e.graph); err != nil {
		return err
	}
	e.unit = sinkUnit
	e.lout = sink.left.out
	e.rout = sink.right.out
	return nil
}

// SendMessage sends a message to to the engine for it to handle within its goroutine
func (e *Engine) SendMessage(msg *Message) error {
	return e.messages.Send(msg)
}

// Errors returns a channel that expresses any errors during operation of the Engine
func (e *Engine) Errors() <-chan error {
	return e.errors
}

// Run starts the Engine; running the audio stream
func (e *Engine) Run() {
	if err := e.backend.Start(e.callback); err != nil {
		e.errors <- err
	}
	<-e.stop

	err := e.closeProcessors()
	if err != nil {
		e.stop <- err
		return
	}
	e.stop <- e.backend.Stop()
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
	case func(g *graph.Graph) (interface{}, error):
		return fn(e.graph)
	default:
		return nil, fmt.Errorf("unhandled function type %T", action)
	}
}

func (e *Engine) sort() {
	processors := e.processors[:0]
	for _, v := range e.graph.Sorted() {
		collectProcessor(&processors, v, e.singleSampleDisabled)
	}
	e.processors = processors
	e.graph.AckChange()
}

func (e *Engine) handle(msg *Message) {
	start := time.Now()
	data, err := e.call(msg.Action)

	if err == nil && e.graph.HasChanged() {
		e.sort()
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

		offset := dsp.FrameSize * k
		for i := 0; i < int(dsp.FrameSize); i++ {
			e.input[i] = float64(in[offset+i])
		}
		for _, p := range e.processors {
			p.ProcessFrame(dsp.FrameSize)
		}
		for i := range out {
			for j := 0; j < dsp.FrameSize; j++ {
				if i%2 == 0 {
					out[i][offset+j] = float32(e.lout[j])
				} else {
					out[i][offset+j] = float32(e.rout[j])
				}
			}
		}
	}
}

func collectProcessor(processors *[]unit.FrameProcessor, nodes []*graph.Node, singleSampleDisabled bool) {
	if len(nodes) > 1 {
		collectGroup(processors, nodes, singleSampleDisabled)
		return
	}

	first := nodes[0]
	if in, ok := first.Value.(*unit.In); ok && !singleSampleDisabled {
		in.Mode = unit.Block
	}
	if p, ok := first.Value.(unit.FrameProcessor); ok {
		if isp, ok := p.(unit.CondProcessor); ok {
			if isp.IsProcessable() {
				*processors = append(*processors, p)
			}
		} else {
			*processors = append(*processors, p)
		}
	}
}

func collectGroup(processors *[]unit.FrameProcessor, nodes []*graph.Node, singleSampleDisabled bool) {
	var g group
	for _, w := range nodes {
		if in, ok := w.Value.(*unit.In); ok && !singleSampleDisabled {
			in.Mode = unit.Sample
		}
		if p, ok := w.Value.(unit.SampleProcessor); ok {
			if isp, ok := p.(unit.CondProcessor); ok {
				if isp.IsProcessable() {
					g.processors = append(g.processors, p)
				}
			} else {
				g.processors = append(g.processors, p)
			}
		}
	}
	*processors = append(*processors, g)
}

type group struct {
	processors []unit.SampleProcessor
}

func (g group) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		for _, p := range g.processors {
			p.ProcessSample(i)
		}
	}
}

func (g group) Close() error {
	for _, p := range g.processors {
		if closer, ok := p.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
