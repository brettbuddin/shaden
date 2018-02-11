package unit

import (
	"fmt"
	"math"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/graph"
)

// InMode is a mode of processing of an In.
type InMode int

// InModes
const (
	Block InMode = iota
	Sample
)

const controlPeriod = 64

// In is a unit input
type In struct {
	Name               string
	Mode               InMode
	normal             dsp.Valuer
	frame, normalFrame []float64
	unit               *Unit
	source             *Out
	node               *graph.Node

	controlLastF float64
	controlLastI int
}

// NewIn returns a new input
func NewIn(name string, v dsp.Valuer) *In {
	f := newFrame()
	in := &In{
		Name:        name,
		frame:       f,
		normalFrame: f,
	}
	in.setNormal(v)
	return in
}

// Read reads a specific sample from the input frame
func (in *In) Read(i int) float64 {
	if isSourceControlRate(in) {
		return in.frame[0]
	}
	if in.Mode == Sample {
		size := len(in.frame)
		i = (i - 1 + size) % size
	}
	return in.frame[i]
}

// ReadSlow reads a specific sample from the input frame at a slow rate
func (in *In) ReadSlow(i int, f func(float64) float64) float64 {
	if i%controlPeriod == 0 {
		in.controlLastF = f(in.Read(i))
	}
	return in.controlLastF
}

// ReadSlowInt reads a specific sample from the input frame at a slow rate
func (in *In) ReadSlowInt(i int, f func(int) int) int {
	if i%controlPeriod == 0 {
		in.controlLastI = f(int(in.Read(i)))
	}
	return in.controlLastI
}

// Fill fills the internal frame with a specific constant value
func (in *In) Fill(v dsp.Valuer) {
	for i := range in.frame {
		in.frame[i] = v.Float64()
	}
}

// Write writes a sample to the internal buffer
func (in *In) Write(i int, v float64) {
	in.frame[i] = v
}

// Couple assigns the internal frame of this input to the frame of an output; binding them together. This in-of-itself
// does not define the connection. That is controlled by the the Nodes and Graph.
func (in *In) Couple(out Output) {
	o := out.Out()
	in.source = o
	in.frame = o.frame
}

// HasSource returns whether or not we have an inbound connection
func (in *In) HasSource() bool {
	return in.source != nil
}

// Reset disconnects an input from an output (if a connection has been established) and fills the frame with the normal
// constant value
func (in *In) Reset() {
	in.source = nil
	in.frame = in.normalFrame
	in.Fill(in.normal)
}

// ExternalNeighborCount returns the count of neighboring nodes outside of the parent Unit
func (in *In) ExternalNeighborCount() int {
	return in.node.InNeighborCount()
}

func (in *In) setNormal(v dsp.Valuer) {
	in.normal = v
	in.Fill(v)
}

func (in *In) String() string {
	return fmt.Sprintf("%s/%s", in.unit.ID, in.Name)
}

func isSourceControlRate(in *In) bool {
	return in.HasSource() && in.source.Rate() == RateControl
}

func ident(v float64) float64   { return v }
func minZero(v float64) float64 { return math.Max(v, 0) }
