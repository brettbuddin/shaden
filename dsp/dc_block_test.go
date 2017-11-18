package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDCBlocking(t *testing.T) {
	block := &DCBlock{}
	require.Equal(t, 0.5, block.Tick(0.5))
	require.Equal(t, -0.5025, block.Tick(-0.5))
}
