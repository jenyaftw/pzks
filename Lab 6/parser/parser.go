package parser

import (
	"fmt"
	"slices"

	"github.com/jenyaftw/lab1/simplifier"
	"github.com/jenyaftw/lab1/token"
)

type Parser struct {
	tokens []token.Token
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

var globalId = 0

func (p *Parser) toPostfix() (output []token.Token) {
	holding := []token.Token{}

	for _, t := range p.tokens {
		switch t.Type {
		case token.NumericLiteralType:
			output = append(output, t)
		case token.ParanthesesOpenType:
			holding = append(holding, t)
		case token.ParanthesesCloseType:
			for len(holding) > 0 && holding[len(holding)-1].Type != token.ParanthesesOpenType {
				output = append(output, holding[len(holding)-1])
				holding = holding[:len(holding)-1]
			}

			if len(holding) == 0 {
				fmt.Println("Mismatched parantheses")
				return
			}

			if len(holding) > 0 && holding[len(holding)-1].Type == token.ParanthesesOpenType {
				holding = holding[:len(holding)-1]
			}
		case token.TrigFunction:
			holding = append(holding, t)
		case token.OperatorType, token.UnaryOperatorType:
			for len(holding) > 0 && holding[len(holding)-1].Type != token.ParanthesesOpenType {
				front := holding[len(holding)-1]
				if front.Type == token.OperatorType || front.Type == token.UnaryOperatorType {
					if front.GetPrecedence() > t.GetPrecedence() {
						output = append(output, front)
						holding = holding[:len(holding)-1]
					} else {
						break
					}
				}
			}

			holding = append(holding, t)
		}
	}

	for i := len(holding) - 1; i >= 0; i-- {
		output = append(output, holding[i])
		holding = holding[:len(holding)-1]
	}

	return
}

func (p *Parser) ParseInfix() *TreeNode {
	postfix := p.toPostfix()

	// Print
	for _, t := range postfix {
		fmt.Printf("%s ", t.Text)
	}
	fmt.Println()

	stack := []TreeNode{}
	for _, t := range postfix {
		switch t.Type {
		case token.NumericLiteralType:
			stack = append(stack, TreeNode{Token: t})
		case token.OperatorType, token.UnaryOperatorType:
			if len(stack) < 2 {
				fmt.Println("Invalid expression")
				return nil
			}

			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			stack = append(stack, TreeNode{Token: t, Left: &left, Right: &right})
		}
	}

	return &stack[len(stack)-1]
}

func divideTokens(tokens []token.Token, operators []string) (operator *token.Token, left, right []token.Token) {
	opIdxs := []int{}

	paraCount := 0
	for i, t := range tokens {
		if t.Type == token.ParanthesesOpenType {
			paraCount++
		} else if t.Type == token.ParanthesesCloseType {
			paraCount--
		}

		if slices.Contains(operators, t.Text) && paraCount == 0 {
			opIdxs = append(opIdxs, i)
		}
	}

	if len(opIdxs) == 0 {
		return nil, left, right
	}

	centerOp := opIdxs[len(opIdxs)/2]
	operator = &tokens[centerOp]
	left = tokens[:centerOp]
	right = tokens[centerOp+1:]

	return
}

func recursivelyParse(tokens []token.Token) *TreeNode {
	if len(tokens) == 0 {
		return nil
	}

	order := [][]string{
		{"+", "-"},
		{"*", "/"},
		{"^"},
	}

	for _, ops := range order {
		operator, leftTokens, rightTokens := divideTokens(tokens, ops)

		if len(rightTokens) > 2 && rightTokens[0].Type == token.ParanthesesOpenType && rightTokens[len(rightTokens)-1].Type == token.ParanthesesCloseType {
			rightTokens = rightTokens[1 : len(rightTokens)-1]
		}

		if len(leftTokens) > 2 && leftTokens[0].Type == token.ParanthesesOpenType && leftTokens[len(leftTokens)-1].Type == token.ParanthesesCloseType {
			leftTokens = leftTokens[1 : len(leftTokens)-1]
		}

		if operator != nil {
			id := globalId + 1
			globalId += 1

			return &TreeNode{
				Id:    id,
				Token: *operator,
				Left:  recursivelyParse(leftTokens),
				Right: recursivelyParse(rightTokens),
			}
		}
	}

	id := globalId + 1
	globalId += 1
	return &TreeNode{Id: id, Token: tokens[0]}
}

func (p *Parser) Parse() *TreeNode {
	simplifier := simplifier.NewSimplifier()
	p.tokens = simplifier.Simplify(p.tokens)

	return recursivelyParse(p.tokens)
}
