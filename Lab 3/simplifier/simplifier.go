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

	return s.OpenParanthases(newTokens)
}

func (s Simplifier) OpenParanthases(tokens []token.Token) []token.Token {
	// ((a + b)) -> (a + b)
	// 1 + (a + b) + 2 -> 1 + a + b + 2

	newTokens := []token.Token{}

	// leftIdx := -1
	// rightIdx := -1
	// count := 0

	for _, v := range tokens {
		newTokens = append(newTokens, v)

		// if v.Type == token.ParanthesesOpenType {
		// 	if leftIdx == -1 {
		// 		leftIdx = len(newTokens) - 1
		// 	}

		// 	count += 1
		// } else if v.Type == token.ParanthesesCloseType {
		// 	count -= 1

		// 	if count == 0 {
		// 		if rightIdx == -1 {
		// 			rightIdx = len(newTokens) - 1
		// 		}
		// 	}
		// }
	}

	// if leftIdx != -1 && rightIdx != -1 {
	// 	if leftIdx-1 < 0 || leftIdx-1 >= 0 && (newTokens[leftIdx-1].Text == "+" || newTokens[leftIdx-1].Text == "-") {
	// 		if rightIdx+1 >= len(newTokens) || rightIdx+1 < len(newTokens) && (newTokens[rightIdx+1].Text == "+" || newTokens[rightIdx+1].Text == "-") {
	// 			left := newTokens[:leftIdx]
	// 			right := newTokens[rightIdx+1:]

	// 			inner := newTokens[leftIdx+1 : rightIdx]

	// 			newTokens = append(append(left, inner...), right...)
	// 		}
	// 	}
	// }

	return newTokens
}
