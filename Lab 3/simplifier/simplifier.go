package simplifier

import (
	"github.com/jenyaftw/lab1/token"
)

type Simplifier struct{}

func NewSimplifier() Simplifier {
	return Simplifier{}
}

func (s Simplifier) Simplify(tokens []token.Token) []token.Token {
	newTokens := []token.Token{}
	skipNext := false
	for i, v := range tokens {
		if skipNext {
			skipNext = false
			continue
		}

		if v.Text == "*" || v.Text == "/" {
			if tokens[i+1].Value == 1 {
				skipNext = true
				continue
			}

			if v.Text == "*" && tokens[i-1].Value == 1 {
				skipNext = true
				continue
			}
		}
		newTokens = append(newTokens, v)
	}

	return newTokens
}
