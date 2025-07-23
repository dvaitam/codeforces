package main

import (
	"bufio"
	"fmt"
	"os"
)

const alphabet = 26

type Node struct {
	next     [alphabet]int
	students []int
	pseuds   []int
}

var (
	nodes []Node
	pairs [][2]int
	total int
)

func newNode() int {
	nodes = append(nodes, Node{})
	return len(nodes) - 1
}

func add(word string, idx int, isStudent bool) {
	v := 0
	for i := 0; i < len(word); i++ {
		c := word[i] - 'a'
		if nodes[v].next[c] == 0 {
			nodes[v].next[c] = newNode()
		}
		v = nodes[v].next[c]
	}
	if isStudent {
		nodes[v].students = append(nodes[v].students, idx)
	} else {
		nodes[v].pseuds = append(nodes[v].pseuds, idx)
	}
}

func dfs(v, depth int) {
	for i := 0; i < alphabet; i++ {
		child := nodes[v].next[i]
		if child != 0 {
			dfs(child, depth+1)
			if len(nodes[child].students) > 0 {
				nodes[v].students = append(nodes[v].students, nodes[child].students...)
			}
			if len(nodes[child].pseuds) > 0 {
				nodes[v].pseuds = append(nodes[v].pseuds, nodes[child].pseuds...)
			}
			nodes[child].students = nil
			nodes[child].pseuds = nil
		}
	}
	for len(nodes[v].students) > 0 && len(nodes[v].pseuds) > 0 {
		s := nodes[v].students[len(nodes[v].students)-1]
		nodes[v].students = nodes[v].students[:len(nodes[v].students)-1]
		p := nodes[v].pseuds[len(nodes[v].pseuds)-1]
		nodes[v].pseuds = nodes[v].pseuds[:len(nodes[v].pseuds)-1]
		pairs = append(pairs, [2]int{s, p})
		total += depth
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	nodes = make([]Node, 1)
	for i := 1; i <= n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		add(s, i, true)
	}
	for i := 1; i <= n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		add(s, i, false)
	}
	dfs(0, 0)
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, total)
	for _, pr := range pairs {
		fmt.Fprintln(writer, pr[0], pr[1])
	}
	writer.Flush()
}
