// Package interpreter interprets ASTs
package interpreter

import (
	"log"

	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Expresser

type Expresser interface {
	Expr() (parser.ASTNode, error)
}

type Interpreter struct {
	pars Expresser
}

func NewInterpreter(pars Expresser) *Interpreter {
	return &Interpreter{
		pars: pars,
	}
}

func (i *Interpreter) Interpret() (int, error) {
	ast, err := i.pars.Expr()
	if err != nil {
		return 0, err
	}
	return ast.Accept(i).(int), nil
}

func (i *Interpreter) VisitNum(node *parser.NumNode) interface{} {
	return node.Value
}

func (i *Interpreter) VisitBinOp(node *parser.BinOpNode) interface{} {
	left := node.Left.Accept(i).(int)
	right := node.Right.Accept(i).(int)

	switch node.Token.Type {
	case lexer.Plus:
		return left + right
	case lexer.Minus:
		return left - right
	case lexer.Mult:
		return left * right
	case lexer.Div:
		return left / right
	}

	log.Fatalf("weird bin op node: %v", *node)
	return nil
}

func (i *Interpreter) VisitUnary(node *parser.UnaryNode) interface{} {
	child := node.Child.Accept(i).(int)

	if node.Token.Type == lexer.Minus {
		return -child
	}

	return child
}
