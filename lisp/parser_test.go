package lisp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected List
		errorMsg string
	}{
		{input: "", expected: List{}},
		{input: "(list 1 2 3)", expected: List{List{Symbol("list"), 1, 2, 3}}},
		{input: "\"hello\"", expected: List{"hello"}},
		{input: "3.67", expected: List{3.67}},
		{input: "nil", expected: List{nil}},
		{input: "true", expected: List{true}},
		{input: "false", expected: List{false}},
		{input: ":hello", expected: List{Keyword("hello")}},
		{input: "((hello)", errorMsg: "unclosed left paren (line 1)"},
		{input: "(hello))", errorMsg: "unexpected right paren (line 1)"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			b := bytes.NewBufferString(test.input)
			v, err := Parse(b)
			if len(test.errorMsg) == 0 {
				require.NoError(t, err)
				require.Equal(t, newRoot(test.expected), v)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), test.errorMsg)
			}
		})
	}
}

func newRoot(nodes List) *root {
	return &root{Nodes: nodes}
}
