package parser

import "github.com/jenyaftw/lab1/token"

type Parser struct {
	tokens []token.Token
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() *TreeNode {
	if len(p.tokens) == 0 {
		return nil
	}

	var stack []*TreeNode

	for _, t := range p.tokens {
		if t.Type == token.OperatorType {
			stack = append(stack, &TreeNode{Token: t})
		} else {
			if len(stack) < 2 {
				return nil
			}

			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			stack = append(stack, &TreeNode{Token: t, Left: left, Right: right})
		}
	}

	if len(stack) != 1 {
		return nil
	}

	return stack[0]
}
