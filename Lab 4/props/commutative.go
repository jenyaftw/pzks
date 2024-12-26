package props

import (
	"fmt"

	"github.com/jenyaftw/lab1/simplifier"
	"github.com/jenyaftw/lab1/token"
	"github.com/jenyaftw/lab1/tokenizer"
)

type CommutativeSimplifier struct{}

func NewCommutativeSimplifier() CommutativeSimplifier {
	return CommutativeSimplifier{}
}

func (c CommutativeSimplifier) Permutate(arr []string, start int) [][]string {
	var results [][]string

	var generate func(int)
	generate = func(start int) {
		if start == len(arr)-1 {
			perm := make([]string, len(arr))
			copy(perm, arr)
			results = append(results, perm)
			return
		}

		for i := start; i < len(arr); i++ {
			arr[start], arr[i] = arr[i], arr[start]
			generate(start + 1)
			arr[start], arr[i] = arr[i], arr[start]
		}
	}
	generate(start)

	return results
}

func (c CommutativeSimplifier) Commutate(tokens []token.Token) [][]token.Token {
	str := ""
	for _, v := range tokens {
		str += v.Text
	}

	simplifier := simplifier.NewSimplifier()
	matches := simplifier.SplitBySigns(str)
	fmt.Println("Matches:", len(matches))
	for _, match := range matches {
		fmt.Println(match)
	}

	perms := c.Permutate(matches, 0)
	newPerms := [][]token.Token{}
	tokenizer := tokenizer.NewTokenizer()
	for _, perm := range perms {
		newPerm := ""
		for i, part := range perm {
			if i == 0 && part[0] == '+' {
				part = part[1:]
			} else if i > 0 && (part[0] != '+' && part[0] != '-') {
				part = "+" + part
			}

			newPerm += part
		}

		tokens, errors := tokenizer.Tokenize(newPerm)
		if len(errors) > 0 {
			continue
		}
		newPerms = append(newPerms, tokens)
	}

	return newPerms
}
