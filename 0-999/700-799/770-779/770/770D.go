package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	children []*Node
	depth    int
	width    int
}

func parse(s string) *Node {
	root := &Node{}
	stack := []*Node{root}
	for _, ch := range s {
		if ch == '[' {
			n := &Node{}
			parent := stack[len(stack)-1]
			parent.children = append(parent.children, n)
			stack = append(stack, n)
		} else if ch == ']' {
			stack = stack[:len(stack)-1]
		}
	}
	return root
}

func computeDepth(n *Node, d int) int {
	n.depth = d
	maxd := d
	for _, c := range n.children {
		if md := computeDepth(c, d+1); md > maxd {
			maxd = md
		}
	}
	return maxd
}

func computeWidth(n *Node) int {
	if len(n.children) == 0 {
		n.width = 5
	} else {
		sum := 0
		for _, c := range n.children {
			sum += computeWidth(c)
		}
		n.width = sum + 2
	}
	return n.width
}

func draw(n *Node, x int, grid [][]byte, H int) {
	if n.depth == 0 {
		cx := x
		for _, c := range n.children {
			draw(c, cx, grid, H)
			cx += c.width
		}
		return
	}
	top := n.depth - 1
	bottom := H - n.depth
	w := n.width
	grid[top][x] = '+'
	grid[top][x+1] = '-'
	grid[top][x+w-2] = '-'
	grid[top][x+w-1] = '+'
	grid[bottom][x] = '+'
	grid[bottom][x+1] = '-'
	grid[bottom][x+w-2] = '-'
	grid[bottom][x+w-1] = '+'
	for r := top + 1; r < bottom; r++ {
		grid[r][x] = '|'
		grid[r][x+w-1] = '|'
	}
	cx := x + 1
	for _, c := range n.children {
		draw(c, cx, grid, H)
		cx += c.width
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	root := parse(s)
	D := computeDepth(root, 0)
	H := 2*D + 1
	computeWidth(root)
	W := 0
	for _, c := range root.children {
		W += c.width
	}
	grid := make([][]byte, H)
	for i := range grid {
		grid[i] = make([]byte, W)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}
	draw(root, 0, grid, H)
	w := bufio.NewWriter(os.Stdout)
	for _, row := range grid {
		fmt.Fprintln(w, string(row))
	}
	w.Flush()
}
