package graph

import (
	"fmt"
)

// Node is a member of a Graph
type Node struct {
	graph           *Graph
	idx             int
	searchState     searchState
	outputs, inputs []connection
	Value           interface{}
}

// Neighbors returns the Nodes connected to this Node in the Graph
func (n *Node) Neighbors() []*Node {
	nodes := make([]*Node, len(n.inputs)+len(n.outputs))
	for i, e := range n.inputs {
		nodes[i] = e.end
	}
	for i, e := range n.outputs {
		nodes[(len(n.inputs)-1)+i] = e.end
	}
	return nodes
}

// InNeighbors returns only neighboring Nodes with an inbound connection to this Node
func (n *Node) InNeighbors() []*Node {
	nodes := make([]*Node, 0, len(n.inputs))
	for _, input := range n.inputs {
		nodes = append(nodes, input.end)
	}
	return nodes
}

// OutNeighbors returns only neighboring Nodes with an outbound connection from this Node
func (n *Node) OutNeighbors() []*Node {
	nodes := make([]*Node, 0, len(n.outputs))
	for _, output := range n.outputs {
		nodes = append(nodes, output.end)
	}
	return nodes
}

// InNeighborCount returns the count of inbound neighbors
func (n *Node) InNeighborCount() int {
	if n == nil {
		return 0
	}
	return len(n.inputs)
}

// NeighborCount returns the total neighbor count
func (n *Node) NeighborCount() int {
	if n == nil {
		return 0
	}
	return len(n.outputs) + len(n.inputs)
}

// OutNeighborCount returns the count of outbound neighbors
func (n *Node) OutNeighborCount() int {
	if n == nil {
		return 0
	}
	return len(n.outputs)
}

func (n *Node) removeOutConnection(to *Node) (bool, error) {
	outputs := n.outputs
	for i := range outputs {
		if outputs[i].end == to {
			if err := removeConnectionIndex(i, &outputs); err != nil {
				return false, err
			}
			n.outputs = outputs
			return true, nil
		}
	}
	return false, nil
}

func (n *Node) removeInConnection(from *Node) (bool, error) {
	inputs := n.inputs
	for i := range inputs {
		if inputs[i].end == from {
			if err := removeConnectionIndex(i, &inputs); err != nil {
				return false, err
			}
			n.inputs = inputs
			return true, nil
		}
	}
	return false, nil
}

func reverseNodes(list []*Node) {
	length := len(list)
	for i := 0; i < length/2; i++ {
		list[i], list[length-i-1] = list[length-i-1], list[i]
	}
}

type searchState int

const (
	unseen searchState = iota
	seen
)

type connection struct {
	start, end *Node
}

func removeConnectionIndex(pos int, e *[]connection) error {
	if pos < 0 || pos > len(*e)-1 {
		return fmt.Errorf("position (%d) out of the range of connection list size (%d)", pos, len(*e))
	}
	(*e)[pos], (*e)[len(*e)-1] = (*e)[len(*e)-1], (*e)[pos]
	*e = (*e)[:len(*e)-1]
	return nil
}
