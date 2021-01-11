package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kieron-dev/lsbasi/arithmetic"
	"github.com/kieron-dev/lsbasi/lexer"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("expr> ")
	for scanner.Scan() {

		line := scanner.Text()
		line = strings.TrimSpace(line)

		interp := arithmetic.NewInterpreter(lexer.NewTokeniser(strings.NewReader(line)))
		val, err := interp.Expr()
		if err != nil {
			fmt.Printf("invalid expression: %q\n", line)
			fmt.Print("expr> ")
			continue
		}

		fmt.Printf("%d\n", val)
		fmt.Print("expr> ")
	}
}