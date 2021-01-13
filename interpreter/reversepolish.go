package interpreter

import (
	"fmt"
	"strconv"

	"github.com/kieron-dev/lsbasi/parser"
)

type ReversePolish struct {
	pars Expresser
}

func NewReversePolish(pars Expresser) ReversePolish {
	return ReversePolish{
		pars: pars,
	}
}

func (rp ReversePolish) Interpret() (string, error) {
	ast, err := rp.pars.Expr()
	if err != nil {
		return "", err
	}

	return ast.Accept(rp).(string), nil
}

func (rp ReversePolish) VisitBinOp(node *parser.BinOpNode) interface{} {
	return fmt.Sprintf("%s %s %s",
		node.Left.Accept(rp),
		node.Right.Accept(rp),
		string(node.Token.Value.(byte)),
	)
}

func (rp ReversePolish) VisitNum(node *parser.NumNode) interface{} {
	return strconv.Itoa(node.Value)
}
