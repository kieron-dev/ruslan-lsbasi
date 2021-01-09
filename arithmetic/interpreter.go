// Package arithmetic implements an arithmetic interpreter
package arithmetic

import (
	"errors"
	"fmt"

	"github.com/kieron-dev/lsbasi/lexer"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Tokeniser

type Tokeniser interface {
	CurrentToken() lexer.Token
	NextToken() (lexer.Token, error)
	Eat(lexer.TokenType) error
}

type Interpreter struct {
	tokeniser Tokeniser
}

func NewInterpreter(tokeniser Tokeniser) Interpreter {
	return Interpreter{
		tokeniser: tokeniser,
	}
}

func (i Interpreter) Expr() (int, error) {
	token, err := i.tokeniser.NextToken()
	if err != nil {
		return 0, err
	}

	leftVal := token.Value
	if err := i.tokeniser.Eat(lexer.NUMBER); err != nil {
		return 0, fmt.Errorf("invalid expression")
	}

	operation := i.tokeniser.CurrentToken()
	if err := i.tokeniser.Eat(operation.Type); err != nil {
		return 0, fmt.Errorf("invalid expression")
	}

	token = i.tokeniser.CurrentToken()
	rightVal := token.Value
	if err := i.tokeniser.Eat(lexer.NUMBER); err != nil {
		return 0, fmt.Errorf("invalid expression")
	}

	if operation.Type == lexer.PLUS {
		return leftVal.(int) + rightVal.(int), nil
	}

	if operation.Type == lexer.MINUS {
		return leftVal.(int) - rightVal.(int), nil
	}

	return 0, errors.New("unexpected op")
}
