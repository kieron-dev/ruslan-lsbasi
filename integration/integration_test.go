package integration_test

import (
	"fmt"
	"strings"

	"github.com/kieron-dev/lsbasi/interpreter"
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	DescribeTable("interpreting expressions", func(expr string, res int) {
		program := fmt.Sprintf(`
BEGIN
	res := %s;
END.
`, expr)
		tokeniser := lexer.NewTokeniser(strings.NewReader(program))
		pars := parser.NewParser(tokeniser)
		interp := interpreter.NewInterpreter(pars)
		err := interp.Interpret()
		Expect(err).NotTo(HaveOccurred())
		Expect(interp.GlobalScope()["res"]).To(Equal(res))
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
		Entry("unary minus", "- 5  + 3", -2),
		Entry("unary plus", "+ 5  + 3", 8),
		Entry("unary minus minus", "- - 5  + 3", 8),
		Entry("unary minus parens", "-(3+2)", -5),
	)

	DescribeTable("interpreting programs", func(program string, res map[string]int) {
		tokeniser := lexer.NewTokeniser(strings.NewReader(program))
		pars := parser.NewParser(tokeniser)
		interp := interpreter.NewInterpreter(pars)
		err := interp.Interpret()
		Expect(err).NotTo(HaveOccurred())
		Expect(interp.GlobalScope()).To(Equal(res))
	},

		Entry("empty block", "BEGIN END.", map[string]int{}),
		Entry("simple assignment", "BEGIN a := 1 END.", map[string]int{"a": 1}),
		Entry("assignment to a var", "BEGIN a := 1; b := a END.", map[string]int{"a": 1, "b": 1}),
		Entry("sample prog from chapter 9", `
BEGIN
    BEGIN
        number := 2;
        a := number;
        b := 10 * a + 10 * number / 4;
        c := a - - b
    END;
    x := 11;
END.
`,
			map[string]int{"number": 2, "a": 2, "b": 25, "c": 27, "x": 11},
		),
	)
})
