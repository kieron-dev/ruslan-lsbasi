// Package lexer separates file bytes into tokens
package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TokenType int

const (
	Unknown TokenType = iota
	Number
	Plus
	Minus
	Mult
	Div
	LParen
	RParen
	EOF
	Begin
	End
	ID
	Dot
	Semi
	Assign
)

func (tt TokenType) String() string {
	return []string{
		"Unknown",
		"Number",
		"Plus",
		"Minus",
		"Multiply",
		"Divide",
		"left paren",
		"right paren",
		"EOF",
		"begin",
		"end",
		"ID",
		"dot",
		"semicolon",
		"assignment",
	}[tt]
}

type Token struct {
	Type  TokenType
	Value interface{}
}

type Tokeniser struct {
	currentToken Token
	buf          *bufio.Reader
}

var reservedWords = map[string]TokenType{
	"BEGIN": Begin,
	"END":   End,
	"DIV":   Div,
}

func NewTokeniser(data io.Reader) *Tokeniser {
	return &Tokeniser{
		buf: bufio.NewReader(data),
	}
}

func (t *Tokeniser) NextToken() (Token, error) {
	c := byte(' ')
	for c == ' ' || c == '\t' || c == '\n' {
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

	var byte2 byte
	nextByte, err := t.buf.Peek(1)
	if err != nil && err != io.EOF {
		return Token{}, fmt.Errorf("error peeking ahead: %w", err)
	}

	if err == nil {
		byte2 = nextByte[0]
	}

	var token Token

	switch {
	case c == '+':
		token = Token{
			Type:  Plus,
			Value: c,
		}

	case c == '-':
		token = Token{
			Type:  Minus,
			Value: c,
		}

	case c == '*':
		token = Token{
			Type:  Mult,
			Value: c,
		}

	case c == '(':
		token = Token{
			Type:  LParen,
			Value: c,
		}

	case c == ')':
		token = Token{
			Type:  RParen,
			Value: c,
		}

	case c == '.':
		token = Token{
			Type:  Dot,
			Value: c,
		}

	case c == ';':
		token = Token{
			Type:  Semi,
			Value: c,
		}

	case c == ':' && byte2 == '=':
		token = Token{
			Type:  Assign,
			Value: ":=",
		}
		_, err := t.buf.Discard(1)
		if err != nil && err != io.EOF {
			return Token{}, fmt.Errorf("trying to discard next byte: %w", err)
		}

	case c >= '0' && c <= '9':
		n, err := t.readNumber(c)
		if err != nil {
			return Token{}, fmt.Errorf("error getting number: %w", err)
		}

		token = Token{
			Type:  Number,
			Value: n,
		}

	case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
		id, err := t.readID(c)
		if err != nil {
			return Token{}, fmt.Errorf("error geting id: %w", err)
		}

		if tokenType, ok := reservedWords[strings.ToUpper(id)]; ok {
			return Token{Type: tokenType, Value: id}, nil
		}

		return Token{Type: ID, Value: id}, nil
	}

	if token.Type == Unknown {
		return token, fmt.Errorf("unexpected character: %q", c)
	}

	t.currentToken = token

	return token, nil
}

func (t *Tokeniser) readID(c byte) (string, error) {
	var s string
	for (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		s += string(c)

		var err error
		c, err = t.buf.ReadByte()
		if err != nil {
			if err != io.EOF {
				return "", err
			}

			return s, nil
		}
	}

	t.buf.UnreadByte()

	return s, nil
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
