// Package interpreter interprets ASTs
package interpreter

import (
	"log"

	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
)

type Interpreter struct {
	pars *parser.Parser
}

func NewInterpreter(pars *parser.Parser) *Interpreter {
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
	case lexer.PLUS:
		return left + right
	case lexer.MINUS:
		return left - right
	case lexer.MULT:
		return left * right
	case lexer.DIV:
		return left / right
	}

	log.Fatalf("weird bin op node: %v", *node)
	return nil
}
