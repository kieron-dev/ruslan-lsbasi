// Package interpreter interprets ASTs
package interpreter

import (
	"fmt"
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
	_, err = expr.Accept(i)

	return err
}

func (i *Interpreter) VisitNum(node *parser.NumNode) (interface{}, error) {
	return node.Value, nil
}

func (i *Interpreter) VisitBinOp(node *parser.BinOpNode) (interface{}, error) {
	leftVal, err := node.Left.Accept(i)
	if err != nil {
		return nil, err
	}
	rightVal, err := node.Right.Accept(i)
	if err != nil {
		return nil, err
	}
	left := leftVal.(int)
	right := rightVal.(int)

	switch node.Token.Type {
	case lexer.Plus:
		return left + right, nil
	case lexer.Minus:
		return left - right, nil
	case lexer.Mult:
		return left * right, nil
	case lexer.Div:
		return left / right, nil
	}

	return nil, fmt.Errorf("weird bin op node: %v", *node)
}

func (i *Interpreter) VisitUnary(node *parser.UnaryNode) (interface{}, error) {
	child, err := node.Child.Accept(i)
	if err != nil {
		return nil, err
	}
	val := child.(int)

	if node.Token.Type == lexer.Minus {
		return -val, nil
	}

	return val, nil
}

func (i *Interpreter) VisitCompound(node *parser.CompoundNode) (interface{}, error) {
	for _, child := range node.Children {
		_, err := child.Accept(i)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitAssign(node *parser.AssignNode) (interface{}, error) {
	varName := strings.ToLower(node.Left.Value)
	value, err := node.Right.Accept(i)
	if err != nil {
		return nil, err
	}
	i.globalSymbols[varName] = value.(int)

	return nil, nil
}

func (i *Interpreter) VisitVar(node *parser.VarNode) (interface{}, error) {
	varName := strings.ToLower(node.Value)
	val, ok := i.globalSymbols[varName]
	if !ok {
		return nil, fmt.Errorf("unknown var %q", node.Value)
	}

	return val, nil
}

func (i *Interpreter) VisitNoOp(node *parser.NoOpNode) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) GlobalScope() map[string]int {
	return i.globalSymbols
}
