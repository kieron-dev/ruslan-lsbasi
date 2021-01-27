package parser_test

import (
	"errors"

	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
	"github.com/kieron-dev/lsbasi/parser/parserfakes"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type parserFn int

const (
	expr parserFn = iota + 1
	program
)

var _ = DescribeTable("tokens to AST",
	func(fn parserFn, tokens []lexer.Token, expectedAST parser.ASTNode, expectedErr error) {
		tokeniser := new(parserfakes.FakeTokeniser)

		tokenPos := -1
		tokeniser.NextTokenStub = func() (lexer.Token, error) {
			tokenPos++
			if tokenPos >= len(tokens) {
				return lexer.Token{Type: lexer.EOF}, nil
			}

			return tokens[tokenPos], nil
		}

		pars := parser.NewParser(tokeniser)
		var val parser.ASTNode
		var err error

		switch fn {
		case expr:
			// line it up
			pars.NextToken()
			val, err = pars.Expr()
		case program:
			val, err = pars.Program()
		}

		if expectedErr != nil {
			Expect(err).To(MatchError(err))
			return
		}

		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal(expectedAST))
	},

	// expressions

	Entry("3+8",
		expr,
		[]lexer.Token{
			{Type: lexer.Number, Value: 3},
			{Type: lexer.Plus},
			{Type: lexer.Number, Value: 8},
		},
		&parser.BinOpNode{
			Left:  &parser.NumNode{Value: 3},
			Right: &parser.NumNode{Value: 8},
			Token: lexer.Token{Type: lexer.Plus},
		},
		nil,
	),

	Entry("3+8+6",
		expr,
		[]lexer.Token{
			{Type: lexer.Number, Value: 3},
			{Type: lexer.Plus},
			{Type: lexer.Number, Value: 8},
			{Type: lexer.Plus},
			{Type: lexer.Number, Value: 6},
		},
		&parser.BinOpNode{
			Left: &parser.BinOpNode{
				Left:  &parser.NumNode{Value: 3},
				Right: &parser.NumNode{Value: 8},
				Token: lexer.Token{Type: lexer.Plus},
			},
			Right: &parser.NumNode{Value: 6},
			Token: lexer.Token{Type: lexer.Plus},
		},
		nil,
	),

	Entry("3-8",
		expr,
		[]lexer.Token{
			{Type: lexer.Number, Value: 3},
			{Type: lexer.Minus},
			{Type: lexer.Number, Value: 8},
		},
		&parser.BinOpNode{
			Left:  &parser.NumNode{Value: 3},
			Right: &parser.NumNode{Value: 8},
			Token: lexer.Token{Type: lexer.Minus},
		},
		nil,
	),

	Entry("3*8",
		expr,
		[]lexer.Token{
			{Type: lexer.Number, Value: 3},
			{Type: lexer.Mult},
			{Type: lexer.Number, Value: 8},
		},
		&parser.BinOpNode{
			Left:  &parser.NumNode{Value: 3},
			Right: &parser.NumNode{Value: 8},
			Token: lexer.Token{Type: lexer.Mult},
		},
		nil,
	),

	Entry("3/8",
		expr,
		[]lexer.Token{
			{Type: lexer.Number, Value: 3},
			{Type: lexer.Div},
			{Type: lexer.Number, Value: 8},
		},
		&parser.BinOpNode{
			Left:  &parser.NumNode{Value: 3},
			Right: &parser.NumNode{Value: 8},
			Token: lexer.Token{Type: lexer.Div},
		},
		nil,
	),

	Entry("25-(5+6)",
		expr,
		[]lexer.Token{
			{Type: lexer.Number, Value: 25},
			{Type: lexer.Minus},
			{Type: lexer.LParen},
			{Type: lexer.Number, Value: 5},
			{Type: lexer.Plus},
			{Type: lexer.Number, Value: 6},
			{Type: lexer.RParen},
		},
		&parser.BinOpNode{
			Left: &parser.NumNode{Value: 25},
			Right: &parser.BinOpNode{
				Left:  &parser.NumNode{Value: 5},
				Right: &parser.NumNode{Value: 6},
				Token: lexer.Token{Type: lexer.Plus},
			},
			Token: lexer.Token{Type: lexer.Minus},
		},
		nil,
	),

	Entry("-5",
		expr,
		[]lexer.Token{
			{Type: lexer.Minus},
			{Type: lexer.Number, Value: 5},
		},
		&parser.UnaryNode{
			Child: &parser.NumNode{Value: 5},
			Token: lexer.Token{Type: lexer.Minus},
		},
		nil,
	),

	Entry("+5",
		expr,
		[]lexer.Token{
			{Type: lexer.Plus},
			{Type: lexer.Number, Value: 5},
		},
		&parser.UnaryNode{
			Child: &parser.NumNode{Value: 5},
			Token: lexer.Token{Type: lexer.Plus},
		},
		nil,
	),

	// full programs

	Entry(`
BEGIN
	bob := 2;
	res := bob;
END.
`,
		program,
		[]lexer.Token{
			{Type: lexer.Begin},
			{Type: lexer.ID, Value: "bob"},
			{Type: lexer.Assign},
			{Type: lexer.Number, Value: 2},
			{Type: lexer.Semi},
			{Type: lexer.ID, Value: "res"},
			{Type: lexer.Assign},
			{Type: lexer.ID, Value: "bob"},
			{Type: lexer.Semi},
			{Type: lexer.End},
			{Type: lexer.Dot},
			{Type: lexer.EOF},
		},
		&parser.CompoundNode{
			Children: []parser.ASTNode{
				&parser.AssignNode{
					Left:  &parser.VarNode{Value: "bob"},
					Right: &parser.NumNode{Value: 2},
				},
				&parser.AssignNode{
					Left:  &parser.VarNode{Value: "res"},
					Right: &parser.VarNode{Value: "bob"},
				},
				&parser.NoOpNode{},
			},
		},
		nil,
	),

	// errors

	Entry("5+",
		expr,
		[]lexer.Token{
			{Type: lexer.Number, Value: 5},
			{Type: lexer.Plus},
		},
		nil,
		errors.New("expected a left parenthesis or a number"),
	),
)
