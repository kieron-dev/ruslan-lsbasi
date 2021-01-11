package integration_test

import (
	"strings"

	"github.com/kieron-dev/lsbasi/arithmetic"
	"github.com/kieron-dev/lsbasi/lexer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	DescribeTable("expressions", func(expr string, res int) {
		tokeniser := lexer.NewTokeniser(strings.NewReader(expr))
		interpreter := arithmetic.NewInterpreter(tokeniser)
		out, err := interpreter.Expr()
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(res))
	},

		Entry("single number", "9", 9),
		Entry("single bigger number", "371", 371),
		Entry("simple addition", "3 + 5", 8),
		Entry("more addition", "3 + 5 + 51", 59),
		Entry("simple subtraction", "31 - 5", 26),
		Entry("add and subtract", "3 + 5 - 51", -43),
		Entry("simple multiplication", "3 *  5", 15),
		Entry("multiple multiplication", "3 * 4 * 5", 60),
		Entry("simple division", "25 /  5", 5),
		Entry("simple division", "3*24/6", 12),
		Entry("precedence I", "3 + 4 * 5", 23),
		Entry("precedence II", "3 * 4 + 5", 17),
		Entry("precedence III", "3 + 24 / 8", 6),
		Entry("precedence IV", "(3 + 24) / 9", 3),
		Entry("precedence V", "(163 + 17) / (9 * 10)", 2),
		Entry("precedence VI", "7 + 3 * (10 / (12 / (3 + 1) - 1))", 22),
	)
})
