package unit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProp_ValueInStringList(t *testing.T) {
	prop := Prop{setter: inStringList([]string{"x", "y", "z"})}
	require.NoError(t, prop.SetValue("x"))
	require.Error(t, prop.SetValue("a"))
}

func TestProp_ValueClampedRange(t *testing.T) {
	prop := Prop{setter: clampRange(1, 10)}
	require.NoError(t, prop.SetValue(1))
	require.Equal(t, 1.0, prop.value)
	require.NoError(t, prop.SetValue(10))
	require.Equal(t, 10.0, prop.value)
	require.NoError(t, prop.SetValue(11))
	require.Equal(t, 10.0, prop.value)
	require.NoError(t, prop.SetValue(-2))
	require.Equal(t, 1.0, prop.value)
	require.Error(t, prop.SetValue("badtype"))
}
