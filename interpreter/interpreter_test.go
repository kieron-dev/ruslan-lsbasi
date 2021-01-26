package interpreter_test

import (
	"github.com/kieron-dev/lsbasi/interpreter"
	"github.com/kieron-dev/lsbasi/interpreter/interpreterfakes"
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interpreter", func() {
	var (
		pars   *interpreterfakes.FakeExpresser
		interp *interpreter.Interpreter
		ast    parser.ASTNode
	)

	BeforeEach(func() {
		pars = new(interpreterfakes.FakeExpresser)
		ast = &parser.NumNode{Value: 10}
	})

	JustBeforeEach(func() {
		interp = interpreter.NewInterpreter(pars)
		pars.ExprReturns(ast, nil)
	})

	It("calcs a single number", func() {
		Expect(interp.Interpret()).To(Equal(10))
	})

	Context("2+5", func() {
		BeforeEach(func() {
			ast2 := &parser.NumNode{Value: 2}
			ast5 := &parser.NumNode{Value: 5}
			ast = &parser.BinOpNode{Left: ast2, Right: ast5, Token: lexer.Token{Type: lexer.Plus, Value: byte('+')}}
		})

		It("gets 7", func() {
			Expect(interp.Interpret()).To(Equal(7))
		})
	})

	Context("2+5*3", func() {
		BeforeEach(func() {
			ast2 := &parser.NumNode{Value: 2}
			ast3 := &parser.NumNode{Value: 3}
			ast5 := &parser.NumNode{Value: 5}
			astMult := &parser.BinOpNode{Left: ast5, Right: ast3, Token: lexer.Token{Type: lexer.Mult, Value: byte('*')}}
			ast = &parser.BinOpNode{Left: ast2, Right: astMult, Token: lexer.Token{Type: lexer.Plus, Value: byte('+')}}
		})

		It("gets 17", func() {
			Expect(interp.Interpret()).To(Equal(17))
		})
	})

	Context("-5", func() {
		BeforeEach(func() {
			ast5 := &parser.NumNode{Value: 5}
			ast = &parser.UnaryNode{Child: ast5, Token: lexer.Token{Type: lexer.Minus, Value: byte('-')}}
		})

		It("gets -5", func() {
			Expect(interp.Interpret()).To(Equal(-5))
		})
	})
})
