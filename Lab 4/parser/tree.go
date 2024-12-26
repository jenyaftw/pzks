package parser

import "github.com/jenyaftw/lab1/token"

type TreeNode struct {
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
		print(t.Token.Text, "\n")
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
