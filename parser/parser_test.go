package parser_test

import (
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
	"github.com/kieron-dev/lsbasi/parser/parserfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parser", func() {
	var (
		pars      *parser.Parser
		tokeniser *parserfakes.FakeTokeniser
		tokens    []lexer.Token
	)

	BeforeEach(func() {
		tokeniser = new(parserfakes.FakeTokeniser)

		tokenPos := -1
		tokeniser.NextTokenStub = func() (lexer.Token, error) {
			tokenPos++
			if tokenPos >= len(tokens) {
				return lexer.Token{Type: lexer.EOF}, nil
			}

			return tokens[tokenPos], nil
		}
	})

	JustBeforeEach(func() {
		pars = parser.NewParser(tokeniser)
	})

	Describe("addition", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.NUMBER, Value: 3},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Value: 8},
			}
		})

		It("can add two numbers", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.PLUS))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("multiple addition", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.NUMBER, Value: 3},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Value: 8},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Value: 6},
			}
		})

		It("can add two numbers", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			binOp2, ok := binOp.Left.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())

			Expect(binOp.Token.Type).To(Equal(lexer.PLUS))
			Expect(binOp2.Token.Type).To(Equal(lexer.PLUS))
			Expect(binOp2.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp2.Right.(*parser.NumNode).Value).To(Equal(8))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(6))
		})
	})

	Describe("subtraction", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.NUMBER, Value: 3},
				{Type: lexer.MINUS},
				{Type: lexer.NUMBER, Value: 8},
			}
		})

		It("can subtract two numbers", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.MINUS))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("multiplication", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.NUMBER, Value: 3},
				{Type: lexer.MULT},
				{Type: lexer.NUMBER, Value: 8},
			}
		})

		It("can multiply two numbers", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.MULT))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("division", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.NUMBER, Value: 25},
				{Type: lexer.DIV},
				{Type: lexer.NUMBER, Value: 8},
			}
		})

		It("can multiply two numbers", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.DIV))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(25))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("parentheses", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.NUMBER, Value: 25},
				{Type: lexer.MINUS},
				{Type: lexer.LPAREN},
				{Type: lexer.NUMBER, Value: 5},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Value: 6},
				{Type: lexer.RPAREN},
			}
		})

		It("does the brackets first", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())

			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			binOp2, ok := binOp.Right.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())

			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(25))
			Expect(binOp.Token.Type).To(Equal(lexer.MINUS))
			Expect(binOp2.Token.Type).To(Equal(lexer.PLUS))
			Expect(binOp2.Left.(*parser.NumNode).Value).To(Equal(5))
			Expect(binOp2.Right.(*parser.NumNode).Value).To(Equal(6))
		})
	})

	Context("invalid input", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.NUMBER, Value: 3},
				{Type: lexer.MULT},
			}
		})

		It("errors", func() {
			_, err := pars.Expr()
			Expect(err).To(MatchError(ContainSubstring("expected a left parenthesis or a number")))
		})
	})
})
