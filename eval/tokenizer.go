package eval

import (
	"fmt"
	"strings"
	"unicode"
)

func tokenizer(expression string) ([]string, error) {

	isOperator := func(char rune) bool {
		return char == '+' || char == '-' || char == '*' || char == '/' || char == '^' || char == '(' || char == ')'
	}

	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, fmt.Errorf("empty expression")
	}

	tokens := []string{}
	var numBuffer string
	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			numBuffer += string(char)
			continue
		}
		isOp := isOperator(char)
		if !isOp {
			return nil, fmt.Errorf("invalid token %c", char)
		}
		tokens = append(tokens, numBuffer)
		tokens = append(tokens, string(char))
		numBuffer = ""
	}
	return nil, fmt.Errorf("empty expression")
}
