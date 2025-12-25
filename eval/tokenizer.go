package eval

import (
	"fmt"
	"strings"
	"unicode"
)

func tokenizer(expression string, vars map[string]float64, functions map[string]func(float64) float64) ([]string, error) {
	
	// Helper to identify supported mathematical operators and delimiters
	isOperator := func(s string) bool {
		return strings.Contains("+-*/^()%", s)
	}

	tokens := []string{}
	var numBuffer string // Accumulates digits for multi-digit/decimal numbers
	var varBuffer string // Accumulates characters for variable/function names

	// addToken handles the logic of placing a new token into the slice,
	// checking for special cases like unary operators and implicit multiplication.
	addToken := func(newToken string) {
		if len(tokens) > 0 {
			lastToken := tokens[len(tokens)-1]

			// Context flags for the previous and current tokens
			isLastClosing := lastToken == ")"
			isCurrentOpening := newToken == "("
			isLastNumber := unicode.IsDigit(rune(lastToken[0])) || lastToken[0] == '.'
			isCurrentNumber := unicode.IsDigit(rune(newToken[0])) || newToken[0] == '.'

			_, isLastVar := vars[lastToken]
			_, isCurrentVar := vars[newToken]
			_, isCurrentFunc := functions[newToken]

			// Handle Unary Operators: 
			// If '+' or '-' appears at the start or after another operator (except ')'),
			// it is treated as a unary sign. We label negative as 'u-'.
			if (newToken == "-" || newToken == "+") && !isLastNumber && !isLastClosing && !isCurrentVar {
				if newToken == "-" {
					tokens = append(tokens, "u-")
				}
				// Note: Unary '+' is ignored as it doesn't change the value
				return
			}

			// Handle Implicit Multiplication:
			// Automatically inserts '*' in cases like '2(x)', '(x)y', or '2x'
			shouldMultiply := (isLastClosing && (isCurrentOpening || isCurrentNumber || isCurrentVar)) ||
				(isLastNumber && (isCurrentOpening || isCurrentVar || isCurrentFunc)) ||
				(isLastVar && (isCurrentOpening || isCurrentVar || isCurrentFunc))
			
			if shouldMultiply {
				tokens = append(tokens, "*")
			}
		} else if newToken == "-" {
			// Special case: Leading negative sign is always unary
			tokens = append(tokens, "u-")
			return
		}
		
		tokens = append(tokens, newToken)
	}

	// Main Loop: Iterate through every character in the input string
	for _, char := range expression {
		if unicode.IsSpace(char) {
			// Space marks the end of a sequence; flush buffers
			if numBuffer != "" {
				addToken(numBuffer)
				numBuffer = ""
			}
			if varBuffer != "" {
				addToken(varBuffer)
				varBuffer = ""
			}
		} else if unicode.IsDigit(char) || char == '.' {
			// Collect numeric characters
			numBuffer += string(char)
		} else if unicode.IsLetter(char) {
			// Collect alphabetic characters (variables or function names)
			varBuffer += string(char)
		} else if isOperator(string(char)) {
			// If we hit an operator, first flush any pending buffers
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
			// Return error for symbols not explicitly handled
			return nil, fmt.Errorf("invalid character: %c", char)
		}
	}

	// Final flush: process any remaining content in buffers after loop ends
	if numBuffer != "" {
		addToken(numBuffer)
	}
	if varBuffer != "" {
		addToken(varBuffer)
	}

	return tokens, nil
}
