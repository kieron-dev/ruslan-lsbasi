package interpreter

import (
	"fmt"
	"strconv"

	"github.com/kieron-dev/lsbasi/parser"
)

type Lisp struct {
	pars Expresser
}

func NewLisp(pars Expresser) Lisp {
	return Lisp{
		pars: pars,
	}
}

func (l Lisp) VisitBinOp(node *parser.BinOpNode) interface{} {
	return fmt.Sprintf("(%s %s %s)",
		string(node.Token.Value.(byte)),
		node.Left.Accept(l).(string),
		node.Right.Accept(l).(string),
	)
}

func (l Lisp) VisitNum(node *parser.NumNode) interface{} {
	return strconv.Itoa(node.Value)
}

func (l Lisp) Interpret() (string, error) {
	ast, err := l.pars.Expr()
	if err != nil {
		return "", err
	}

	return ast.Accept(l).(string), nil
}
