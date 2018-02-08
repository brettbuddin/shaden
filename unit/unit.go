// Package unit provides built-in units for synthesis.
package unit

import (
	"io"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/errors"
	"buddin.us/shaden/graph"
)

// FrameProcessor processes a block of samples of a given size.
type FrameProcessor interface {
	ProcessFrame(n int)
}

// SampleProcessor processes a single sample of signal.
type SampleProcessor interface {
	ProcessSample(i int)
}

// CondProcessor informs another party whether or not it should be processed.
type CondProcessor interface {
	IsProcessable() bool
}

// Rate is a rate in which signals will be processed by units
type Rate int

// Supported processing rates
const (
	RateAudio Rate = iota
	RateControl
)

// Unit is a synthesizer unit
type Unit struct {
	*IO
	SampleProcessor
	rate Rate
	node *graph.Node
}

// NewUnit creates a new Unit that defaults to audio rate.
func NewUnit(io *IO, p SampleProcessor) *Unit {
	return &Unit{IO: io, SampleProcessor: p}
}

// IsProcessable determines whether or not this Unit's ProcessFrame method should be called by the engine
func (u *Unit) IsProcessable() bool {
	return u.node != nil && u.SampleProcessor != nil
}

// ProcessFrame calculates a block of samples.
func (u *Unit) ProcessFrame(n int) {
	if p, ok := u.SampleProcessor.(FrameProcessor); ok {
		p.ProcessFrame(n)
		return
	}
	if u.rate == RateControl {
		u.ProcessSample(0)
		return
	}
	for i := 0; i < n; i++ {
		u.ProcessSample(i)
	}
}

// Close closes the Processor if it is an io.Closer. It also closes any Outs that it has.
func (u *Unit) Close() error {
	if c, ok := u.SampleProcessor.(io.Closer); ok {
		if err := c.Close(); err != nil {
			return errors.Wrap(err, "close processor failed")
		}
	}
	for _, o := range u.Out {
		if c, ok := o.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return errors.Wrap(err, "close output failed")
			}
		}
	}
	return nil
}

// Attach connects this unit and its inputs/outputs to a Graph
func (u *Unit) Attach(g *graph.Graph) error {
	n := g.NewNode(u)
	u.node = n

	for _, e := range u.In {
		e.unit = u
		c := g.NewNode(e)
		e.node = c
		if err := g.NewConnection(c, n); err != nil {
			return errors.Wrap(err, "new input connection failed")
		}
	}
	for _, e := range u.Out {
		e.Out().unit = u
		c := g.NewNode(e)
		e.Out().node = c
		if err := g.NewConnection(n, c); err != nil {
			return errors.Wrap(err, "new output connection failed")
		}
	}
	return nil
}

// Detach removes this unit and its inputs/outputs from a Graph
func (u *Unit) Detach(g *graph.Graph) error {
	if err := g.RemoveNode(u.node); err != nil {
		return errors.Wrap(err, "remove node failed")
	}
	for _, e := range u.In {
		if e.HasSource() {
			sources := e.node.InNeighbors()
			for _, n := range sources {
				if err := g.RemoveConnection(n, e.node); err != nil {
					return errors.Wrap(err, "remove input connection failed")
				}
			}
		}
		if err := g.RemoveNode(e.node); err != nil {
			return errors.Wrap(err, "remove input node failed")
		}
	}
	for _, e := range u.Out {
		if out := e.Out(); out.DestinationCount() > 0 {
			dests := out.node.OutNeighbors()
			for _, n := range dests {
				n.Value.(*In).Reset()
				if err := g.RemoveConnection(out.node, n); err != nil {
					return errors.Wrap(err, "remove output connection failed")
				}
			}
		}
		if err := g.RemoveNode(e.Out().node); err != nil {
			return errors.Wrap(err, "remove output node failed")
		}
	}
	return nil
}

// ExternalNeighborCount returns the count of neighboring nodes outside of the Unit
func (u *Unit) ExternalNeighborCount() int {
	var n int
	for _, in := range u.In {
		n += in.ExternalNeighborCount()
	}
	for _, out := range u.Out {
		n += out.Out().ExternalNeighborCount()
	}
	return n
}

func isTrig(last, current float64) bool {
	return last <= 0 && current > 0
}

func isHigh(v float64) bool {
	return v > 0
}

func newFrame() []float64 {
	return make([]float64, dsp.FrameSize)
}
