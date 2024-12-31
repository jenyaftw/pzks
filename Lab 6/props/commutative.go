package props

import (
	"slices"

	prmt "github.com/gitchander/permutation"
	"github.com/jenyaftw/lab1/parser"
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

// func (c CommutativeSimplifier) Commutate(tokens []token.Token) [][]token.Token {
// 	str := ""
// 	for _, v := range tokens {
// 		str += v.Text
// 	}

// 	simplifier := simplifier.NewSimplifier()
// 	matches := simplifier.SplitBySigns(str)

// 	perms := c.Permutate(matches, 0)
// 	newPerms := [][]token.Token{}
// 	tokenizer := tokenizer.NewTokenizer()
// 	for _, perm := range perms {
// 		newPerm := ""
// 		for i, part := range perm {
// 			if i == 0 && part[0] == '+' {
// 				part = part[1:]
// 			} else if i > 0 && (part[0] != '+' && part[0] != '-') {
// 				part = "+" + part
// 			}

// 			newPerm += part
// 		}

// 		tokens, errors := tokenizer.Tokenize(newPerm)
// 		if len(errors) > 0 {
// 			continue
// 		}
// 		newPerms = append(newPerms, tokens)
// 	}

// 	return newPerms
// }

func (c CommutativeSimplifier) Commutate(tokens []token.Token) [][]token.Token {
	str := ""
	for _, v := range tokens {
		str += v.Text
	}

	simplifier := simplifier.NewSimplifier()
	matches := simplifier.SplitBySigns(str)

	p := prmt.New(prmt.StringSlice(matches))
	// perms := c.Permutate(matches, 0)
	heightWidth := map[int][]int{}
	newPerms := [][]token.Token{}
	tokenizer := tokenizer.NewTokenizer()
	for p.Next() {
		newPerm := ""
		for i, part := range matches {
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

		tree := parser.NewParser(tokens).Parse()
		height := tree.GetHeight()
		width := tree.GetWidth()

		if _, ok := heightWidth[height]; !ok {
			heightWidth[height] = []int{}
		}

		if slices.Contains(heightWidth[height], width) {
			continue
		}
		heightWidth[height] = append(heightWidth[height], width)
		newPerms = append(newPerms, tokens)
	}

	return newPerms
}
