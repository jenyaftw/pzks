package simplifier

import (
	"github.com/jenyaftw/lab1/token"
)

type Simplifier struct{}

func NewSimplifier() Simplifier {
	return Simplifier{}
}

func (s Simplifier) RemoveWhitespaces(from string) string {
	newStr := ""
	for _, v := range from {
		if v != ' ' {
			newStr += string(v)
		}
	}
	return newStr
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
	removeOpen := -1
	removeClose := -1

	for i, v := range tokens {
		if v.Text == "(" {
			if i-1 < 0 || i-1 >= 0 && (tokens[i-1].Text == "+" || tokens[i-1].Text == "-" || tokens[i-1].Text == "(") {
				count := 1

				for j := i + 1; j < len(tokens); j++ {
					if tokens[j].Text == "(" {
						count++
					}

					if tokens[j].Text == ")" {
						count--
					}

					if count == 0 {
						if j+1 >= len(tokens) || j+1 < len(tokens) && (tokens[j+1].Text == "+" || tokens[j+1].Text == "-" || tokens[j+1].Text == ")") {
							removeOpen = i
							removeClose = j
						}
						break
					}
				}
			}

			if removeOpen != -1 {
				break
			}
		}
	}

	if removeOpen == -1 || removeClose == -1 {
		return tokens
	}

	newTokens := []token.Token{}
	negative := false

	if removeOpen-1 >= 0 && tokens[removeOpen-1].Text == "-" {
		negative = true
	}

	count := 0
	neededCount := -1
	switchSign := false
	for i, v := range tokens {
		if v.Text == "(" {
			count++
		} else if v.Text == ")" {
			count--
		}

		if count == neededCount && switchSign {
			if v.Text == "+" {
				v.Text = "-"
			} else if v.Text == "-" {
				v.Text = "+"
			}
		}

		if i == removeOpen {
			if negative {
				switchSign = true
				neededCount = count
			}
			continue
		} else if i == removeClose {
			switchSign = false
			continue
		}

		newTokens = append(newTokens, v)
	}

	return s.OpenParanthases(newTokens)
}

func (s Simplifier) SplitBySigns(expr string) []string {
	splits := []string{}

	bracketCount := 0
	currentSplit := ""
	for _, c := range expr {
		if c == '(' {
			bracketCount++
		} else if c == ')' {
			bracketCount--
		}

		if bracketCount == 0 && (c == '+' || c == '-') {
			if len(currentSplit) > 0 {
				splits = append(splits, currentSplit)
			}

			currentSplit = string(c)
		} else {
			currentSplit += string(c)
		}
	}

	if len(currentSplit) > 0 {
		splits = append(splits, currentSplit)
	}

	return splits
}
