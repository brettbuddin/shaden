package lisp

import (
	"testing"

	assert "gopkg.in/go-playground/assert.v1"
)

func TestLex(t *testing.T) {
	var tests = []struct {
		input  string
		output []token
	}{
		{"", []token{newToken(tokenEOF, "")}},
		{"word", []token{newToken(tokenSymbol, "word"), newToken(tokenEOF, "")}},
		{"\"word\"", []token{newToken(tokenString, "\"word\""), newToken(tokenEOF, "")}},
		{"\"word", []token{
			newToken(tokenError, "unclosed quoted string"),
		}},
		{"1", []token{newToken(tokenInt, "1"), newToken(tokenEOF, "")}},
		{"1.05", []token{newToken(tokenFloat, "1.05"), newToken(tokenEOF, "")}},
		{"(word)", []token{
			newToken(tokenLeftParen, "("),
			newToken(tokenSymbol, "word"),
			newToken(tokenRightParen, ")"),
			newToken(tokenEOF, ""),
		}},
		{"(word", []token{
			newToken(tokenLeftParen, "("),
			newToken(tokenSymbol, "word"),
			newToken(tokenEOF, ""),
			newToken(tokenError, "unclosed left paren"),
		}},
		{")", []token{
			newToken(tokenRightParen, ")"),
			newToken(tokenError, "unexpected right paren"),
		}},
		{"(word 5)", []token{
			newToken(tokenLeftParen, "("),
			newToken(tokenSymbol, "word"),
			newToken(tokenInt, "5"),
			newToken(tokenRightParen, ")"),
			newToken(tokenEOF, ""),
		}},
		{"nil", []token{newToken(tokenNil, "nil"), newToken(tokenEOF, "")}},
		{"true", []token{newToken(tokenBool, "true"), newToken(tokenEOF, "")}},
		{"false", []token{newToken(tokenBool, "false"), newToken(tokenEOF, "")}},
		{":word", []token{newToken(tokenKeyword, ":word"), newToken(tokenEOF, "")}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := newLexer(test.input)
			go l.run()
			defer l.drain()

			var output []token
			for token := range l.tokens {
				output = append(output, token)
			}

			for i, out := range output {
				assert.Equal(t, out.typ, test.output[i].typ)
				assert.Equal(t, out.value, test.output[i].value)
			}
		})
	}
}

func newToken(typ tokenType, value string) token {
	return token{typ: typ, value: value}
}
