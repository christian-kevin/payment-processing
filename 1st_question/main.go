package main

import (
	"fmt"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func searchBST(root *TreeNode, val int) bool {
	if root == nil {
		return false
	}

	if root.Val == val {
		return true
	}

	if root.Left == nil && root.Right == nil {
		return false
	}

	if searchBST(root.Left, val) || searchBST(root.Right, val) {
		return true
	}

	return false
}

func buildDummyTree() (root *TreeNode) {
	r := TreeNode{
		Val:   7,
	}


	r.Left = createTreeNode(6)
	r.Left.Left = createTreeNode(3)
	r.Left.Right = createTreeNode(20)

	r.Right = createTreeNode(100)
	r.Right.Right = createTreeNode(140)
	return &r
}

func createTreeNode(val int) *TreeNode {
	return &TreeNode{Val: val}
}

func main() {
	root := buildDummyTree()

	fmt.Println(searchBST(root, 6))
	fmt.Println(searchBST(root, 7))
	fmt.Println(searchBST(root, 8))
	fmt.Println(searchBST(root, 9))
	fmt.Println(searchBST(root, 20))
	fmt.Println(searchBST(root, 140))
}
