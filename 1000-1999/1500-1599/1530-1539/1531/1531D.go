package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	lock  [2][7]int
	color [2][7]int
}

var colorIndex = map[string]int{
	"red":    0,
	"orange": 1,
	"yellow": 2,
	"green":  3,
	"blue":   4,
	"indigo": 5,
	"violet": 6,
}

var colorNames = []string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"}

func makeNode(msg string) Node {
	var node Node
	for l := 0; l < 2; l++ {
		for c := 0; c < 7; c++ {
			switch msg {
			case "lock":
				node.lock[l][c] = 1
				node.color[l][c] = c
			case "unlock":
				node.lock[l][c] = 0
				node.color[l][c] = c
			default:
				node.lock[l][c] = l
				if l == 0 {
					node.color[l][c] = colorIndex[msg]
				} else {
					node.color[l][c] = c
				}
			}
		}
	}
	return node
}

func combine(a, b Node) Node {
	var res Node
	for l := 0; l < 2; l++ {
		for c := 0; c < 7; c++ {
			ml := a.lock[l][c]
			mc := a.color[l][c]
			res.lock[l][c] = b.lock[ml][mc]
			res.color[l][c] = b.color[ml][mc]
		}
	}
	return res
}

var tree []Node
var msgs []string

func build(v, l, r int) {
	if l == r {
		tree[v] = makeNode(msgs[l])
		return
	}
	m := (l + r) / 2
	build(v*2, l, m)
	build(v*2+1, m+1, r)
	tree[v] = combine(tree[v*2], tree[v*2+1])
}

func update(v, l, r, pos int, msg string) {
	if l == r {
		tree[v] = makeNode(msg)
		return
	}
	m := (l + r) / 2
	if pos <= m {
		update(v*2, l, m, pos, msg)
	} else {
		update(v*2+1, m+1, r, pos, msg)
	}
	tree[v] = combine(tree[v*2], tree[v*2+1])
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	msgs = make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &msgs[i])
	}

	tree = make([]Node, 4*n)
	build(1, 0, n-1)

	var t int
	fmt.Fscan(in, &t)

	initialColor := colorNames[tree[1].color[0][colorIndex["blue"]]]
	fmt.Fprintln(out, initialColor)

	for ; t > 0; t-- {
		var idx int
		var msg string
		fmt.Fscan(in, &idx, &msg)
		msgs[idx-1] = msg
		update(1, 0, n-1, idx-1, msg)
		resColor := colorNames[tree[1].color[0][colorIndex["blue"]]]
		fmt.Fprintln(out, resColor)
	}
}
