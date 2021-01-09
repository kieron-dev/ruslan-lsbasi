package arithmetic_test

import (
	"errors"

	"github.com/kieron-dev/lsbasi/arithmetic"
	"github.com/kieron-dev/lsbasi/arithmetic/arithmeticfakes"
	"github.com/kieron-dev/lsbasi/lexer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("intepreter", func() {
	var (
		interpreter arithmetic.Interpreter
		tokeniser   *arithmeticfakes.FakeTokeniser
	)

	BeforeEach(func() {
		tokeniser = new(arithmeticfakes.FakeTokeniser)
	})

	JustBeforeEach(func() {
		interpreter = arithmetic.NewInterpreter(tokeniser)
	})

	Describe("addition", func() {
		BeforeEach(func() {
			tokeniser.NextTokenReturns(lexer.Token{Value: 3}, nil)
			tokeniser.CurrentTokenReturnsOnCall(0, lexer.Token{Type: lexer.PLUS})
			tokeniser.CurrentTokenReturnsOnCall(1, lexer.Token{Value: 8})
		})

		It("can add two numbers", func() {
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(11))
		})
	})

	Describe("subtraction", func() {
		BeforeEach(func() {
			tokeniser.NextTokenReturns(lexer.Token{Value: 3}, nil)
			tokeniser.CurrentTokenReturnsOnCall(0, lexer.Token{Type: lexer.MINUS})
			tokeniser.CurrentTokenReturnsOnCall(1, lexer.Token{Value: 8})
		})

		It("can add two numbers", func() {
			val, err := interpreter.Expr()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(-5))
		})
	})

	Context("invalid input", func() {
		BeforeEach(func() {
			tokeniser.EatReturns(errors.New("oops"))
		})

		It("errors", func() {
			_, err := interpreter.Expr()
			Expect(err).To(MatchError(ContainSubstring("invalid expression")))
		})
	})
})
