package eval

import (
	"fmt"
	"strings"
	"unicode"
)

func tokenizer(expression string) ([]string, error) {
	isOperator := func(s string) bool {
		return strings.Contains("+-*/^()", s)
	}

	tokens := []string{}
	var numBuffer string

	// Helper to safely append tokens and check for implicit multiplication
	addToken := func(newToken string) {
		if len(tokens) > 0 {
			lastToken := tokens[len(tokens)-1]

			// Logic: Insert '*' if:
			// 1. Last was ')' and current is a Number/Variable: (2)3 -> (2)*3
			// 2. Last was ')' and current is '(': (2)(3) -> (2)*(3)
			// 3. Last was Number/Variable and current is '(': 2(3) -> 2*(3)

			isLastClosing := lastToken == ")"
			isCurrentOpening := newToken == "("
			isLastNumber := unicode.IsDigit(rune(lastToken[0]))
			isCurrentNumber := unicode.IsDigit(rune(newToken[0]))
			if newToken == "-" && isOperator(lastToken) && !isLastClosing {
				tokens = append(tokens, "u-")
				return
			}
			if (isLastClosing && (isCurrentOpening || isCurrentNumber)) ||
				(isLastNumber && isCurrentOpening) {
				tokens = append(tokens, "*")
			}
		}
		tokens = append(tokens, newToken)
	}

	for _, char := range strings.ReplaceAll(expression, " ", "") {
		if unicode.IsDigit(char) || char == '.' {
			numBuffer += string(char)
		} else if isOperator(string(char)) {
			if numBuffer != "" {
				addToken(numBuffer)
				numBuffer = ""
			}
			addToken(string(char))
		} else {
			return nil, fmt.Errorf("invalid character: %c", char)
		}
	}

	if numBuffer != "" {
		addToken(numBuffer)
	}

	return tokens, nil
}
