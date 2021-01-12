package parser

import "github.com/kieron-dev/lsbasi/lexer"

type Visitor interface {
	VisitBinOp(*BinOpNode) interface{}
	VisitNum(*NumNode) interface{}
}

type ASTNode interface {
	Accept(Visitor) interface{}
}

type BinOpNode struct {
	Left  ASTNode
	Right ASTNode
	Token lexer.Token
}

func (n *BinOpNode) Accept(v Visitor) interface{} {
	return v.VisitBinOp(n)
}

type NumNode struct {
	Token lexer.Token
	Value int
}

func (n *NumNode) Accept(v Visitor) interface{} {
	return v.VisitNum(n)
}
