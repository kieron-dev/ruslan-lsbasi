package lexer_test

import (
	"strings"

	"github.com/kieron-dev/lsbasi/lexer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("individual tokens", func(expr string, t lexer.TokenType, v interface{}) {
	tokeniser := lexer.NewTokeniser(strings.NewReader(expr))
	token, err := tokeniser.NextToken()
	Expect(err).NotTo(HaveOccurred())
	Expect(token.Type).To(Equal(t))
	if v != nil {
		Expect(token.Value).To(Equal(v))
	}
},

	Entry("multi-digit number", "31", lexer.NUMBER, 31),
	Entry("minus", "-", lexer.MINUS, nil),
	Entry("mult", "*", lexer.MULT, nil),
	Entry("div", "/", lexer.DIV, nil),
	Entry("lparen", "(", lexer.LPAREN, nil),
	Entry("rparen", ")", lexer.RPAREN, nil),
	Entry("begin", "BEGIN", lexer.BEGIN, nil),
	Entry("end", "END", lexer.END, nil),
	Entry("an ID", "foo8", lexer.ID, "foo8"),
	Entry("dot", ".", lexer.DOT, nil),
	Entry("semi", ";", lexer.SEMI, nil),
	Entry("assignment", ":=", lexer.ASSIGN, nil),
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

		Context("invalid input", func() {
			BeforeEach(func() {
				expr = "_asdf"
			})

			It("returns an error from NextToken()", func() {
				_, err := tokeniser.NextToken()
				Expect(err).To(MatchError(ContainSubstring("unexpected character: '_'")))
			})
		})

		Describe("NextToken", func() {
			BeforeEach(func() {
				expr = "BEGIN a := 3 * - 9; END"
			})

			It("next functions as expected", func() {
				expected := []lexer.Token{
					{Type: lexer.BEGIN, Value: "BEGIN"},
					{Type: lexer.ID, Value: "a"},
					{Type: lexer.ASSIGN, Value: ":="},
					{Type: lexer.NUMBER, Value: 3},
					{Type: lexer.MULT, Value: byte('*')},
					{Type: lexer.MINUS, Value: byte('-')},
					{Type: lexer.NUMBER, Value: 9},
					{Type: lexer.SEMI, Value: byte(';')},
					{Type: lexer.END, Value: "END"},
				}

				for _, e := range expected {
					t, err := tokeniser.NextToken()
					Expect(err).NotTo(HaveOccurred())
					Expect(t).To(Equal(e))
				}
			})
		})
	})
})
