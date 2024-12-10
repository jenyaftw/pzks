package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/jenyaftw/lab1/parser"
	"github.com/jenyaftw/lab1/tokenizer"
)

func main() {
	print("> ")

	reader := bufio.NewReader(os.Stdin)
	expression, _ := reader.ReadString('\n')

	t := tokenizer.NewTokenizer()
	tokens, tokenizerErrors := t.Tokenize(expression)

	print("Вираз: ")
	for i := 0; i < len(expression); i++ {
		for _, e := range tokenizerErrors {
			if i == e.StartIdx {
				print("\033[31m")
			}
		}

		print(string(expression[i]))

		for _, e := range tokenizerErrors {
			if i == e.EndIdx {
				print("\033[0m")
			}
		}

		if i == len(expression)-1 {
			print("\033[0m\n")
		}
	}

	errorsByIdx := make(map[int]string)

	for _, e := range tokenizerErrors {
		errorsByIdx[e.StartIdx] = e.Message
	}

	keys := make([]int, 0)
	for k := range errorsByIdx {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("Позиція %d: %s\n", k, errorsByIdx[k])
	}

	tree := parser.NewParser(tokens).Parse()
	tree.PostOrder(4)
}
