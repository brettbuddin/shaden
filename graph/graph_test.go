package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNode(t *testing.T) {
	g := New()

	n1 := g.NewNode("oscillator")
	n2 := g.NewNode("filter")
	n3 := g.NewNode("sink")

	require.True(t, g.Exists(n1))
	require.True(t, g.Exists(n2))
	require.True(t, g.Exists(n3))
}

func TestNodeRemoval(t *testing.T) {
	g := New()

	a := g.NewNode("a")
	b := g.NewNode("b")
	g.NewNode("c")

	require.Nil(t, g.NewConnection(a, b))
	require.Equal(t, 1, a.OutNeighborCount())
	require.Equal(t, 1, b.InNeighborCount())

	path := []string{}
	for _, n := range g.topSort() {
		path = append(path, n.Value.(string))
	}
	require.Equal(t, []string{"c", "a", "b"}, path)

	require.Nil(t, g.RemoveNode(a))
	require.Equal(t, 0, b.InNeighborCount())
	require.NotNil(t, g.RemoveNode(&Node{}))

	path = []string{}
	for _, n := range g.topSort() {
		path = append(path, n.Value.(string))
	}
	require.Equal(t, []string{"c", "b"}, path)
}

func TestNewConnection(t *testing.T) {
	g := New()

	a := g.NewNode("a")
	b := g.NewNode("b")
	c := g.NewNode("c")
	g.NewNode("d")

	require.Nil(t, g.NewConnection(a, b))
	require.Nil(t, g.NewConnection(b, c))

	path := []string{}
	for _, n := range g.topSort() {
		path = append(path, n.Value.(string))
	}
	require.Equal(t, []string{"d", "a", "b", "c"}, path)
}

func TestConnectionRemoval(t *testing.T) {
	g := New()

	a := g.NewNode("a")
	b := g.NewNode("b")

	require.Nil(t, g.NewConnection(a, b))
	require.Equal(t, 1, a.OutNeighborCount())
	require.Equal(t, 1, b.InNeighborCount())

	path := []string{}
	for _, n := range g.topSort() {
		path = append(path, n.Value.(string))
	}
	require.Equal(t, []string{"a", "b"}, path)

	require.Nil(t, g.RemoveConnection(a, b))
	require.Equal(t, 0, a.OutNeighborCount())
	require.Equal(t, 0, b.InNeighborCount())

	require.Nil(t, g.RemoveConnection(b, a))
	require.NotNil(t, g.RemoveConnection(b, &Node{}))

	path = []string{}
	for _, n := range g.topSort() {
		path = append(path, n.Value.(string))
	}
	require.Equal(t, []string{"b", "a"}, path)
}

func TestSCC(t *testing.T) {
	g := New()

	a := g.NewNode("a")
	b := g.NewNode("b")
	c := g.NewNode("c")
	g.NewNode("d")

	require.Nil(t, g.NewConnection(a, b))
	require.Nil(t, g.NewConnection(b, c))
	require.Nil(t, g.NewConnection(c, a))

	scc := g.Sorted()
	path := make([][]string, len(scc))
	for i, g := range scc {
		for _, n := range g {
			path[i] = append(path[i], n.Value.(string))
		}
	}
	require.Equal(t, 2, len(path))
	require.Equal(t, path[0], []string{"d"})
	require.Equal(t, path[1], []string{"b", "c", "a"})
}

func BenchmarkSCC(b *testing.B) {
	g := New()

	a := g.NewNode("a")
	bn := g.NewNode("b")
	c := g.NewNode("c")
	g.NewNode("d")

	err := g.NewConnection(a, bn)
	require.Nil(b, err)
	err = g.NewConnection(bn, c)
	require.Nil(b, err)
	err = g.NewConnection(c, a)
	require.Nil(b, err)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.Sorted()
	}
}
