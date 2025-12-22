package eval

import (
	"fmt"
	"strconv"
)

func infixToPostfix(expression string, vars map[string]float64, functions map[string]func(float64) float64) ([]string, error) {
	tokens, err := tokenizer(expression, vars, functions)

	if err != nil {
		return nil, err
	}
	output := []string{}
	stack := []string{}

	getPrecedence := func(op string) (int, error) {
		switch op {
		case "+", "-":
			return 1, nil
		case "*", "/", "%":
			return 2, nil
		case "u-":
			return 3, nil
		case "^":
			return 4, nil
		}
		if _, ok := functions[op]; ok {
			return 5, nil
		}

		return 0, fmt.Errorf("undefined variable or invalid token %q", op)

	}

	isRightAssoc := func(op string) bool {
		switch op {
		case "^", "u-":
			return true
		}
		if _, ok := functions[op]; ok {
			return true
		}
		return false
	}

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			output = append(output, token)
			continue
		}
		if _, ok := vars[token]; ok {
			output = append(output, token)
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

		precToken, err := getPrecedence(token)
		if err != nil {
			return nil, err
		}
		for len(stack) > 0 {
			top := stack[len(stack)-1]

			if top == "(" {
				break
			}

			precTop, err := getPrecedence(top)
			if err != nil {
				return nil, err
			}
			shouldPop := precTop > precToken || (precTop == precToken && !isRightAssoc(token))
			if shouldPop {
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
