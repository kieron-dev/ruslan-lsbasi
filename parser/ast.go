package parser

import "github.com/kieron-dev/lsbasi/lexer"

type Visitor interface {
	VisitBinOp(*BinOpNode) (interface{}, error)
	VisitNum(*NumNode) (interface{}, error)
	VisitUnary(*UnaryNode) (interface{}, error)
	VisitCompound(*CompoundNode) (interface{}, error)
	VisitAssign(*AssignNode) (interface{}, error)
	VisitVar(*VarNode) (interface{}, error)
	VisitNoOp(*NoOpNode) (interface{}, error)
}

type ASTNode interface {
	Accept(Visitor) (interface{}, error)
}

type BinOpNode struct {
	Left  ASTNode
	Right ASTNode
	Token lexer.Token
}

func (n *BinOpNode) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinOp(n)
}

type NumNode struct {
	Token lexer.Token
	Value int
}

func (n *NumNode) Accept(v Visitor) (interface{}, error) {
	return v.VisitNum(n)
}

type UnaryNode struct {
	Token lexer.Token
	Child ASTNode
}

func (n *UnaryNode) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnary(n)
}

type CompoundNode struct {
	Children []ASTNode
}

func (n *CompoundNode) Accept(v Visitor) (interface{}, error) {
	return v.VisitCompound(n)
}

type AssignNode struct {
	Left  *VarNode
	Right ASTNode
}

func (n *AssignNode) Accept(v Visitor) (interface{}, error) {
	return v.VisitAssign(n)
}

type VarNode struct {
	Value string
}

func (n *VarNode) Accept(v Visitor) (interface{}, error) {
	return v.VisitVar(n)
}

type NoOpNode struct{}

func (n *NoOpNode) Accept(v Visitor) (interface{}, error) {
	return v.VisitNoOp(n)
}
