package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kieron-dev/lsbasi/interpreter"
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("expr> ")
	for scanner.Scan() {

		line := scanner.Text()
		line = strings.TrimSpace(line)

		pars := parser.NewParser(lexer.NewTokeniser(strings.NewReader(line)))
		interp := interpreter.NewInterpreter(pars)
		val, err := interp.Interpret()
		if err != nil {
			fmt.Printf("invalid expression: %q\n", line)
			fmt.Print("expr> ")
			continue
		}

		rpn, _ := interpreter.NewReversePolish(parser.NewParser(lexer.NewTokeniser(strings.NewReader(line)))).Interpret()
		lisp, _ := interpreter.NewLisp(parser.NewParser(lexer.NewTokeniser(strings.NewReader(line)))).Interpret()

		fmt.Printf("result: %d\n", val)
		fmt.Printf("rpn: %s\n", rpn)
		fmt.Printf("lisp: %s\n\n", lisp)
		fmt.Print("expr> ")
	}
}
