// Package arithmetic implements an arithmetic interpreter
package arithmetic

import (
	"errors"

	"github.com/kieron-dev/lsbasi/lexer"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Tokeniser

type Tokeniser interface {
	NextToken() (lexer.Token, error)
}

type Interpreter struct {
	tokeniser    Tokeniser
	currentToken lexer.Token
}

func NewInterpreter(tokeniser Tokeniser) *Interpreter {
	return &Interpreter{
		tokeniser: tokeniser,
	}
}

func (i *Interpreter) NextToken() (lexer.Token, error) {
	token, err := i.tokeniser.NextToken()
	i.currentToken = token

	return token, err
}

func (i *Interpreter) Expr() (int, error) {
	val, err := i.Term()
	if err != nil {
		return 0, err
	}

	for i.currentToken.Type == lexer.PLUS || i.currentToken.Type == lexer.MINUS {
		op := i.currentToken

		nextVal, err := i.Term()
		if err != nil {
			return 0, err
		}

		if op.Type == lexer.PLUS {
			val += nextVal
		} else {
			val -= nextVal
		}
	}

	return val, nil
}

func (i *Interpreter) Term() (int, error) {
	token, err := i.NextToken()
	if err != nil {
		return 0, err
	}

	if token.Type != lexer.NUMBER {
		return 0, errors.New("expected a number")
	}

	val := token.Value.(int)

	_, err = i.NextToken()
	if err != nil {
		return 0, err
	}

	for i.currentToken.Type == lexer.MULT || i.currentToken.Type == lexer.DIV {
		op := i.currentToken

		next, err := i.NextToken()
		if err != nil {
			return 0, err
		}

		if next.Type != lexer.NUMBER {
			return 0, errors.New("expected a number")
		}

		if op.Type == lexer.MULT {
			val *= next.Value.(int)
		}

		if _, err = i.NextToken(); err != nil {
			return 0, err
		}
	}

	return val, nil
}
