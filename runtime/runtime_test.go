package runtime

import (
	"log"
	"os"
	"testing"
	"time"

	"buddin.us/shaden/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const timeout = 10 * time.Second

func TestEnvironmentClearing(t *testing.T) {
	var (
		be       = &backend{calls: 1} // Execute the callback once
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		run, err := New(eng, logger)
		require.NoError(t, err)
		run.Eval([]byte(`
			(define noop (unit/noop))
			(clear)
		`))
		_, err = run.Eval([]byte(`noop`))
		assert.Error(t, err)
		require.NoError(t, eng.Stop())
	}()

	go func() {
		eng.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		t.Error("timeout waiting for completion")
	}
}

func TestEmitting(t *testing.T) {
	var (
		be       = &backend{calls: 3} // Execute the callback twice
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		run, err := New(eng, logger)
		require.NoError(t, err)
		run.Eval([]byte(`
			(define noop (unit/noop))
			(-> noop (table :x 1))
			(emit (<- noop))
		`))
		assert.NoError(t, err)
		left := be.written[0]
		assert.NotEqual(t, 0, left[len(left)-1])
		require.NoError(t, eng.Stop())
	}()

	go func() {
		eng.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		t.Error("timeout waiting for completion")
	}
}
