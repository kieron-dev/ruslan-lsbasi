// Package parser implements an arithmetic parser
package parser

import (
	"errors"

	"github.com/kieron-dev/lsbasi/lexer"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Tokeniser

type Tokeniser interface {
	NextToken() (lexer.Token, error)
}

type Parser struct {
	tokeniser    Tokeniser
	currentToken lexer.Token
}

func NewParser(tokeniser Tokeniser) *Parser {
	return &Parser{
		tokeniser: tokeniser,
	}
}

func (p *Parser) NextToken() (lexer.Token, error) {
	token, err := p.tokeniser.NextToken()
	p.currentToken = token

	return token, err
}

func (p *Parser) Expr() (ASTNode, error) {
	val, err := p.Term()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == lexer.PLUS || p.currentToken.Type == lexer.MINUS {
		op := p.currentToken

		nextVal, err := p.Term()
		if err != nil {
			return nil, err
		}

		val = &BinOpNode{Left: val, Right: nextVal, Token: op}
	}

	return val, nil
}

func (p *Parser) Term() (ASTNode, error) {
	val, err := p.Factor()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == lexer.MULT || p.currentToken.Type == lexer.DIV {
		op := p.currentToken

		nextVal, err := p.Factor()
		if err != nil {
			return nil, err
		}

		val = &BinOpNode{Left: val, Right: nextVal, Token: op}
	}

	return val, nil
}

func (p *Parser) Factor() (ASTNode, error) {
	token, err := p.NextToken()
	if err != nil {
		return nil, err
	}

	if token.Type == lexer.PLUS || token.Type == lexer.MINUS {
		factor, err := p.Factor()
		if err != nil {
			return nil, err
		}

		return &UnaryNode{Token: token, Child: factor}, nil
	}

	if token.Type == lexer.LPAREN {
		val, err := p.Expr()
		if err != nil {
			return nil, err
		}

		if p.currentToken.Type != lexer.RPAREN {
			return nil, errors.New("expected closing parenthesis")
		}

		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

		return val, nil
	}

	if token.Type != lexer.NUMBER {
		return nil, errors.New("expected a left parenthesis or a number")
	}

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	return &NumNode{Value: token.Value.(int)}, nil
}
