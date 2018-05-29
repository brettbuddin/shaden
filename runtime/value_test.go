package runtime

import (
	"log"
	"os"
	"testing"

	"github.com/brettbuddin/shaden/dsp"
	"github.com/brettbuddin/shaden/engine"
	"github.com/stretchr/testify/require"
)

func TestHz(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(&backend{
			sampleRate: sampleRate,
			frameSize:  frameSize,
		}, frameSize, engine.WithMessageChannel(messages))
		logger = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	v, err := run.Eval([]byte(`(hz 440)`))
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	v, err = run.Eval([]byte(`(hz 440.0)`))
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	v, err = run.Eval([]byte(`(hz "A4")`))
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	v, err = run.Eval([]byte(`(hz :A4)`))
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	v, err = run.Eval([]byte(`(hz (theory/pitch "A4"))`))
	require.NoError(t, err)
	require.Equal(t, 0.009977324263038548, v.(dsp.Valuer).Float64())

	_, err = run.Eval([]byte(`(hz "111")`))
	require.Error(t, err)

	_, err = run.Eval([]byte(`(hz)`))
	require.Error(t, err)
}

func TestMS(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(&backend{
			sampleRate: sampleRate,
			frameSize:  frameSize,
		}, frameSize, engine.WithMessageChannel(messages))
		logger = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	v, err := run.Eval([]byte(`(ms 1)`))
	require.NoError(t, err)
	require.Equal(t, 44.1, v.(dsp.Valuer).Float64())

	v, err = run.Eval([]byte(`(ms 1.0)`))
	require.NoError(t, err)
	require.Equal(t, 44.1, v.(dsp.Valuer).Float64())

	_, err = run.Eval([]byte(`(ms "1")`))
	require.Error(t, err)

	_, err = run.Eval([]byte(`(ms)`))
	require.Error(t, err)
}

func TestBPM(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(&backend{
			sampleRate: sampleRate,
			frameSize:  frameSize,
		}, frameSize, engine.WithMessageChannel(messages))
		logger = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	v, err := run.Eval([]byte(`(bpm 60)`))
	require.NoError(t, err)
	require.Equal(t, 2.2675736961451248e-05, v.(dsp.Valuer).Float64())

	v, err = run.Eval([]byte(`(bpm 60.0)`))
	require.NoError(t, err)
	require.Equal(t, 2.2675736961451248e-05, v.(dsp.Valuer).Float64())

	_, err = run.Eval([]byte(`(bpm "1")`))
	require.Error(t, err)

	_, err = run.Eval([]byte(`(bpm)`))
	require.Error(t, err)
}

func TestDB(t *testing.T) {
	var (
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(&backend{
			sampleRate: sampleRate,
			frameSize:  frameSize,
		}, frameSize, engine.WithMessageChannel(messages))
		logger = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)

	v, err := run.Eval([]byte(`(db 0)`))
	require.NoError(t, err)
	require.Equal(t, 1.0, v)

	v, err = run.Eval([]byte(`(db -6)`))
	require.NoError(t, err)
	require.Equal(t, 0.5011872336272722, v)

	v, err = run.Eval([]byte(`(db -6.0)`))
	require.NoError(t, err)
	require.Equal(t, 0.5011872336272722, v)

	v, err = run.Eval([]byte(`(db -12)`))
	require.Equal(t, 0.25118864315095796, v)

	_, err = run.Eval([]byte(`(db "0")`))
	require.Error(t, err)

	_, err = run.Eval([]byte(`(db)`))
	require.Error(t, err)
}
