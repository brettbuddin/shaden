package runtime

import (
	"log"
	"os"
	"testing"

	"github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/shaden/engine"
	"github.com/stretchr/testify/require"
)

func TestPitch(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(&backend{}, sampleRate, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	v, err := run.Eval([]byte(`(theory/pitch "A4")`))
	require.NoError(t, err)
	require.Equal(t, 440.0, v.(musictheory.Pitch).Freq())

	_, err = run.Eval([]byte(`(theory/pitch "1")`))
	require.Error(t, err)

	_, err = run.Eval([]byte(`(theory/pitch)`))
	require.Error(t, err)
}

func TestInterval(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(newBackend(0), frameSize, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	tests := []struct {
		input    string
		expected musictheory.Interval
	}{
		{"(theory/interval :perfect 1)", musictheory.Perfect(1)},
		{"(theory/interval :minor 2)", musictheory.Minor(2)},
		{"(theory/interval :major 3)", musictheory.Major(3)},
		{"(theory/interval :augmented 4)", musictheory.Augmented(4)},
		{"(theory/interval :diminished 5)", musictheory.Diminished(5)},
		{`(theory/interval "diminished" 5)`, musictheory.Diminished(5)},
	}

	for _, test := range tests {
		v, err := run.Eval([]byte(test.input))
		require.NoError(t, err)
		require.Equal(t, test.expected, v)
	}

	_, err = run.Eval([]byte(`(theory/interval :unknown)`))
	require.Error(t, err)

	_, err = run.Eval([]byte(`(theory/interval)`))
	require.Error(t, err)
}

func TestTranspose(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(newBackend(0), sampleRate, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	v, err := run.Eval([]byte(`(theory/transpose (theory/pitch "A4") (theory/interval :minor 3))`))
	require.NoError(t, err)

	actualPitch, err := musictheory.ParsePitch("C5")
	require.NoError(t, err)
	require.True(t, actualPitch.Eq(v.(musictheory.Pitch)))

	_, err = run.Eval([]byte(`(theory/pitch)`))
	require.Error(t, err)
}
