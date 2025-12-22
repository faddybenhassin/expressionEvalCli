package main

import (
	"bufio"
	"expressionEvalCli/eval"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// add asigniing variables and make it create ans variable each time
	vars := map[string]float64{
		"pi": 3.14159,
	}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Calculator started. Type 'exit' or 'quit' to stop.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 1. Handle exit and empty input
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" || input == "q" {
			break
		}

		// 2. Handle "assign <var> <val>"
		if strings.HasPrefix(input, "assign") {
			parts := strings.Fields(input) // Splits by whitespace

			if len(parts) != 3 {
				fmt.Println("Error: Invalid assign syntax. Use 'assign <name> <value>'")
				continue
			}

			varName := parts[1]
			rawVal := parts[2]

			// Evaluate the value (this allows 'assign x 5+5')
			val, err := strconv.ParseFloat(rawVal, 64)
			if err != nil {
				fmt.Printf("Error evaluating value: %q\n", rawVal)
				continue
			}

			// Store in the map
			vars[varName] = val
			fmt.Printf("Assigned %v to %s\n", val, varName)
			continue
		}

		// 3. Regular expression evaluation
		res, err := eval.Eval(input, vars)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Result: %v\n", res)
	}
}
