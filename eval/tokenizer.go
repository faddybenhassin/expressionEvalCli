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
	var varBuffer string

	// Helper to safely append tokens and check for implicit multiplication
	addToken := func(newToken string) {
		if len(tokens) > 0 {
			lastToken := tokens[len(tokens)-1]

			isLastClosing := lastToken == ")"
			isCurrentOpening := newToken == "("
			isLastNumber := unicode.IsDigit(rune(lastToken[0])) || lastToken[0] == '.'
			isCurrentNumber := unicode.IsDigit(rune(newToken[0])) || newToken[0] == '.'

			// Check for Name: Starts with letter, but IS NOT our internal "u-"
			isLastName := unicode.IsLetter(rune(lastToken[0])) && lastToken != "u-"
			isCurrentName := unicode.IsLetter(rune(newToken[0])) && newToken != "u-"
			if (newToken == "-" || newToken == "+") && !isLastNumber && !isLastClosing && !isCurrentName {
				if newToken == "-" {
					tokens = append(tokens, "u-")
				}
				return
			}
			shouldMultiply := (isLastClosing && (isCurrentOpening || isCurrentNumber || isCurrentName)) ||
				(isLastNumber && (isCurrentOpening || isCurrentName)) ||
				(isLastName && (isCurrentOpening || isCurrentName))
			if shouldMultiply {
				tokens = append(tokens, "*")
			}
		} else if newToken == "-" {
			tokens = append(tokens, "u-")
			return
		}
		tokens = append(tokens, newToken)
	}

	for _, char := range expression {
		if unicode.IsSpace(char) {
			if numBuffer != "" {
				addToken(numBuffer)
				numBuffer = ""
			}
			if varBuffer != "" {
				addToken(varBuffer)
				varBuffer = ""
			}
		} else if unicode.IsDigit(char) || char == '.' {
			numBuffer += string(char)
		} else if unicode.IsLetter(char) {
			varBuffer += string(char)
		} else if isOperator(string(char)) {
			if numBuffer != "" {
				addToken(numBuffer)
				numBuffer = ""
			}
			if varBuffer != "" {
				addToken(varBuffer)
				varBuffer = ""
			}
			addToken(string(char))
		} else {
			return nil, fmt.Errorf("invalid character: %c", char)
		}
	}

	if numBuffer != "" {
		addToken(numBuffer)
	}
	if varBuffer != "" {
		addToken(varBuffer)
	}

	return tokens, nil
}
