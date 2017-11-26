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
	require.Equal(t, 3, g.Size())
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

	unknown := &Node{}

	require.NoError(t, g.NewConnection(a, b))
	require.NoError(t, g.NewConnection(b, c))
	require.Error(t, g.NewConnection(unknown, b))
	require.Error(t, g.NewConnection(b, unknown))

	require.Equal(t, 1, a.NeighborCount())
	require.Equal(t, 0, a.InNeighborCount())
	require.Equal(t, 1, a.OutNeighborCount())
	require.Equal(t, b, a.OutNeighbors()[0])

	require.Equal(t, 2, b.NeighborCount())
	require.Equal(t, 1, b.InNeighborCount())
	require.Equal(t, a, b.InNeighbors()[0])
	require.Equal(t, 1, b.OutNeighborCount())
	require.Equal(t, c, b.OutNeighbors()[0])

	require.Equal(t, 1, c.NeighborCount())
	require.Equal(t, 1, c.InNeighborCount())
	require.Equal(t, b, c.InNeighbors()[0])
	require.Equal(t, 0, c.OutNeighborCount())

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

	require.NoError(t, g.NewConnection(a, b))
	require.Equal(t, 1, a.OutNeighborCount())
	require.Equal(t, 1, b.InNeighborCount())

	path := []string{}
	for _, n := range g.topSort() {
		path = append(path, n.Value.(string))
	}
	require.Equal(t, []string{"a", "b"}, path)

	unknown := &Node{}
	require.Error(t, g.RemoveConnection(a, unknown))
	require.Error(t, g.RemoveConnection(unknown, a))

	require.NoError(t, g.RemoveConnection(a, b))
	require.Equal(t, 0, a.OutNeighborCount())
	require.Equal(t, 0, b.InNeighborCount())

	require.NoError(t, g.RemoveConnection(b, a))

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

	require.NoError(t, g.NewConnection(a, b))
	require.NoError(t, g.NewConnection(b, c))
	require.NoError(t, g.NewConnection(c, a))

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
	require.NoError(b, err)
	err = g.NewConnection(bn, c)
	require.NoError(b, err)
	err = g.NewConnection(c, a)
	require.NoError(b, err)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.Sorted()
	}
}

func TestDirty(t *testing.T) {
	g := New()
	require.False(t, g.HasChanged())
	a := g.NewNode("a")
	b := g.NewNode("b")
	require.False(t, g.HasChanged())
	require.NoError(t, g.NewConnection(a, b))
	require.True(t, g.HasChanged())
	g.AckChange()
	require.False(t, g.HasChanged())
}
