// Package graph provides a directed graph implementation.
package graph

import (
	"fmt"
)

// Graph is a directed graph.
type Graph struct {
	nodes      []*Node
	sorted     []*Node
	components [][]*Node
	dirty      bool
}

// New returns a new Graph.
func New() *Graph {
	return &Graph{
		sorted:     make([]*Node, 0, 1024),
		components: make([][]*Node, 0, 1024),
	}
}

// Size returns the number of nodes in the graph.
func (g *Graph) Size() int { return len(g.nodes) }

// NewNode creates a new Node in the Graph.
func (g *Graph) NewNode(v interface{}) *Node {
	n := &Node{
		graph: g,
		idx:   len(g.nodes),
		Value: v,
	}
	g.nodes = append(g.nodes, n)
	return n
}

// RemoveNode removes a Node from the Graph.
func (g *Graph) RemoveNode(n *Node) error {
	if !g.Exists(n) {
		return NotInGraphError{Node: n}
	}

	for _, node := range g.nodes {
		if node == n {
			continue
		}
		removed, err := node.removeOutConnection(n)
		if err != nil {
			return err
		}
		g.maybeDirty(removed)

		removed, err = node.removeInConnection(n)
		if err != nil {
			return err
		}
		g.maybeDirty(removed)

		if node.idx > n.idx {
			g.dirty = true
			node.idx--
		}
	}
	copy(g.nodes[n.idx:], g.nodes[n.idx+1:])
	g.nodes = g.nodes[:len(g.nodes)-1]

	return nil
}

// NewConnection creates a connection between two Nodes.
func (g *Graph) NewConnection(from, to *Node) error {
	if !g.Exists(from) {
		return NotInGraphError{Node: from}
	}
	if !g.Exists(to) {
		return NotInGraphError{Node: to}
	}
	for i := range from.outputs {
		if from.outputs[i].end == to {
			return nil
		}
	}
	g.dirty = true
	from.outputs = append(from.outputs, connection{from, to})
	to.inputs = append(to.inputs, connection{to, from})
	return nil
}

// RemoveConnection removes a connection between two Nodes.
func (g *Graph) RemoveConnection(from, to *Node) error {
	if !g.Exists(from) {
		return NotInGraphError{Node: from}
	}
	if !g.Exists(to) {
		return NotInGraphError{Node: to}
	}

	removed, err := from.removeOutConnection(to)
	if err != nil {
		return err
	}
	g.maybeDirty(removed)

	removed, err = to.removeInConnection(from)
	if err != nil {
		return err
	}
	g.maybeDirty(removed)

	return nil
}

// Exists checks whether the Node exists in the graph.
func (g *Graph) Exists(n *Node) bool {
	return n != nil && n.idx < len(g.nodes) && g.nodes[n.idx] == n
}

// Sorted returns a topologically sorted list of strongly connected components in the Graph.
func (g *Graph) Sorted() [][]*Node {
	var (
		components = g.components[:0]
		sorted     = g.topSort()
	)
	for i := range sorted {
		sorted[i].searchState = unseen
	}
	for _, sink := range sorted {
		if g.nodes[sink.idx].searchState == unseen {
			var group []*Node
			g.dfsInputs(g.nodes[sink.idx], &group)
			components = append(components, group)
		}
	}
	return components
}

// TopologicalSort returns a list of nodes in topological order.
func (g *Graph) topSort() []*Node {
	for i := range g.nodes {
		g.nodes[i].searchState = unseen
	}
	sorted := g.sorted[:0]
	for _, node := range g.nodes {
		if node.searchState == unseen {
			g.dfsOutputs(node, &sorted)
		}
	}
	reverseNodes(sorted)
	return sorted
}

func (g *Graph) dfsOutputs(node *Node, list *[]*Node) {
	node.searchState = seen
	for _, output := range node.outputs {
		if output.end.searchState == unseen {
			g.dfsOutputs(output.end, list)
		}
	}
	*list = append(*list, node)
}

func (g *Graph) dfsInputs(node *Node, list *[]*Node) {
	node.searchState = seen
	for _, input := range node.inputs {
		if input.end.searchState == unseen {
			g.dfsInputs(input.end, list)
		}
	}
	*list = append(*list, node)
}

func (g *Graph) maybeDirty(changed bool) {
	if changed {
		g.dirty = true
	}
}

// HasChanged returns whether or not the Graph state is dirty.
func (g *Graph) HasChanged() bool {
	return g.dirty
}

// AckChange acknowledges changes to the Graph, resetting the dirty flag.
func (g *Graph) AckChange() {
	g.dirty = false
}

// NotInGraphError is an error that will be returned if a Node cannot be found in a Graph
type NotInGraphError struct {
	Node *Node
}

func (n NotInGraphError) Error() string {
	return fmt.Sprintf("node: %v not in graph", n.Node)
}
