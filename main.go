package main

import (
	"bufio"
	"expressionEvalCli/eval"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	expression := strings.TrimSpace(scanner.Text())

	if expression == "" {
		fmt.Println("Usage: <\"expression\">")
		os.Exit(1)
	}

	res, err := eval.Eval(expression)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(res)
}
