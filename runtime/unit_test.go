package runtime

import (
	"log"
	"os"
	"testing"
	"time"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/engine"
	"buddin.us/shaden/lisp"
	"buddin.us/shaden/unit"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBuildFuncDefinition(t *testing.T) {
	var (
		env         = lisp.NewEnvironment()
		builders    = unit.Builders()
		noopBuilder = builders["noop"]
		be          = backend{}
		messages    = messageChannel{make(chan *engine.Message)}
		eng, err    = engine.New(be, engine.WithMessageChannel(messages))
		logger      = log.New(os.Stdout, "", -1)
	)

	defineBuildFunc(env, noopBuilder, eng, logger, "unit/noop")

	fn, err := env.GetSymbol("unit/noop")
	require.NoError(t, err)

	v, err := fn.(func(lisp.List) (interface{}, error))(lisp.List{})
	require.NoError(t, err)

	lu := v.(*lazyUnit)
	require.Equal(t, "noop-1", lu.id)
	require.Equal(t, []string{"x"}, lu.inputs)
	require.Equal(t, []string{"out"}, lu.outputs)
	require.Equal(t, false, lu.mount)
	require.Equal(t, "noop", lu.typ)
	require.Equal(t, logger, lu.logger)
	require.Equal(t, "noop-1", lu.created.ID)
	require.Equal(t, "noop-1(mounted=false)", lu.String())
}

func TestUnitInspection(t *testing.T) {
	var (
		env         = lisp.NewEnvironment()
		builders    = unit.Builders()
		noopBuilder = builders["noop"]
		be          = backend{}
		messages    = messageChannel{make(chan *engine.Message)}
		eng, err    = engine.New(be, engine.WithMessageChannel(messages))
		logger      = log.New(os.Stdout, "", -1)
	)

	defineBuildFunc(env, noopBuilder, eng, logger, "unit/noop")

	fn, err := env.GetSymbol("unit/noop")
	require.NoError(t, err)

	v, err := fn.(func(lisp.List) (interface{}, error))(lisp.List{})
	require.NoError(t, err)

	typ, err := unitTypeFn(lisp.List{v})
	require.NoError(t, err)
	require.Equal(t, "noop", typ)

	id, err := unitIDFn(lisp.List{v})
	require.NoError(t, err)
	require.Equal(t, "noop-3", id)

	inputs, err := unitInputsFn(lisp.List{v})
	require.NoError(t, err)
	require.Equal(t, []string{"x"}, inputs)

	outputs, err := unitOutputsFn(lisp.List{v})
	require.NoError(t, err)
	require.Equal(t, []string{"out"}, outputs)
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
		lisp.List{
			lisp.Table{"y": 10},
			lisp.Table{lisp.Keyword("z-keyword"): 11},
			lisp.Table{lisp.Symbol("z-symbol"): 12},
		},
		lisp.List{12, 13},
	})
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"w":           1,
		"x-keyword":   2,
		"x-symbol":    3,
		"0/y":         10,
		"1/z-keyword": 11,
		"2/z-symbol":  12,
		"0":           12,
		"1":           13,
	}, m)
}

func TestPatching(t *testing.T) {
	var (
		env         = lisp.NewEnvironment()
		builders    = unit.Builders()
		noopBuilder = builders["noop"]
		be          = backend{2} // Execute the callback twice
		messages    = messageChannel{make(chan *engine.Message)}
		eng, err    = engine.New(be, engine.WithMessageChannel(messages))
		logger      = log.New(os.Stdout, "", -1)
	)

	require.NoError(t, err)

	defineBuildFunc(env, noopBuilder, eng, logger, "unit/noop")

	fn, err := env.GetSymbol("unit/noop")
	require.NoError(t, err)

	v, err := fn.(func(lisp.List) (interface{}, error))(lisp.List{})
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		v, err = patchFn(eng, logger, false)(lisp.List{
			v,
			lisp.Table{"x": 1},
		})
		require.NoError(t, err)
		require.IsType(t, &lazyUnit{}, v)
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

type backend struct {
	calls int
}

func (b backend) Start(cb func([]float32, [][]float32)) error {
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
