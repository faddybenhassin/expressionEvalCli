package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Eval(expression string) (float64, error) {
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
			return 0, fmt.Errorf("operator %q requires two operands", token)
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
		case "u-":
			stack = append(stack, a, -b)
		case "^":
			stack = append(stack, math.Pow(a, b))
		default:
			return 0, fmt.Errorf("unknown operator %q", token)

		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression leftover values in stack")
	}

	return stack[0], nil
}
