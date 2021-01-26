package parser_test

import (
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
	"github.com/kieron-dev/lsbasi/parser/parserfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parser AST construction", func() {
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
				{Type: lexer.Number, Value: 3},
				{Type: lexer.Plus},
				{Type: lexer.Number, Value: 8},
			}
		})

		It("creates a BinOpNode with 3 and 8", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.Plus))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("multiple addition", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Number, Value: 3},
				{Type: lexer.Plus},
				{Type: lexer.Number, Value: 8},
				{Type: lexer.Plus},
				{Type: lexer.Number, Value: 6},
			}
		})

		It("has a BinOp with BinOp(3, 8) and 6 as children", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			binOp2, ok := binOp.Left.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())

			Expect(binOp.Token.Type).To(Equal(lexer.Plus))
			Expect(binOp2.Token.Type).To(Equal(lexer.Plus))
			Expect(binOp2.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp2.Right.(*parser.NumNode).Value).To(Equal(8))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(6))
		})
	})

	Describe("subtraction", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Number, Value: 3},
				{Type: lexer.Minus},
				{Type: lexer.Number, Value: 8},
			}
		})

		It("creates a minus BinOp with 3 and 8 as children", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.Minus))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("multiplication", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Number, Value: 3},
				{Type: lexer.Mult},
				{Type: lexer.Number, Value: 8},
			}
		})

		It("creates a Mult BinOp with 3 and 8 as children", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.Mult))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(3))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("division", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Number, Value: 25},
				{Type: lexer.Div},
				{Type: lexer.Number, Value: 8},
			}
		})

		It("creates a Div BinOp with 25 and 8 as children", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())
			binOp, ok := val.(*parser.BinOpNode)
			Expect(ok).To(BeTrue())
			Expect(binOp.Token.Type).To(Equal(lexer.Div))
			Expect(binOp.Left.(*parser.NumNode).Value).To(Equal(25))
			Expect(binOp.Right.(*parser.NumNode).Value).To(Equal(8))
		})
	})

	Describe("parentheses", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Number, Value: 25},
				{Type: lexer.Minus},
				{Type: lexer.LParen},
				{Type: lexer.Number, Value: 5},
				{Type: lexer.Plus},
				{Type: lexer.Number, Value: 6},
				{Type: lexer.RParen},
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
			Expect(binOp.Token.Type).To(Equal(lexer.Minus))
			Expect(binOp2.Token.Type).To(Equal(lexer.Plus))
			Expect(binOp2.Left.(*parser.NumNode).Value).To(Equal(5))
			Expect(binOp2.Right.(*parser.NumNode).Value).To(Equal(6))
		})
	})

	Describe("unary minus", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Minus},
				{Type: lexer.Number, Value: 5},
			}
		})

		It("gives a unary minus with a 5", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())

			op, ok := val.(*parser.UnaryNode)
			Expect(ok).To(BeTrue())
			Expect(op.Token.Type).To(Equal(lexer.Minus))
			Expect(op.Child.(*parser.NumNode).Value).To(Equal(5))
		})
	})

	Describe("unary plus", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Plus},
				{Type: lexer.Number, Value: 4},
			}
		})

		It("gives a unary plus with a 4", func() {
			val, err := pars.Expr()
			Expect(err).NotTo(HaveOccurred())

			op, ok := val.(*parser.UnaryNode)
			Expect(ok).To(BeTrue())
			Expect(op.Token.Type).To(Equal(lexer.Plus))
			Expect(op.Child.(*parser.NumNode).Value).To(Equal(4))
		})
	})

	Context("invalid input", func() {
		BeforeEach(func() {
			tokens = []lexer.Token{
				{Type: lexer.Number, Value: 3},
				{Type: lexer.Mult},
			}
		})

		It("errors", func() {
			_, err := pars.Expr()
			Expect(err).To(MatchError(ContainSubstring("expected a left parenthesis or a number")))
		})
	})
})
