package eval

import (
	"fmt"
	"strconv"
)

func infixToPostfix(expression string, vars map[string]float64, functions map[string]func(float64) float64) ([]string, error) {
	// 1. Convert the raw string into a slice of discrete tokens
	tokens, err := tokenizer(expression, vars, functions)
	if err != nil {
		return nil, err
	}

	output := []string{} // The "Output Queue" (Postfix result)
	stack := []string{}  // The "Operator Stack" (Temporary holding for operators)

	// Defines the hierarchy of operations (PEMDAS/BODMAS logic)
	getPrecedence := func(op string) (int, error) {
		switch op {
		case "+", "-":
			return 1, nil
		case "*", "/", "%":
			return 2, nil
		case "u-": // Unary minus (negative sign) has high priority
			return 3, nil
		case "^":
			return 4, nil
		}
		// Functions (like sin, cos) have the highest precedence
		if _, ok := functions[op]; ok {
			return 5, nil
		}
		return 0, fmt.Errorf("undefined variable or invalid token %q", op)
	}

	// Right-associative operators (like exponents and unary minus) 
	// are evaluated from right to left (e.g., 2^3^2 is 2^(3^2))
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
		// A. If token is a number, add it directly to the output
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			output = append(output, token)
			continue
		}
		// B. If token is a known variable, add it to the output
		if _, ok := vars[token]; ok {
			output = append(output, token)
			continue
		}

		// C. Handle Opening Parenthesis: push to stack
		if token == "(" {
			stack = append(stack, token)
			continue
		}

		// D. Handle Closing Parenthesis: pop everything until we hit "("
		if token == ")" {
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == "(" {
					found = true
					stack = stack[:len(stack)-1] // Discard the "("
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

		// E. Handle Operators and Functions
		precToken, err := getPrecedence(token)
		if err != nil {
			return nil, err
		}
		
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top == "(" {
				break
			}

			precTop, _ := getPrecedence(top)
			
			// Decide if the operator on the stack should be processed first:
			// 1. Stack operator has higher precedence
			// 2. Precedences are equal AND the current token is Left-Associative
			shouldPop := precTop > precToken || (precTop == precToken && !isRightAssoc(token))
			
			if shouldPop {
				output = append(output, top)
				stack = stack[:len(stack)-1]
			} else {
				break
			}
		}
		// Push the current operator onto the stack
		stack = append(stack, token)
	}

	// 2. Final Step: Pop all remaining operators from the stack to the output
	for i := len(stack) - 1; i >= 0; i-- {
		if stack[i] == "(" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[i])
	}

	return output, nil
}
