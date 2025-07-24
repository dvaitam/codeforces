package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	sum int
	min int
	max int
}

var tree []Node
var size int

func merge(a, b Node) Node {
	res := Node{}
	res.sum = a.sum + b.sum
	res.min = a.min
	if a.sum+b.min < res.min {
		res.min = a.sum + b.min
	}
	res.max = a.max
	if a.sum+b.max > res.max {
		res.max = a.sum + b.max
	}
	return res
}

func update(v, tl, tr, pos, val int) {
	if tl == tr {
		tree[v].sum = val
		if val < 0 {
			tree[v].min = val
		} else {
			tree[v].min = 0
		}
		if val > 0 {
			tree[v].max = val
		} else {
			tree[v].max = 0
		}
		return
	}
	tm := (tl + tr) / 2
	if pos <= tm {
		update(v*2, tl, tm, pos, val)
	} else {
		update(v*2+1, tm+1, tr, pos, val)
	}
	tree[v] = merge(tree[v*2], tree[v*2+1])
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)

	size = n + 2
	tree = make([]Node, 4*size)
	arr := make([]int8, size+1)
	pos := 1

	maxPos := 1
	for i := 0; i < n; i++ {
		c := s[i]
		if c == 'L' {
			if pos > 1 {
				pos--
			}
		} else if c == 'R' {
			pos++
			if pos > maxPos {
				maxPos = pos
			}
		} else {
			val := 0
			if c == '(' {
				val = 1
			} else if c == ')' {
				val = -1
			}
			if arr[pos] != int8(val) {
				arr[pos] = int8(val)
				update(1, 1, size, pos, val)
			}
			if pos > maxPos {
				maxPos = pos
			}
		}
		node := tree[1]
		if node.sum != 0 || node.min < 0 {
			fmt.Fprint(writer, -1)
		} else {
			fmt.Fprint(writer, node.max)
		}
		if i+1 < n {
			fmt.Fprint(writer, " ")
		}
	}
}
