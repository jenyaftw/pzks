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
	findLeftRight := func(i int) ([]token.Token, []token.Token, int, int, bool, bool) {
		var leftTokens []token.Token
		var rightTokens []token.Token

		leftIdx := 0
		rightIdx := len(tokens) - 1

		prevIdx := i - 1
		nextIdx := i + 1

		foundLeft := false
		foundRight := false

		if len(tokens) > prevIdx {
			count := 0
			if tokens[prevIdx].Type == token.ParanthesesCloseType {
				count += 1
				foundLeft = true

				for j := prevIdx - 1; j >= 0; j-- {
					if tokens[j].Type == token.ParanthesesOpenType {
						count -= 1

						if count == 0 {
							leftTokens = tokens[j+1 : prevIdx]
							leftIdx = j
							break
						}
					} else if tokens[j].Type == token.ParanthesesCloseType {
						count += 1
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
			count := 0
			if tokens[nextIdx].Type == token.ParanthesesOpenType {
				count += 1
				foundRight = true

				for j := nextIdx + 1; j < len(tokens); j++ {
					if tokens[j].Type == token.ParanthesesCloseType {
						count -= 1

						if count == 0 {
							rightTokens = tokens[nextIdx+1 : j]
							rightIdx = j
							break
						}
					} else if tokens[j].Type == token.ParanthesesOpenType {
						count += 1
					}
				}
			} else if tokens[nextIdx].Type == token.NumericLiteralType || tokens[nextIdx].Type == token.Symbol {
				rightTokens = []token.Token{
					tokens[nextIdx],
				}

				rightIdx = nextIdx
			}
		}

		return leftTokens, rightTokens, leftIdx, rightIdx, foundLeft, foundRight
	}

	lastDivide := -1
	numDivides := 0

	count := 0
	for i, v := range tokens {
		if v.Type == token.ParanthesesOpenType {
			count += 1
		} else if v.Type == token.ParanthesesCloseType {
			count -= 1
		}

		if v.Text == "*" {
			operation := v

			leftTokens, rightTokens, leftIdx, rightIdx, foundLeft, foundRight := findLeftRight(i)

			if !foundLeft && !foundRight {
				continue
			}

			return leftTokens, rightTokens, &operation, leftIdx, rightIdx
		} else if v.Text == "/" {
			lastDivide = i
			numDivides += 1
		} else if count == 0 && (v.Text == "+" || v.Text == "-") {
			if numDivides <= 1 {
				numDivides = 0
			} else {
				break
			}
		}
	}

	if numDivides >= 2 {
		operation := tokens[lastDivide]

		leftTokens, rightTokens, leftIdx, rightIdx, _, _ := findLeftRight(lastDivide)

		return leftTokens, rightTokens, &operation, leftIdx, rightIdx
	}

	return nil, nil, nil, 0, 0
}

func (d DistributiveShortener) OpenExpression(leftTokens, rightTokens []token.Token, operation token.Token) []token.Token {
	result := []token.Token{}
	leftExpr := []token.Token{}
	rightExpr := []token.Token{}

	if operation.Text == "/" {
		wrappedLeft := d.WrapInParanthases(leftTokens)
		wrappedRight := d.WrapInParanthases(rightTokens)
		result = append(result, wrappedLeft...)
		result = append(result, token.Token{Type: token.OperatorType, Text: "*"})
		result = append(result, wrappedRight...)
		return result
	}

	var curOperator *token.Token

	addToResult := func() {
		result = append(result, leftExpr...)
		result = append(result, operation)
		result = append(result, rightExpr...)
	}

	for i, v := range leftTokens {
		doStuff := func() {
			count := 0
			for j, k := range rightTokens {
				if k.Type == token.ParanthesesOpenType {
					count += 1
				} else if k.Type == token.ParanthesesCloseType {
					count -= 1
				}

				if k.Type == token.UnaryOperatorType || (k.Text != "-" && k.Text != "+") {
					rightExpr = append(rightExpr, k)

					if j+1 >= len(rightTokens) {
						addToResult()
						rightExpr = []token.Token{}
					}
				} else {
					if j+1 < len(rightTokens) && (rightTokens[j+1].Text == "*" || rightTokens[j+1].Text == "/") {
						rightExpr = append(rightExpr, k)
						continue
					}

					if count == 0 {
						addToResult()
					} else {
						rightExpr = append(rightExpr, k)
						continue
					}

					if curOperator == nil {
						result = append(result, k)
					} else if curOperator.Text == "+" && k.Text == "+" {
						result = append(result, k)
					} else if curOperator.Text == "+" && k.Text == "-" {
						result = append(result, token.Token{Type: token.OperatorType, Text: "-"})
					} else if curOperator.Text == "-" && k.Text == "+" {
						result = append(result, token.Token{Type: token.OperatorType, Text: "-"})
					} else if curOperator.Text == "-" && k.Text == "-" {
						result = append(result, token.Token{Type: token.OperatorType, Text: "+"})
					} else {
						fmt.Print("Unexpected operator: ", v.Text, curOperator.Text, k.Text, " | ")
						fmt.Print("Left: ")
						for _, v := range leftTokens {
							fmt.Print(v.Text)
						}
						fmt.Print(" | Right: ")
						for _, v := range rightTokens {
							fmt.Print(v.Text)
						}
						fmt.Println()

						result = append(result, k)
					}

					rightExpr = []token.Token{}
				}
			}

			leftExpr = []token.Token{}
		}

		if v.Type == token.UnaryOperatorType || (v.Text != "-" && v.Text != "+") {
			leftExpr = append(leftExpr, v)

			if i-1 >= 0 && leftTokens[i-1].Type == token.OperatorType {
				if leftTokens[i-1].Text == "+" || leftTokens[i-1].Text == "-" {
					curOperator = &leftTokens[i-1]
				}
			}

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

var allVariants [][]token.Token

func (d DistributiveShortener) Shorten(tokens []token.Token) ([]token.Token, [][]token.Token) {
	leftTokens, rightTokens, operation, leftIdx, rightIdx := d.FindExpression(tokens)
	simplifier := simplifier.NewSimplifier()

	if len(leftTokens) == 0 && len(rightTokens) == 0 {
		return simplifier.Simplify(tokens), allVariants
	}

	opened := d.OpenExpression(leftTokens, rightTokens, *operation)
	simplified := simplifier.Simplify(opened)
	wrapped := d.WrapInParanthases(simplified)

	newTokens := []token.Token{}
	newTokens = append(newTokens, tokens[0:leftIdx]...)
	newTokens = append(newTokens, wrapped...)
	newTokens = append(newTokens, tokens[rightIdx+1:]...)

	allVariants = append(allVariants, simplifier.Simplify(newTokens))
	shortened, _ := d.Shorten(newTokens)
	final := simplifier.Simplify(shortened)
	return final, allVariants
}
