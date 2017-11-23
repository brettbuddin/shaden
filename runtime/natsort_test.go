package runtime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNaturalSort(t *testing.T) {
	slice := []string{
		"a",
		"A",
		"b",
		"B",
		"11",
		"11",
		"11a11",
		"1",
		"10",
		"a03",
		"a10",
		"a1",
		"/b",
	}
	natsort(slice)
	require.Equal(t, []string{
		"/b",
		"1",
		"10",
		"11",
		"11",
		"11a11",
		"A",
		"B",
		"a",
		"a1",
		"a03",
		"a10",
		"b",
	}, slice)
}
