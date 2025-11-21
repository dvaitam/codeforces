package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 100000

type Node struct {
	edges map[byte]int
	fail  int
	cnt   int64
}

var automaton []Node
var ptr int
var tree1Parent []int
var tree1Char []byte
var tree2Parent []int
var tree2Char []byte
var tree1Depth []int
var powBase []uint64
var prefixHash []uint64

func initAutomaton(n int) {
	automaton = make([]Node, n+5)
	automaton[0].edges = make(map[byte]int)
	ptr = 1
}

func goEdge(v int, ch byte) int {
	if nxt, ok := automaton[v].edges[ch]; ok {
		return nxt
	}
	automaton[ptr].edges = make(map[byte]int)
	automaton[v].edges[ch] = ptr
	ptr++
	return ptr - 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	tree1Parent = make([]int, maxN+2)
	tree1Char = make([]byte, maxN+2)
	tree1Depth = make([]int, maxN+2)

	tree2Parent = make([]int, maxN+2)
	tree2Char = make([]byte, maxN+2)

	initAutomaton(maxN)

	type Operation struct {
		t int
		v int
		c byte
	}

	ops := make([]Operation, n)
	for i := 0; i < n; i++ {
		var t, v int
		var c string
		fmt.Fscan(in, &t, &v, &c)
		ops[i] = Operation{t: t, v: v, c: c[0]}
	}

	tree1Size := 1
	tree2Size := 1
	ans := int64(0)
	for _, op := range ops {
		if op.t == 1 {
			tree1Size++
			tree1Parent[tree1Size] = op.v
			tree1Char[tree1Size] = op.c
			tree1Depth[tree1Size] = tree1Depth[op.v] + 1

			node := 0
			temp := tree1Size
			for temp > 1 {
				node = goEdge(node, tree1Char[temp])
				temp = tree1Parent[temp]
			}
			automaton[node].cnt++
		} else {
			tree2Size++
			tree2Parent[tree2Size] = op.v
			tree2Char[tree2Size] = op.c

			temp := tree2Size
			node := 0
			for temp > 0 {
				if nxt, ok := automaton[node].edges[tree2Char[temp]]; ok {
					node = nxt
					ans += automaton[node].cnt
				} else {
					break
				}
				temp = tree2Parent[temp]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
