package parser

import (
	"fmt"

	"github.com/jenyaftw/lab1/token"
)

type TreeNode struct {
	Id    int
	Token token.Token
	Left  *TreeNode
	Right *TreeNode
}

func (t *TreeNode) PostOrder(indent int) {
	if t != nil {
		if t.Right != nil {
			t.Right.PostOrder(indent + 4)
		}
		if indent > 0 {
			for i := 0; i < indent; i++ {
				print(" ")
			}
		}
		if t.Right != nil {
			print(" /\n")
			for i := 0; i < indent; i++ {
				print(" ")
			}
		}
		fmt.Printf("%d(%s)\n", t.Id, t.Token.Text)
		if t.Left != nil {
			for i := 0; i < indent; i++ {
				print(" ")
			}
			print(" \\\n")
			t.Left.PostOrder(indent + 4)
		}
	}
}

func (t *TreeNode) GetHeight() int {
	if t == nil {
		return 0
	}

	leftHeight := t.Left.GetHeight()
	rightHeight := t.Right.GetHeight()

	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

func (t *TreeNode) GetWidth() int {
	if t == nil {
		return 0
	}

	leftWidth := t.Left.GetWidth()
	rightWidth := t.Right.GetWidth()

	return leftWidth + rightWidth + 1
}

func (t *TreeNode) ReverseLevelOrderTraversal() []*TreeNode {
	queue := make(chan *TreeNode, 300)
	queue <- t

	stack := make([]*TreeNode, 0)
	for len(queue) > 0 {
		node := <-queue
		stack = append(stack, node)

		if node.Right != nil {
			queue <- node.Right
		}
		if node.Left != nil {
			queue <- node.Left
		}
	}

	reverse := make([]*TreeNode, 0)
	for i := len(stack) - 1; i >= 0; i-- {
		reverse = append(reverse, stack[i])
	}
	return reverse
}
