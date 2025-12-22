package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Eval(expression string, vars map[string]float64, functions map[string]func(float64) float64) (float64, error) {

	if strings.TrimSpace(expression) == "" {
		return 0, fmt.Errorf("empty expression")
	}

	tokens, err := infixToPostfix(expression, vars, functions)
	if err != nil {
		return 0, err
	}

	stack := []float64{}
	for _, token := range tokens {
		// 1. Is it a Number?
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
			continue
		}

		// 2. Is it a Variable?
		if val, ok := vars[token]; ok {
			stack = append(stack, val)
			continue
		}

		if fn, ok := functions[token]; ok {
			if len(stack) < 1 {
				return 0, fmt.Errorf("stack empty for function %s", token)
			}
			// Pop the top value, apply function, push result back
			val := stack[len(stack)-1]
			stack[len(stack)-1] = fn(val)
			continue
		}
		// 3. Is it Unary Minus?
		if token == "u-" {
			if len(stack) < 1 {
				return 0, fmt.Errorf("stack empty for u-")
			}
			stack[len(stack)-1] = -stack[len(stack)-1]
			continue
		}

		// 4. Standard Operators
		if len(stack) < 2 {
			return 0, fmt.Errorf("missing operands for %s", token)
		}
		b := stack[len(stack)-1]
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
			stack = append(stack, a/b)
		case "%":
			stack = append(stack, math.Mod(a, b))
		case "^":
			stack = append(stack, math.Pow(a, b))
		default:
			return 0, fmt.Errorf("unknown token: %s", token)
		}
	}

	if len(stack) != 1 {
		print(len(stack))
		return 0, fmt.Errorf("invalid expression leftover values in stack")
	}

	// create or change ans variable storing last result
	vars["ans"] = stack[0]
	return stack[0], nil
}
