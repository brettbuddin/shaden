package unit

import (
	"fmt"

	"github.com/brettbuddin/shaden/graph"
)

// Output is an abstract Out provider. Allows us to wrap an Out in some other behavior if we choose to. Out itself
// implements this interface.
type Output interface {
	Out() *Out
}

// OutputProcessor is an output that can be ticked by the engine.
type OutputProcessor interface {
	Output
	SampleProcessor
	FrameProcessor
}

// Out is a unit output
type Out struct {
	name  string
	unit  *Unit
	node  *graph.Node
	frame []float64
}

// NewOut returns a new output
func NewOut(name string, f []float64) *Out {
	return &Out{
		name:  name,
		frame: f,
	}
}

// Out implements Output interface
func (out *Out) Out() *Out {
	return out
}

// Rate returns the rate of the parent Unit
func (out *Out) Rate() Rate {
	return out.unit.rate
}

// DestinationCount returns the number of outbound connections to this output
func (out *Out) DestinationCount() int {
	return out.node.OutNeighborCount()
}

// Write writes a sample to the output frame if there are downstream consumers of the output
func (out *Out) Write(i int, v float64) {
	out.frame[i] = v
}

// Read reads a sample from the internal buffer
func (out *Out) Read(i int) float64 {
	return out.frame[i]
}

// ExternalNeighborCount returns the count of neighboring nodes outside of the parent Unit
func (out *Out) ExternalNeighborCount() int {
	return out.node.OutNeighborCount()
}

// Destinations returns a list of all destination inputs of the output
func (out *Out) Destinations() []*In {
	var (
		nodes = out.node.OutNeighbors()
		dests = make([]*In, len(nodes))
	)
	for i, n := range nodes {
		dests[i] = n.Value.(*In)
	}
	return dests
}

func (out *Out) String() string {
	return fmt.Sprintf("%s/%s", out.unit.ID, out.name)
}

// OutRef is an unresolved reference to a Unit's Out.
type OutRef struct {
	Unit   *Unit
	Output string
}

func (r OutRef) String() string {
	return fmt.Sprintf("%s/%s", r.Unit.ID, r.Output)
}
