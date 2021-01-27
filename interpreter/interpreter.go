// Package interpreter interprets ASTs
package interpreter

import (
	"log"
	"strings"

	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Programmer

type Programmer interface {
	Program() (parser.ASTNode, error)
}

type Interpreter struct {
	pars          Programmer
	globalSymbols map[string]int
}

func NewInterpreter(pars Programmer) *Interpreter {
	return &Interpreter{
		pars:          pars,
		globalSymbols: map[string]int{},
	}
}

func (i *Interpreter) Interpret() error {
	expr, err := i.pars.Program()
	if err != nil {
		return err
	}
	expr.Accept(i)

	return nil
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

func (i *Interpreter) VisitCompound(node *parser.CompoundNode) interface{} {
	for _, child := range node.Children {
		child.Accept(i)
	}

	return nil
}

func (i *Interpreter) VisitAssign(node *parser.AssignNode) interface{} {
	varName := strings.ToLower(node.Left.Value)
	value := node.Right.Accept(i)
	i.globalSymbols[varName] = value.(int)

	return nil
}

func (i *Interpreter) VisitVar(node *parser.VarNode) interface{} {
	varName := strings.ToLower(node.Value)
	val, ok := i.globalSymbols[varName]
	if !ok {
		log.Printf("unknown var %q", node.Value)
	}

	return val
}

func (i *Interpreter) VisitNoOp(node *parser.NoOpNode) interface{} {
	return nil
}

func (i *Interpreter) GlobalScope() map[string]int {
	return i.globalSymbols
}
