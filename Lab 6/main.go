package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	"github.com/jenyaftw/lab1/evaluator"
	"github.com/jenyaftw/lab1/parser"
	"github.com/jenyaftw/lab1/props"
	"github.com/jenyaftw/lab1/simplifier"
	"github.com/jenyaftw/lab1/token"
	"github.com/jenyaftw/lab1/tokenizer"
)

func main() {
	print("> ")

	simplifier := simplifier.NewSimplifier()
	reader := bufio.NewReader(os.Stdin)
	expression, _ := reader.ReadString('\n')
	expression = simplifier.RemoveWhitespaces(expression)

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

	divider := ""
	dividerLength := 75
	for i := 0; i < dividerLength; i++ {
		divider += "-"
	}

	tree := parser.NewParser(tokens).Parse()
	fmt.Println("Generated tree with height: ", tree.GetHeight(), " and width: ", tree.GetWidth())
	tree.PostOrder(4)

	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
	fmt.Println("Running base variant")
	fmt.Println()
	ev := evaluator.NewEvaluator()
	chart, operators, parts := ev.GenerateGanttChart(*tree)
	chart.Print(operators)
	chart.PrintStats(operators, parts)
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	records := [][]string{
		{"t", "input", "s1", "s2", "s3", "s4", "s5", "s6", "output"},
	}

	f, err := os.OpenFile("yeet.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for i := 0; i < len(chart.S[len(chart.S)-1]); i++ {
		strings := []string{}
		strings = append(strings, fmt.Sprint(i+1))
		strings = append(strings, fmt.Sprint(""))

		for j := 0; j < len(chart.S); j++ {
			if i < len(chart.S[j]) {
				text := operators[chart.S[j][i]].Token.Text
				if operators[chart.S[j][i]].Token.Type == token.OperatorType {
					strings = append(strings, fmt.Sprint(text))
				} else {
					strings = append(strings, "")
				}
			} else {
				strings = append(strings, "")
			}
		}
		strings = append(strings, fmt.Sprint(""))

		records = append(records, strings)
	}
	w.WriteAll(records)
	if err := w.Error(); err != nil {
		fmt.Println("error writing csv:", err)
	}

	tokensDistClone := make([]token.Token, len(tokens))
	copy(tokensDistClone, tokens)

	tokensCommutativeClone := make([]token.Token, len(tokens))
	copy(tokensCommutativeClone, tokens)

	dist := props.NewDistributiveShortener()
	shortened, all := dist.Shorten(tokensDistClone)

	fmt.Print("Shortened with distributive property: ")
	for i := 0; i < len(shortened); i++ {
		print(shortened[i].Text)
	}
	fmt.Println()

	fmt.Println("Generated " + fmt.Sprint(len(all)) + " distributive property variants")
	for i := 0; i < len(all); i++ {
		fmt.Print("Variant ", i+1, ": ")
		for j := 0; j < len(all[i]); j++ {
			print(all[i][j].Text)
		}
		fmt.Println()
	}

	// fmt.Println()
	// fmt.Println(divider)
	// fmt.Println()
	// fmt.Println("Running distributive variants")
	// fmt.Println()

	// for i := 0; i < len(all); i++ {
	// 	fmt.Print("Variant: ")
	// 	for j := 0; j < len(all[i]); j++ {
	// 		print(all[i][j].Text)
	// 	}
	// 	fmt.Println()
	// 	tree := parser.NewParser(all[i]).Parse()
	// 	chart, operators, parts := ev.GenerateGanttChart(*tree)
	// 	chart.PrintStats(operators, parts)
	// 	fmt.Println()
	// }
	// fmt.Println(divider)
	// fmt.Println()

	comm := props.NewCommutativeSimplifier()
	perms := comm.Commutate(tokensCommutativeClone)

	fmt.Println("Generated " + fmt.Sprint(len(perms)) + " commutative property variants")
	for k, v := range perms {
		tree := parser.NewParser(v).Parse()
		height := tree.GetHeight()
		width := tree.GetWidth()

		fmt.Print("Variant ", k+1, ": ")
		for j := 0; j < len(v); j++ {
			print(v[j].Text)
		}
		fmt.Print(" | Height: ")
		fmt.Print(height)
		fmt.Print(" | Width: ")
		fmt.Print(width)
		fmt.Println()
	}

	// fmt.Println()
	// fmt.Println(divider)
	// fmt.Println()
	// fmt.Println("Running associative variants")
	// fmt.Println()
	// for i := 0; i < len(perms); i++ {
	// 	fmt.Print("Variant: ")
	// 	for j := 0; j < len(perms[i]); j++ {
	// 		print(perms[i][j].Text)
	// 	}
	// 	fmt.Println()
	// 	tree := parser.NewParser(perms[i]).Parse()
	// 	chart, operators, parts := ev.GenerateGanttChart(*tree)
	// 	chart.PrintStats(operators, parts)
	// 	fmt.Println()
	// }
	// fmt.Println(divider)
	// fmt.Println()
}
