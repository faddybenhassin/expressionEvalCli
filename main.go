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

	fmt.Println("Calculator started. Type 'exit' or 'quit' to stop.")

	for {
		fmt.Print("> ") // Visual cue for the user
		if !scanner.Scan() {
			break // Exit if the scanner hits an error or EOF
		}

		expression := strings.TrimSpace(scanner.Text())

		// Handle empty input or exit commands
		if expression == "" {
			continue
		}
		if expression == "exit" || expression == "quit" {
			break
		}

		res, err := eval.Eval(expression)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			// We use continue so one error doesn't crash the whole loop
			continue
		}

		fmt.Printf("Result: %v\n", res)
	}

	fmt.Println("Exiting...")
}
