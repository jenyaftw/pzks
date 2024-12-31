package main

import (
	"bufio"
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
	fmt.Println()
	chart.PrintStats(operators, parts)
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	tokensDistClone := make([]token.Token, len(tokens))
	copy(tokensDistClone, tokens)

	tokensAssociativeClone := make([]token.Token, len(tokens))
	copy(tokensAssociativeClone, tokens)

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

	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
	fmt.Println("Running distributive variants")
	fmt.Println()

	for i := 0; i < len(all); i++ {
		fmt.Print("Variant: ")
		for j := 0; j < len(all[i]); j++ {
			print(all[i][j].Text)
		}
		fmt.Println()
		chart, operators, parts := ev.GenerateGanttChart(*tree)
		chart.PrintStats(operators, parts)
		fmt.Println()
	}
	fmt.Println(divider)
	fmt.Println()

	comm := props.NewCommutativeSimplifier()
	perms := comm.Commutate(tokensAssociativeClone)

	fmt.Println("Generated " + fmt.Sprint(len(perms)) + " associative property variants")

	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
	fmt.Println("Running associative variants")
	fmt.Println()
	for i := 0; i < len(perms); i++ {
		fmt.Print("Variant: ")
		for j := 0; j < len(perms[i]); j++ {
			print(perms[i][j].Text)
		}
		fmt.Println()
		chart, operators, parts := ev.GenerateGanttChart(*tree)
		chart.PrintStats(operators, parts)
	}
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	// Print variants, all should have unique height AND width, so make a map to keep track of already found widths and heights

	// fmt.Print("Expanded after distributive property: ")
	// for i := 0; i < len(shortened); i++ {
	// 	print(shortened[i].Text)
	// }
	// println()

	// comm := props.NewCommutativeSimplifier()
	// perms := comm.Commutate(tokens)

	// fmt.Println("Generated " + fmt.Sprint(len(perms)) + " permutations")

	// type Perm struct {
	// 	perm   []token.Token
	// 	width  int
	// 	height int
	// }

	// var permsWithProps []Perm
	// for _, perm := range perms {
	// 	tree := parser.NewParser(perm).Parse()
	// 	height := tree.GetHeight()
	// 	width := tree.GetWidth()
	// 	permsWithProps = append(permsWithProps, Perm{perm: perm, width: width, height: height})
	// }

	// sort.Slice(permsWithProps, func(i, j int) bool {
	// 	if permsWithProps[i].height == permsWithProps[j].height {
	// 		return permsWithProps[i].width < permsWithProps[j].width
	// 	}
	// 	return permsWithProps[i].height < permsWithProps[j].height
	// })

	// // Print top 10 permutations with different height and width
	// width := 0
	// height := 0
	// new := 0
	// for _, perm := range permsWithProps {
	// 	if perm.height != height || perm.width != width {
	// 		width = perm.width
	// 		height = perm.height
	// 		new++

	// 		fmt.Print("Height: ", perm.height, " Width: ", perm.width, " | Permutation: ")
	// 		for _, token := range perm.perm {
	// 			fmt.Print(token.Text)
	// 		}
	// 		fmt.Println()
	// 	}

	// 	if new > 10 {
	// 		break
	// 	}
	// }

}
