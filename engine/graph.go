package engine

import (
	"fmt"
	"io"

	"github.com/brettbuddin/shaden/dsp"
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/graph"
	"github.com/brettbuddin/shaden/unit"
)

func NewGraph(frameSize int) *Graph {
	return &Graph{
		graph:      graph.New(),
		processors: make([]unit.FrameProcessor, 100),
		in:         make([]float64, frameSize),
		leftOut:    make([]float64, frameSize),
		rightOut:   make([]float64, frameSize),
	}
}

type Graph struct {
	singleSampleDisabled  bool
	graph                 *graph.Graph
	processors            []unit.FrameProcessor
	sink                  *unit.Unit
	in, leftOut, rightOut []float64
}

func (g *Graph) Sort() {
	if !g.graph.HasChanged() {
		return
	}
	processors := g.processors[:0]
	for _, v := range g.graph.Sorted() {
		collectProcessor(&processors, v, g.singleSampleDisabled)
	}
	g.processors = processors
	g.graph.AckChange()
}

func (g *Graph) Reset(fadeIn, frameSize, sampleRate int) error {
	if err := g.Close(); err != nil {
		return err
	}
	g.graph = graph.New()

	if err := g.createSink(fadeIn, frameSize, sampleRate); err != nil {
		return err
	}
	g.Sort()

	return nil
}

func (g *Graph) createSink(fadeIn, frameSize, sampleRate int) error {
	var (
		io       = unit.NewIO("sink", frameSize)
		sink     = newSink(io, fadeIn, sampleRate, frameSize)
		sinkUnit = unit.NewUnit(io, sink)
	)
	if err := sinkUnit.Attach(g.graph); err != nil {
		return err
	}
	g.sink = sinkUnit
	g.leftOut = sink.left.out
	g.rightOut = sink.right.out
	return nil
}

func (g *Graph) Close() error {
	for _, p := range g.processors {
		if closer, ok := p.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Graph) Patch(v interface{}, in *unit.In) error {
	switch v := v.(type) {
	case float64:
		if err := unit.Unpatch(g.graph, in); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unpatch %q", in))
		}
		in.Fill(dsp.Float64(v))
	case int:
		if err := unit.Unpatch(g.graph, in); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unpatch %q", in))
		}
		in.Fill(dsp.Float64(v))
	case dsp.Valuer:
		if err := unit.Unpatch(g.graph, in); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unpatch %q", in))
		}
		in.Fill(v)
	case unit.OutRef:
		out, ok := v.Unit.Out[v.Output]
		if !ok {
			return errors.Errorf("unit %q has no output %q", v.Unit.ID, v.Output)
		}
		if err := unit.Patch(g.graph, out, in); err != nil {
			return errors.Wrap(err, fmt.Sprintf("patch %q into %q", out.Out(), in))
		}
	}
	return nil
}

func (g *Graph) Mount(u *unit.Unit) error { return u.Attach(g.graph) }

func (g *Graph) Unmount(u *unit.Unit) error {
	if err := u.Close(); err != nil {
		return err
	}
	if err := u.Detach(g.graph); err != nil {
		switch err := err.(type) {
		case graph.NotInGraphError:
			return errors.Errorf("unit %q not in graph", u.ID)
		default:
			return err
		}
	}
	return nil
}

func (g *Graph) HasChanged() bool { return g.graph.HasChanged() }
func (g *Graph) Size() int        { return g.graph.Size() }

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
