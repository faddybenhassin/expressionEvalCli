package eval

import (
	"fmt"
	"strconv"
)

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
		case "u-":
			return 3, nil
		case "^":
			return 4, nil
		default:
			return 0, fmt.Errorf("unknown operator %q", op)
		}
	}

	isRightAssoc := func(op string) bool {
		return op == "^" || op == "u-"
	}

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
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
