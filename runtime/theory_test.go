package runtime

import (
	"testing"

	"buddin.us/musictheory"
	"buddin.us/shaden/lisp"
	"github.com/stretchr/testify/require"
)

func TestPitch(t *testing.T) {
	v, err := pitchFn(lisp.List{"A4"})
	require.NoError(t, err)
	require.Equal(t, 440.0, v.(musictheory.Pitch).Freq())

	_, err = pitchFn(lisp.List{"111"})
	require.Error(t, err)

	_, err = pitchFn(lisp.List{})
	require.Error(t, err)
}

func TestInterval(t *testing.T) {
	v, err := intervalFn(lisp.List{"perfect", 1})
	require.NoError(t, err)
	require.Equal(t, musictheory.Perfect(1), v)

	_, err = pitchFn(lisp.List{"111"})
	require.Error(t, err)

	_, err = pitchFn(lisp.List{})
	require.Error(t, err)
}

func TestTranspose(t *testing.T) {
	pitch, err := pitchFn(lisp.List{"A4"})
	require.NoError(t, err)

	intvl, err := intervalFn(lisp.List{"minor", 3})
	require.NoError(t, err)
	require.Equal(t, musictheory.Minor(3), intvl)

	v, err := transposeFn(lisp.List{pitch, intvl})
	require.NoError(t, err)

	actualPitch, err := musictheory.ParsePitch("C5")
	require.NoError(t, err)
	require.True(t, actualPitch.Eq(v.(musictheory.Pitch)))

	_, err = pitchFn(lisp.List{"111"})
	require.Error(t, err)

	_, err = pitchFn(lisp.List{})
	require.Error(t, err)
}
