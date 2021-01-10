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

		Context("MULT", func() {
			BeforeEach(func() {
				expr = "*"
			})

			It("recognises *", func() {
				token, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token.Type).To(Equal(lexer.MULT))
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

		Describe("workflow", func() {
			BeforeEach(func() {
				expr = "3 + 5 * 9 - 2"
			})

			It("next functions as expected", func() {
				token1, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token1.Type).To(Equal(lexer.NUMBER))
				Expect(token1.Value).To(Equal(3))

				token2, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token2.Type).To(Equal(lexer.PLUS))

				for i := 0; i < 5; i++ {
					_, err := tokeniser.NextToken()
					Expect(err).NotTo(HaveOccurred())
				}

				tokenEOF, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(tokenEOF.Type).To(Equal(lexer.EOF))

				tokenEOF2, err := tokeniser.NextToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(tokenEOF2.Type).To(Equal(lexer.EOF))
			})
		})
	})
})
