package main

import (
	"fmt"
	"os"

	"github.com/kieron-dev/lsbasi/interpreter"
	"github.com/kieron-dev/lsbasi/lexer"
	"github.com/kieron-dev/lsbasi/parser"
)

func main() {
	pars := parser.NewParser(lexer.NewTokeniser(os.Stdin))
	interp := interpreter.NewInterpreter(pars)
	err := interp.Interpret()
	if err != nil {
		fmt.Printf("invalid expression: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("result: %#v\n", interp.GlobalScope())
}
