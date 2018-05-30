package unit

import (
	"github.com/brettbuddin/shaden/graph"
)

// Patch connects a one Unit's Out to another's In. It also creates an edge on a graph to track the connection.
func Patch(g *graph.Graph, out Output, in *In) error {
	if err := g.NewConnection(out.Out().node, in.node); err != nil {
		return err
	}
	in.Couple(out)
	return nil
}

// Unpatch disconnects all inbound neighbors (Outs) from an In. All graph edges are removed as well to track the
// disconnection. Once all, if any, Outs are disconnected the In is reset to its default value constant.
func Unpatch(g *graph.Graph, in *In) error {
	if in.HasSource() {
		inputs := in.node.InNeighbors()
		for _, n := range inputs {
			if err := g.RemoveConnection(n, in.node); err != nil {
				return err
			}
		}
	}
	in.Reset()
	return nil
}
