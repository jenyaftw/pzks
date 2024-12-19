package props

import (
	"fmt"

	"github.com/jenyaftw/lab1/simplifier"
	"github.com/jenyaftw/lab1/token"
)

type DistributiveShortener struct{}

func NewDistributiveShortener() DistributiveShortener {
	return DistributiveShortener{}
}

func (d DistributiveShortener) FindExpression(tokens []token.Token) ([]token.Token, []token.Token, *token.Token, int, int) {
	for i, v := range tokens {
		if v.Text == "*" || v.Text == "/" {
			operation := v

			var leftTokens []token.Token
			var rightTokens []token.Token

			leftIdx := 0
			rightIdx := len(tokens) - 1

			prevIdx := i - 1
			nextIdx := i + 1

			foundLeft := false
			foundRight := false

			if len(tokens) > prevIdx {
				if tokens[prevIdx].Type == token.ParanthesesCloseType {
					foundLeft = true

					for j := prevIdx - 1; j >= 0; j-- {
						if tokens[j].Type == token.ParanthesesOpenType {
							leftTokens = tokens[j+1 : prevIdx]
							leftIdx = j
							break
						}
					}
				} else if tokens[prevIdx].Type == token.NumericLiteralType || tokens[prevIdx].Type == token.Symbol {
					leftTokens = []token.Token{
						tokens[prevIdx],
					}

					if prevIdx-1 >= 0 && tokens[prevIdx-1].Type == token.UnaryOperatorType {
						leftTokens = []token.Token{
							tokens[prevIdx-1],
							tokens[prevIdx],
						}
					}

					leftIdx = prevIdx
				}
			}

			if len(tokens) > nextIdx {
				if tokens[nextIdx].Type == token.ParanthesesOpenType {
					foundRight = true

					for j := nextIdx + 1; j < len(tokens); j++ {
						if tokens[j].Type == token.ParanthesesCloseType {
							rightTokens = tokens[nextIdx+1 : j]
							rightIdx = j
							break
						}
					}
				} else if tokens[nextIdx].Type == token.NumericLiteralType || tokens[nextIdx].Type == token.Symbol {
					rightTokens = []token.Token{
						tokens[nextIdx],
					}

					rightIdx = nextIdx
				}
			}

			if !foundLeft && !foundRight {
				continue
			}

			return leftTokens, rightTokens, &operation, leftIdx, rightIdx
		}
	}

	return nil, nil, nil, 0, 0
}

func (d DistributiveShortener) OpenExpression(leftTokens, rightTokens []token.Token, operation token.Token) []token.Token {
	result := []token.Token{}
	leftExpr := []token.Token{}
	rightExpr := []token.Token{}

	addToResult := func() {
		result = append(result, leftExpr...)
		result = append(result, operation)
		result = append(result, rightExpr...)
	}

	for i, v := range leftTokens {
		doStuff := func() {
			for j, k := range rightTokens {
				if k.Type == token.UnaryOperatorType || (k.Text != "-" && k.Text != "+") {
					rightExpr = append(rightExpr, k)

					if j+1 >= len(rightTokens) {
						addToResult()
						rightExpr = []token.Token{}
					}
				} else {
					addToResult()
					result = append(result, k)
					rightExpr = []token.Token{}
				}
			}

			leftExpr = []token.Token{}
		}

		if v.Type == token.UnaryOperatorType || (v.Text != "-" && v.Text != "+") {
			fmt.Println(v.String())
			leftExpr = append(leftExpr, v)

			if i+1 >= len(leftTokens) {
				doStuff()
				continue
			}
		} else {
			doStuff()

			result = append(result, v)
		}
	}

	return result
}

func (d DistributiveShortener) WrapInParanthases(tokens []token.Token) []token.Token {
	var newTokens []token.Token

	newTokens = append(newTokens, token.Token{Type: token.ParanthesesOpenType, Text: "("})

	for _, v := range tokens {
		newTokens = append(newTokens, v)
	}

	return append(newTokens, token.Token{Type: token.ParanthesesCloseType, Text: ")"})
}

func (d DistributiveShortener) Shorten(tokens []token.Token) []token.Token {
	leftTokens, rightTokens, operation, leftIdx, rightIdx := d.FindExpression(tokens)
	simplifier := simplifier.NewSimplifier()

	if len(leftTokens) == 0 && len(rightTokens) == 0 {
		return simplifier.Simplify(tokens)
	}

	for i := 0; i < len(leftTokens); i++ {
		print(leftTokens[i].Text)
	}
	println()

	for i := 0; i < len(rightTokens); i++ {
		print(rightTokens[i].Text)
	}
	println()

	opened := d.OpenExpression(leftTokens, rightTokens, *operation)
	simplified := simplifier.Simplify(opened)
	wrapped := d.WrapInParanthases(simplified)

	newTokens := []token.Token{}
	newTokens = append(newTokens, tokens[0:leftIdx]...)
	newTokens = append(newTokens, wrapped...)
	newTokens = append(newTokens, tokens[rightIdx+1:]...)

	return simplifier.Simplify(d.Shorten(newTokens))
}
