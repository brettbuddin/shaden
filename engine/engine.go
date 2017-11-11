package engine

import (
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"

	"buddin.us/lumen/dsp"
	"buddin.us/lumen/graph"
	"buddin.us/lumen/unit"
)

// Engine is the connection of the synthesizer to PortAudio
type Engine struct {
	backend              *portAudio
	module               *unit.Unit
	left, right          *unit.In
	graph                *graph.Graph
	processors           []unit.FrameProcessor
	messages             chan *Message
	errors, stop         chan error
	input                []float64
	chunks               int
	singleSampleDisabled bool
}

// New returns a new Sink
func New(inDeviceIndex, outDeviceIndex int, latency string, frameSize int, singleSampleDisabled bool) (*Engine, error) {
	if frameSize%dsp.FrameSize != 0 {
		return nil, fmt.Errorf("frame size (%d) must be a multiple of %d", frameSize, dsp.FrameSize)
	}

	g := graph.New()

	sink := newSink()
	if err := sink.Attach(g); err != nil {
		return nil, err
	}

	backend, err := newPortAudio(inDeviceIndex, outDeviceIndex, latency, frameSize)
	if err != nil {
		return nil, err
	}

	return &Engine{
		backend:              backend,
		graph:                g,
		messages:             make(chan *Message),
		module:               sink,
		left:                 sink.In["l"],
		right:                sink.In["r"],
		errors:               make(chan error),
		stop:                 make(chan error),
		input:                make([]float64, dsp.FrameSize),
		chunks:               int(dsp.Float64(frameSize) / dsp.FrameSize),
		singleSampleDisabled: singleSampleDisabled,
	}, nil
}

// UnitBuilders returns all unit.BuildFuncs for Units provided by the Engine.
func (e *Engine) UnitBuilders() map[string]unit.BuildFunc {
	return unitBuilders(e)
}

// Reset clears the state of the Engine. This includes clearing the audio graph.
func (e *Engine) Reset() error {
	e.graph = graph.New()

	sink := newSink()
	if err := sink.Attach(e.graph); err != nil {
		return err
	}
	e.module = sink
	e.left = sink.In["l"]
	e.right = sink.In["r"]

	e.sort()

	return nil
}

// Devices return input and output devices being used
func (e *Engine) Devices() (in *portaudio.DeviceInfo, out *portaudio.DeviceInfo) {
	return e.backend.inDevice, e.backend.outDevice
}

// Messages provides a send-only channel that can be used to execute code on the main audio goroutine
func (e *Engine) Messages() chan<- *Message {
	return e.messages
}

// Errors returns a channel that expresses any errors during operation of the Engine
func (e *Engine) Errors() <-chan error {
	return e.errors
}

// Run starts the Engine; running the audio stream
func (e *Engine) Run() {
	if err := e.backend.Start(e.callback); err != nil {
		e.errors <- err
		return
	}
	<-e.stop
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
	case func(g *graph.Graph, l, r *unit.In) (interface{}, error):
		return fn(e.graph, e.left, e.right)
	default:
		return nil, fmt.Errorf("unhandled function type %T", action)
	}
}

func (e *Engine) sort() {
	processors := e.processors[:0]
	for _, v := range e.graph.Sorted() {
		e.collectProcessor(&processors, v)
	}
	e.processors = processors
	e.graph.AckChange()
}

func (e *Engine) collectProcessor(processors *[]unit.FrameProcessor, nodes []*graph.Node) {
	if len(nodes) > 1 {
		e.collectGroup(processors, nodes)
		return
	}

	first := nodes[0]
	if in, ok := first.Value.(*unit.In); ok && !e.singleSampleDisabled {
		in.Mode = unit.Block
	}
	if p, ok := first.Value.(frameProcessor); ok && p.ExternalNeighborCount() > 0 {
		if isp, ok := p.(condProcessor); ok {
			if isp.IsProcessable() {
				*processors = append(*processors, p)
			}
		} else {
			*processors = append(*processors, p)
		}
	}
}

func (e *Engine) collectGroup(processors *[]unit.FrameProcessor, nodes []*graph.Node) {
	var g group
	for _, w := range nodes {
		if in, ok := w.Value.(*unit.In); ok && !e.singleSampleDisabled {
			in.Mode = unit.Sample
		}
		if p, ok := w.Value.(unit.SampleProcessor); ok {
			if isp, ok := p.(condProcessor); ok {
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
		select {
		case msg := <-e.messages:
			e.handle(msg)
		default:
		}

		offset := int(dsp.FrameSize * k)
		for i := 0; i < int(dsp.FrameSize); i++ {
			e.input[i] = float64(in[offset+i])
		}
		for _, p := range e.processors {
			p.ProcessFrame(dsp.FrameSize)
		}
		for i := range out {
			for j := 0; j < dsp.FrameSize; j++ {
				if i%2 == 0 {
					out[i][offset+j] = float32(e.left.Read(j))
				} else {
					out[i][offset+j] = float32(e.right.Read(j))
				}
			}
		}
	}
}

type frameProcessor interface {
	unit.FrameProcessor
	ExternalNeighborCount() int
}

type condProcessor interface {
	IsProcessable() bool
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
