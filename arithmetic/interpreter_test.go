package arithmetic_test

import (
	"github.com/kieron-dev/lsbasi/arithmetic"
	"github.com/kieron-dev/lsbasi/arithmetic/arithmeticfakes"
	"github.com/kieron-dev/lsbasi/lexer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("intepreter", func() {
	var (
		interpreter *arithmetic.Interpreter
		tokeniser   *arithmeticfakes.FakeTokeniser
		tokens      []lexer.Token
	)

	BeforeEach(func() {
		tokeniser = new(arithmeticfakes.FakeTokeniser)

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
		interpreter = arithmetic.NewInterpreter(tokeniser)
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
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(11))
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
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(17))
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
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(-5))
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
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(24))
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
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(3))
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
				{Type: lexer.NUMBER, Value: 5},
				{Type: lexer.RPAREN},
			}
		})

		It("does the brackets first", func() {
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(15))
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
			_, err := interpreter.Expr()
			Expect(err).To(MatchError(ContainSubstring("expected a left parenthesis or a number")))
		})
	})
})
