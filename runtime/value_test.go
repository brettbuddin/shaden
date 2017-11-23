package runtime

import (
	"testing"

	"buddin.us/musictheory"
	"buddin.us/shaden/dsp"
	"buddin.us/shaden/lisp"
	"github.com/stretchr/testify/require"
)

func TestHz(t *testing.T) {
	v, err := hzFn(lisp.List{440.0})
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	v, err = hzFn(lisp.List{440})
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	v, err = hzFn(lisp.List{"A4"})
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	v, err = hzFn(lisp.List{lisp.Keyword("A4")})
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	pitch, err := musictheory.ParsePitch("A4")
	require.NoError(t, err)
	v, err = hzFn(lisp.List{pitch})
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	_, err = hzFn(lisp.List{"111"})
	require.Error(t, err)

	_, err = hzFn(lisp.List{})
	require.Error(t, err)
}

func TestMS(t *testing.T) {
	v, err := msFn(lisp.List{1})
	require.NoError(t, err)
	require.Equal(t, 44.1, v.(dsp.Valuer).Float64())

	v, err = msFn(lisp.List{1.0})
	require.NoError(t, err)
	require.Equal(t, 44.1, v.(dsp.Valuer).Float64())

	_, err = msFn(lisp.List{"1"})
	require.Error(t, err)

	_, err = msFn(lisp.List{})
	require.Error(t, err)
}

func TestBPM(t *testing.T) {
	v, err := bpmFn(lisp.List{1})
	require.NoError(t, err)
	require.Equal(t, 3.779289493575208e-07, v.(dsp.Valuer).Float64())

	v, err = bpmFn(lisp.List{1.0})
	require.NoError(t, err)
	require.Equal(t, 3.779289493575208e-07, v.(dsp.Valuer).Float64())

	_, err = bpmFn(lisp.List{"aoeu"})
	require.Error(t, err)

	_, err = bpmFn(lisp.List{})
	require.Error(t, err)
}

func TestDB(t *testing.T) {
	v, err := dbFn(lisp.List{0})
	require.NoError(t, err)
	require.Equal(t, 1.0, v)

	v, err = dbFn(lisp.List{-6})
	require.NoError(t, err)
	require.Equal(t, 0.5011872336272722, v)

	v, err = dbFn(lisp.List{-12})
	require.NoError(t, err)
	require.Equal(t, 0.25118864315095796, v)

	_, err = dbFn(lisp.List{"-6"})
	require.Error(t, err)

	_, err = dbFn(lisp.List{})
	require.Error(t, err)
}
