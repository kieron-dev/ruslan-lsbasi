package lexer_test

import (
	"strings"

	"github.com/kieron-dev/lsbasi/lexer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokeniser", func() {
	var (
		tokeniser *lexer.Tokeniser
		expr      string
	)

	BeforeEach(func() {
		expr = "3+5"
	})

	JustBeforeEach(func() {
		tokeniser = lexer.NewTokeniser(strings.NewReader(expr))
	})

	Context("tokenisation", func() {
		It("can tokenise '3+5'", func() {
			token, err := tokeniser.NextToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(lexer.Token{
				Type:  lexer.NUMBER,
				Value: 3,
			}))

			token, err = tokeniser.NextToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(lexer.Token{
				Type:  lexer.PLUS,
				Value: byte('+'),
			}))

			token, err = tokeniser.NextToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(lexer.Token{
				Type:  lexer.NUMBER,
				Value: 5,
			}))

			token, err = tokeniser.NextToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(lexer.Token{
				Type: lexer.EOF,
			}))
		})

		Context("multi-digit numbers", func() {
			BeforeEach(func() {
				expr = "31"
			})

			It("gets a number token with value 31", func() {
				token, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token.Value).To(Equal(31))
			})
		})

		Context("spaces are ignored", func() {
			BeforeEach(func() {
				expr = "31   178"
			})

			It("gets a number token with value 31", func() {
				_, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())

				token, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token.Value).To(Equal(178))
			})
		})

		Context("MINUS", func() {
			BeforeEach(func() {
				expr = "-"
			})

			It("recognises -", func() {
				token, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token.Type).To(Equal(lexer.MINUS))
			})
		})

		Context("invalid input", func() {
			BeforeEach(func() {
				expr = "asdf"
			})

			It("returns an error from NextToken()", func() {
				_, err := tokeniser.NextToken()
				Expect(err).To(MatchError(ContainSubstring("unexpected character: 'a'")))
			})
		})

		Describe("eating", func() {
			It("can eat a token type", func() {
				token, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token).To(Equal(lexer.Token{
					Type:  lexer.NUMBER,
					Value: 3,
				}))

				err = tokeniser.Eat(lexer.NUMBER)

				Expect(err).NotTo(HaveOccurred())
				Expect(tokeniser.CurrentToken()).To(Equal(lexer.Token{Type: lexer.PLUS, Value: byte('+')}))
			})

			It("errors when eating an incorrect type", func() {
				tokeniser.NextToken()
				err := tokeniser.Eat(lexer.PLUS)

				Expect(err).To(MatchError(ContainSubstring("expected current token to be of type")))
			})
		})
	})
})
