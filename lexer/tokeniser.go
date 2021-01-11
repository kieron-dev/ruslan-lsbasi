// Package lexer separates file bytes into tokens
package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type TokenType int

const (
	UNKNOWN TokenType = iota
	NUMBER
	PLUS
	MINUS
	MULT
	DIV
	LPAREN
	RPAREN
	EOF
)

type Token struct {
	Type  TokenType
	Value interface{}
}

type Tokeniser struct {
	currentToken Token
	buf          *bufio.Reader
}

func NewTokeniser(data io.Reader) *Tokeniser {
	return &Tokeniser{
		buf: bufio.NewReader(data),
	}
}

func (t *Tokeniser) NextToken() (Token, error) {
	c := byte(' ')
	for c == ' ' {
		var err error
		c, err = t.buf.ReadByte()
		if err != nil {
			if err != io.EOF {
				return Token{}, fmt.Errorf("read error getting next char: %w", err)
			}

			return Token{
				Type: EOF,
			}, nil
		}
	}

	var token Token

	switch {
	case c == '+':
		token = Token{
			Type:  PLUS,
			Value: c,
		}

	case c == '-':
		token = Token{
			Type:  MINUS,
			Value: c,
		}

	case c == '*':
		token = Token{
			Type:  MULT,
			Value: c,
		}

	case c == '/':
		token = Token{
			Type:  DIV,
			Value: c,
		}

	case c == '(':
		token = Token{
			Type:  LPAREN,
			Value: c,
		}

	case c == ')':
		token = Token{
			Type:  RPAREN,
			Value: c,
		}

	case c >= '0' && c <= '9':
		n, err := t.readNumber(c)
		if err != nil {
			return Token{}, fmt.Errorf("error getting number: %w", err)
		}

		token = Token{
			Type:  NUMBER,
			Value: n,
		}
	}

	if token.Type == UNKNOWN {
		return token, fmt.Errorf("unexpected character: %q", c)
	}

	t.currentToken = token

	return token, nil
}

func (t *Tokeniser) readNumber(c byte) (int, error) {
	var s string
	for c >= '0' && c <= '9' {
		s += string(c)

		var err error
		c, err = t.buf.ReadByte()
		if err != nil {
			if err != io.EOF {
				return 0, err
			}

			return strconv.Atoi(s)
		}
	}

	t.buf.UnreadByte()

	return strconv.Atoi(s)
}
