// Package lisp provides a lisp interpreter
package lisp

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/brettbuddin/shaden/errors"
)

// Parse parses lisp expressions in the content from an io.Reader.
func Parse(r io.Reader) (any, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	l := newLexer(string(buf))
	go l.run()
	return newParser(l).Parse()
}

type parser struct {
	lexer *lexer
}

func newParser(l *lexer) *parser {
	return &parser{lexer: l}
}

type root struct {
	Nodes List
}

func (p *parser) Parse() (any, error) {
	nodes := List{}
	if err := p.parse(&nodes); err != nil {
		p.lexer.drain()
		return nil, err
	}
	return &root{Nodes: nodes}, nil
}

func (p *parser) parse(tree *List) error {
	for token := range p.lexer.tokens {
		if token.typ == tokenEOF {
			return nil
		}
		switch typ := token.typ; typ {
		case tokenError:
			return newLineError(errors.New(token.value), token.line)
		case tokenLeftParen:
			nodes := List{}
			if err := p.parse(&nodes); err != nil {
				return newLineError(err, token.line)
			}
			*tree = append(*tree, nodes)
		case tokenRightParen:
			return nil
		case tokenNil:
			*tree = append(*tree, nil)
		case tokenBool:
			b, err := strconv.ParseBool(token.value)
			if err != nil {
				return newLineError(err, token.line)
			}
			*tree = append(*tree, b)
		case tokenString:
			*tree = append(*tree, strings.Trim(token.value, "\""))
		case tokenKeyword:
			*tree = append(*tree, Keyword(strings.TrimPrefix(token.value, ":")))
		case tokenSymbol:
			*tree = append(*tree, Symbol(token.value))
		case tokenFloat:
			f, err := strconv.ParseFloat(token.value, 64)
			if err != nil {
				return newLineError(err, token.line)
			}
			*tree = append(*tree, f)
		case tokenInt:
			i, err := strconv.ParseInt(token.value, 10, 64)
			if err != nil {
				return newLineError(err, token.line)
			}
			*tree = append(*tree, int(i))
		default:
			return newLineError(fmt.Errorf("unknown token type %v", token), token.line)
		}
	}
	return nil
}
