package unit

import (
	"fmt"

	"buddin.us/shaden/graph"
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
}

// Out is a module output
type Out struct {
	Name  string
	unit  *Unit
	node  *graph.Node
	frame []float64
	last  float64
}

// NewOut returns a new output
func NewOut(name string, f []float64) *Out {
	return &Out{
		Name:  name,
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
	if out.node.OutNeighborCount() == 0 {
		return
	}
	out.frame[i] = v
}

// ExternalNeighborCount returns the count of neighboring nodes outside of the parent Unit
func (out *Out) ExternalNeighborCount() int {
	return out.node.OutNeighborCount()
}

func (out *Out) String() string {
	return fmt.Sprintf("%s/%s", out.unit.ID, out.Name)
}

// OutRef is an unresolved reference to a Unit's Out.
type OutRef struct {
	Unit   *Unit
	Output string
}

func (r OutRef) String() string {
	return fmt.Sprintf("%s/%s", r.Unit.ID, r.Output)
}
