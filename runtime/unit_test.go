package runtime

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/engine"
	"buddin.us/shaden/lisp"
	"buddin.us/shaden/unit"
)

func TestUnitOutputs(t *testing.T) {
	var (
		be       = &backend{}
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)
	v, err := run.Eval([]byte(`
		(define noop (unit/noop))
		(unit-outputs noop)
	`))
	require.NoError(t, err)
	require.Equal(t, []string{"out"}, v)
}

func TestUnitInputs(t *testing.T) {
	var (
		be       = &backend{}
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)
	v, err := run.Eval([]byte(`
		(define noop (unit/noop))
		(unit-inputs noop)
	`))
	require.NoError(t, err)
	require.Equal(t, []string{"x"}, v)
}

func TestUnitID(t *testing.T) {
	var (
		be       = &backend{}
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)
	v, err := run.Eval([]byte(`
		(define noop (unit/noop))
		(unit-id noop)
	`))
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(v.(string), "noop-"))
}

func TestUnitType(t *testing.T) {
	var (
		be       = &backend{}
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)
	run, err := New(eng, logger)
	require.NoError(t, err)
	v, err := run.Eval([]byte(`
		(define noop (unit/noop))
		(unit-type noop)
	`))
	require.NoError(t, err)
	require.Equal(t, "noop", v)
}

func TestUnitOutput(t *testing.T) {
	var (
		be       = &backend{calls: 1} // execute callback once
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		run, err := New(eng, logger)
		require.NoError(t, err)
		v, err := run.Eval([]byte(`
			(define noop (unit/noop))
			(list (<- noop) (<- noop :out))
		`))
		assert.NoError(t, err)
		list := v.(lisp.List)
		assert.Equal(t, "out", list[0].(unit.OutRef).Output)
		assert.Equal(t, "out", list[1].(unit.OutRef).Output)
		require.NoError(t, eng.Stop())
	}()

	go func() {
		eng.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for completion")
	}
}

func TestUnitPatch(t *testing.T) {
	var (
		be       = &backend{calls: 2} // execute callback twice
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		run, err := New(eng, logger)
		require.NoError(t, err)
		_, err = run.Eval([]byte(`
			(define noop (unit/noop))
			(-> noop (table :x 1))
		`))
		assert.NoError(t, err)
		require.NoError(t, eng.Stop())
	}()

	go func() {
		eng.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for completion")
	}
}

func TestUnitUnmount(t *testing.T) {
	var (
		be       = &backend{calls: 3}
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		run, err := New(eng, logger)
		require.NoError(t, err)
		v, err := run.Eval([]byte(`
			(define noop (unit/noop))
			(-> noop (table :x 1))
			(unit-unmount noop)
			noop
		`))
		assert.NoError(t, err)
		assert.False(t, v.(*lazyUnit).mount)
		require.NoError(t, eng.Stop())
	}()

	go func() {
		eng.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for completion")
	}
}

func TestUnitRemove(t *testing.T) {
	var (
		be       = &backend{calls: 3}
		messages = messageChannel{make(chan *engine.Message)}
		eng, err = engine.New(be, engine.WithMessageChannel(messages))
		logger   = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		run, err := New(eng, logger)
		require.NoError(t, err)
		_, err = run.Eval([]byte(`
			(define noop (unit/noop))
			(-> noop (table :x 1))
			(unit-remove noop)
			noop
		`))
		assert.Error(t, err)
		require.NoError(t, eng.Stop())
	}()

	go func() {
		eng.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for completion")
	}
}

func TestInputConversions(t *testing.T) {
	m, err := patchableInputs(lisp.List{})
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{}, m)

	_, err = patchableInputs(lisp.List{"invalid"})
	require.Error(t, err)

	m, err = patchableInputs(lisp.List{
		lisp.Table{"w": 1},
		lisp.Table{lisp.Keyword("x-keyword"): 2},
		lisp.Table{lisp.Symbol("x-symbol"): 3},

		// Flatten lists into map
		lisp.List{
			lisp.Table{"y": 10},
			lisp.Table{lisp.Keyword("z-keyword"): 11},
			lisp.Table{lisp.Symbol("z-symbol"): 12},
		},
		lisp.List{12, 13},

		// We only go 1 layer deep in the naming
		lisp.Table{"nested-list": lisp.List{"a", "b", "c"}},
		lisp.Table{"nested-table": lisp.Table{"a": "b"}},
	})
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"w":            1,
		"x-keyword":    2,
		"x-symbol":     3,
		"0/y":          10,
		"1/z-keyword":  11,
		"2/z-symbol":   12,
		"0":            12,
		"1":            13,
		"nested-list":  []interface{}{"a", "b", "c"},
		"nested-table": map[string]interface{}{"a": "b"},
	}, m)
}

type backend struct {
	calls   int
	written [][]float32
}

func (b *backend) Start(cb func([]float32, [][]float32)) error {
	var (
		in  = make([]float32, dsp.FrameSize)
		out = [][]float32{
			make([]float32, dsp.FrameSize),
			make([]float32, dsp.FrameSize),
		}
	)
	for i := 0; i < b.calls; i++ {
		cb(in, out)
	}
	b.written = out
	return nil
}
func (backend) Stop() error    { return nil }
func (backend) FrameSize() int { return dsp.FrameSize }

type messageChannel struct {
	messages chan *engine.Message
}

func (c messageChannel) Send(msg *engine.Message) error {
	select {
	case c.messages <- msg:
	case <-time.After(10 * time.Second):
		return errors.New("timeout sending message")
	}
	return nil
}

func (c messageChannel) Receive() *engine.Message {
	return <-c.messages
}

func (c messageChannel) Close() { close(c.messages) }
