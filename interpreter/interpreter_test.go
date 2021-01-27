package interpreter_test

import (
	"github.com/kieron-dev/lsbasi/interpreter"
	"github.com/kieron-dev/lsbasi/interpreter/interpreterfakes"
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interpreter", func() {
	DescribeTable("expressions", func(expr parser.ASTNode, expectedValue int) {
		ast := &parser.CompoundNode{
			Children: []parser.ASTNode{
				&parser.AssignNode{
					Left:  &parser.VarNode{Value: "res"},
					Right: expr,
				},
			},
		}

		pars := new(interpreterfakes.FakeProgrammer)
		pars.ProgramReturns(ast, nil)

		interp := interpreter.NewInterpreter(pars)
		err := interp.Interpret()
		Expect(err).NotTo(HaveOccurred())

		Expect(interp.GlobalScope()["res"]).To(Equal(expectedValue))
	},

		Entry("2+5",
			&parser.BinOpNode{
				Left:  &parser.NumNode{Value: 2},
				Right: &parser.NumNode{Value: 5},
				Token: lexer.Token{Type: lexer.Plus, Value: byte('+')},
			},
			7,
		),

		Entry("2+5*3",
			&parser.BinOpNode{
				Left: &parser.NumNode{Value: 2},
				Right: &parser.BinOpNode{
					Left:  &parser.NumNode{Value: 5},
					Right: &parser.NumNode{Value: 3},
					Token: lexer.Token{Type: lexer.Mult, Value: byte('*')},
				},
				Token: lexer.Token{Type: lexer.Plus, Value: byte('+')},
			},
			17,
		),

		Entry("-5",
			&parser.UnaryNode{
				Child: &parser.NumNode{Value: 5},
				Token: lexer.Token{Type: lexer.Minus, Value: byte('-')},
			},
			-5,
		),
	)

	DescribeTable("programs", func(program *parser.CompoundNode, expectedValue map[string]int) {
		pars := new(interpreterfakes.FakeProgrammer)
		pars.ProgramReturns(program, nil)
		interp := interpreter.NewInterpreter(pars)
		err := interp.Interpret()
		Expect(err).NotTo(HaveOccurred())

		Expect(interp.GlobalScope()).To(Equal(expectedValue))
	},

		Entry("empty",
			&parser.CompoundNode{
				Children: []parser.ASTNode{
					&parser.NoOpNode{},
				},
			},
			map[string]int{},
		),

		Entry("var assignment",
			&parser.CompoundNode{
				Children: []parser.ASTNode{
					&parser.AssignNode{
						Left: &parser.VarNode{
							Value: "a",
						},
						Right: &parser.NumNode{
							Value: 42,
						},
					},
				},
			},
			map[string]int{"a": 42},
		),

		Entry("a := 42; b := a - 1",
			&parser.CompoundNode{
				Children: []parser.ASTNode{
					&parser.AssignNode{
						Left: &parser.VarNode{
							Value: "a",
						},
						Right: &parser.NumNode{
							Value: 42,
						},
					},
					&parser.AssignNode{
						Left: &parser.VarNode{
							Value: "b",
						},
						Right: &parser.BinOpNode{
							Left: &parser.VarNode{
								Value: "a",
							},
							Right: &parser.NumNode{
								Value: 1,
							},
							Token: lexer.Token{Type: lexer.Minus},
						},
					},
				},
			},
			map[string]int{"a": 42, "b": 41},
		),
	)
})
