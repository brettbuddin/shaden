package lisp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvironment_Lookup(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", 42)
	v, err := env.GetSymbol("hello")
	require.NoError(t, err)
	require.Equal(t, 42, v)
}

func TestEnvironment_LookupUndefined(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", 42)
	_, err := env.GetSymbol("helloworld")
	require.Error(t, err)
}

func TestEnvironment_Set(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", 42)
	err := env.SetSymbol("hello", 41)
	require.NoError(t, err)
	v, err := env.GetSymbol("hello")
	require.NoError(t, err)
	require.Equal(t, 41, v)
}

func TestEnvironment_Unset(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", 42)
	err := env.UnsetSymbol("hello")
	require.NoError(t, err)
	_, err = env.GetSymbol("hello")
	require.Error(t, err)
}

func TestEnvironment_UnsetUndefined(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", 42)
	err := env.UnsetSymbol("helloworld")
	require.Error(t, err)
}

func TestEnvironment_SetUndefined(t *testing.T) {
	env := NewEnvironment()
	err := env.SetSymbol("hello", 42)
	require.Error(t, err)
}

func TestEnvironment_SymbolAlreadyDefined(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", 42)
	v, err := env.GetSymbol("hello")
	require.NoError(t, err)
	require.Equal(t, 42, v)

	env.DefineSymbol("hello", 41)
	v, err = env.GetSymbol("hello")
	require.NoError(t, err)
	require.Equal(t, 41, v)
}

func TestEnvironment_BranchParentLookup(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", 42)
	env.DefineSymbol("world", 2)

	env = env.Branch()

	env.DefineSymbol("hello", 41)

	v, err := env.GetSymbol("hello")
	require.NoError(t, err)
	require.Equal(t, 41, v)

	v, err = env.GetSymbol("world")
	require.NoError(t, err)
	require.Equal(t, 2, v)
}

func TestEnvironment_FunctionCall(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", func(args List) (any, error) {
		return "hello " + args[0].(string), nil
	})

	v, err := env.Eval(List{Symbol("hello"), "world"})
	require.NoError(t, err)
	require.Equal(t, "hello world", v)
}

func TestEnvironment_FunctionInterfaceCall(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", function(func(args List) (any, error) {
		return "hello " + args[0].(string), nil
	}))

	v, err := env.Eval(List{Symbol("hello"), "world"})
	require.NoError(t, err)
	require.Equal(t, "hello world", v)
}

func TestEnvironment_EnvFunctionCall(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", func(env *Environment, args List) (any, error) {
		v, err := env.Eval(args[0])
		if err != nil {
			return nil, err
		}
		return "hello " + v.(string), nil
	})

	v, err := env.Eval(List{Symbol("hello"), "world"})
	require.NoError(t, err)
	require.Equal(t, "hello world", v)
}

func TestEnvironment_EnvFunctionInterfaceCall(t *testing.T) {
	env := NewEnvironment()
	env.DefineSymbol("hello", envFunction(func(env *Environment, args List) (any, error) {
		v, err := env.Eval(args[0])
		if err != nil {
			return nil, err
		}
		return "hello " + v.(string), nil
	}))

	v, err := env.Eval(List{Symbol("hello"), "world"})
	require.NoError(t, err)
	require.Equal(t, "hello world", v)
}

func TestEnvironment_ValueIdentity(t *testing.T) {
	tests := []struct {
		input  any
		result any
	}{
		{nil, nil},
		{1, 1},
		{1.0, 1.0},
		{"hello", "hello"},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.input), func(t *testing.T) {
			env := NewEnvironment()
			v, err := env.Eval(test.input)
			require.NoError(t, err)
			require.Equal(t, test.result, v)
		})
	}
}

type function func(List) (any, error)

func (function) Name() string { return "function" }

func (f function) Func(l List) (any, error) {
	return f(l)
}

type envFunction func(*Environment, List) (any, error)

func (envFunction) Name() string { return "env function" }

func (f envFunction) EnvFunc(env *Environment, l List) (any, error) {
	return f(env, l)
}
