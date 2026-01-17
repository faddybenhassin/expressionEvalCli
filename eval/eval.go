package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Eval(expression string, vars map[string]float64, functions map[string]func(float64) float64) (float64, error) {

	// Guard clause: handle empty strings
	if strings.TrimSpace(expression) == "" {
		return 0, fmt.Errorf("empty expression")
	}

	// Step 1: Convert from human-readable (Infix) to machine-friendly (Postfix)
	tokens, err := infixToPostfix(expression, vars, functions)
	if err != nil {
		return 0, err
	}

	// The stack stores intermediate numerical results
	stack := []float64{}

	for _, token := range tokens {
		// 1. Check if token is a Number (e.g., "3.14")
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
			continue
		}

		// 2. Check if token is a Variable (e.g., "x")
		if val, ok := vars[token]; ok {
			stack = append(stack, val)
			continue
		}

		// 3. Check if token is a Function (e.g., "sin")
		if fn, ok := functions[token]; ok {
			if len(stack) < 1 {
				return 0, fmt.Errorf("stack empty for function %s", token)
			}
			// Functions are unary here: pop 1 value, apply function, push result
			val := stack[len(stack)-1]
			stack[len(stack)-1] = fn(val)
			continue
		}

		// 4. Handle Unary Minus (the 'u-' tag from our tokenizer)
		if token == "u-" {
			if len(stack) < 1 {
				return 0, fmt.Errorf("stack empty for u-")
			}
			// Negate the top value on the stack
			stack[len(stack)-1] = -stack[len(stack)-1]
			continue
		}

		// 5. Handle Binary Operators (+, -, *, etc.)
		// These require exactly two values from the stack
		if len(stack) < 2 {
			return 0, fmt.Errorf("missing operands for %s", token)
		}

		// Order matters: 'b' is the top (right operand), 'a' is below it (left operand)
		b := stack[len(stack)-1]
		a := stack[len(stack)-2]
		stack = stack[:len(stack)-2] // Remove both from stack

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
		case "%":
			stack = append(stack, math.Mod(a, b))
		case "^":
			stack = append(stack, math.Pow(a, b))
		default:
			return 0, fmt.Errorf("unknown token: %s", token)
		}
	}

	// After processing all tokens, exactly one value should remain (the result)
	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression: leftover values in stack")
	}

	// Bonus: Store the result in a special "ans" variable for future calculations
	vars["ans"] = stack[0]
	return stack[0], nil
}
