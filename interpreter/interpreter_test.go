package interpreter_test

import (
	"github.com/kieron-dev/lsbasi/interpreter"
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interpreter", func() {
	var (
		ast    parser.ASTNode
		interp *interpreter.Interpreter
	)

	BeforeEach(func() {
		ast2 := &parser.NumNode{Value: 2}
		ast5 := &parser.NumNode{Value: 5}
		ast = &parser.BinOpNode{Left: ast2, Right: ast5, Token: lexer.Token{Type: lexer.PLUS}}

		interp = interpreter.NewInterpreter(ast)
	})

	It("does 2 + 5", func() {
		Expect(interp.Interpret()).To(Equal(7))
	})
})
