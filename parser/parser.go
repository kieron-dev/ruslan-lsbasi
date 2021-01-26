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

func (p *Parser) Program() (ASTNode, error) {
	// program : compound-statement DOT

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	val, err := p.CompoundStatement()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.Dot {
		return nil, errors.New("expected a DOT")
	}

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	return val, nil
}

func (p *Parser) CompoundStatement() (ASTNode, error) {
	// compound-statement: BEGIN statement-list END

	if p.currentToken.Type != lexer.Begin {
		return nil, errors.New("expected BEGIN")
	}

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	val, err := p.StatementList()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.End {
		return nil, errors.New("expected END")
	}

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	return val, nil
}

func (p *Parser) StatementList() (ASTNode, error) {
	// statement-list: statement
	//               | statement SEMI statement_list

	statement, err := p.Statement()
	if err != nil {
		return nil, err
	}

	val := &CompoundNode{
		Children: []ASTNode{statement},
	}

	for p.currentToken.Type == lexer.Semi {
		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

		next, err := p.Statement()
		if err != nil {
			return nil, err
		}

		val.Children = append(val.Children, next)
	}

	return val, nil
}

func (p *Parser) Statement() (ASTNode, error) {
	// statement : compound_statement
	//           | assignment_statement
	//           | empty

	if p.currentToken.Type == lexer.Begin {
		return p.CompoundStatement()
	}

	if p.currentToken.Type == lexer.ID {
		return p.AssignmentStatement()
	}

	return p.Empty()
}

func (p *Parser) AssignmentStatement() (ASTNode, error) {
	// assignment_statement : variable ASSIGN expr

	left, err := p.Variable()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.Assign {
		return nil, errors.New("expected assignment")
	}

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	right, err := p.Expr()
	if err != nil {
		return nil, err
	}

	return &AssignNode{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) Variable() (*VarNode, error) {
	// variable : ID

	if p.currentToken.Type != lexer.ID {
		return nil, errors.New("expected an ID")
	}

	node := &VarNode{
		Value: p.currentToken.Value.(string),
	}

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	return node, nil
}

func (p *Parser) Empty() (ASTNode, error) {
	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	return &NoOpNode{}, nil
}

func (p *Parser) Expr() (ASTNode, error) {
	val, err := p.Term()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == lexer.Plus || p.currentToken.Type == lexer.Minus {
		op := p.currentToken

		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

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

	for p.currentToken.Type == lexer.Mult || p.currentToken.Type == lexer.Div {
		op := p.currentToken

		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

		nextVal, err := p.Factor()
		if err != nil {
			return nil, err
		}

		val = &BinOpNode{Left: val, Right: nextVal, Token: op}
	}

	return val, nil
}

func (p *Parser) Factor() (ASTNode, error) {
	token := p.currentToken

	if token.Type == lexer.Plus || token.Type == lexer.Minus {
		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

		factor, err := p.Factor()
		if err != nil {
			return nil, err
		}

		return &UnaryNode{Token: token, Child: factor}, nil
	}

	if token.Type == lexer.LParen {
		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

		val, err := p.Expr()
		if err != nil {
			return nil, err
		}

		if p.currentToken.Type != lexer.RParen {
			return nil, errors.New("expected closing parenthesis")
		}

		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

		return val, nil
	}

	if token.Type == lexer.ID {
		return p.Variable()
	}

	if token.Type != lexer.Number {
		return nil, errors.New("expected a left parenthesis, ID or a number")
	}

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	return &NumNode{Value: token.Value.(int)}, nil
}
