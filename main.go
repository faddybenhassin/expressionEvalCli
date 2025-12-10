package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
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

	res, err := eval(expression)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(res)
}

func tokenizer(expression string) ([]string, error) {
	return strings.Split(expression, " "), nil
}

func infixToPostfix(expression string) ([]string, error) {
	tokens, err := tokenizer(expression)
	if err != nil {
		return nil, err
	}
	output := []string{}
	stack := []string{}

	getPrecedence := func(op string) (int, error) {
		switch op {
		case "+", "-":
			return 1, nil
		case "*", "/":
			return 2, nil
		case "^":
			return 3, nil
		default:
			return 0, fmt.Errorf("unknown operator %q", op)
		}
	}

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, token)
			continue
		}

		if token == "(" {
			stack = append(stack, token)
			continue
		}

		if token == ")" {
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == "(" {
					found = true
					stack = stack[:len(stack)-1]
					break
				}
				output = append(output, top)
				stack = stack[:len(stack)-1]
			}
			if !found {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			continue
		}

		for len(stack) > 0 {
			lastElemPrecdence, err := getPrecedence(stack[len(stack)-1])
			if err != nil {
				return nil, err
			}
			tokenPrecdence, err := getPrecedence(token)
			if err != nil {
				return nil, err
			}
			if lastElemPrecdence > tokenPrecdence || (lastElemPrecdence == tokenPrecdence && token != "^") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			} else {
				break
			}
		}
		stack = append(stack, token)

	}
	for i := len(stack) - 1; i >= 0; i-- {
		output = append(output, stack[i])
	}

	return output, nil
}

func eval(expression string) (float64, error) {
	if strings.TrimSpace(expression) == "" {
		return 0, fmt.Errorf("empty expression")
	}
	tokens, err := infixToPostfix(expression)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}
	stack := []float64{}
	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
			continue
		}
		// otherwise its an operator
		if len(stack) < 2 {
			return 0, fmt.Errorf("operator \"%q\" requires two operands", token)
		}
		b := stack[len(stack)-1] // pop the last two numbers
		a := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		switch token {
		case "+":
			stack = append(stack, a+b)
		case "-":
			stack = append(stack, a-b)
		case "*":
			stack = append(stack, a*b)
		case "/":
			if b == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			stack = append(stack, a/b)
		case "^":
			stack = append(stack, math.Pow(a, b))
		default:
			return 0, fmt.Errorf("unknown operator %q", token)

		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression: leftover values in stack")
	}

	return stack[0], nil
}
