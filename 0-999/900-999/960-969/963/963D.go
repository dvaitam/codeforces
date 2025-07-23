package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	next [26]int
	link int
	out  []int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	fmt.Fscan(reader, &s)
	var n int
	fmt.Fscan(reader, &n)

	kVals := make([]int, n)
	patterns := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &kVals[i], &patterns[i])
	}

	nodes := make([]Node, 1)
	patLen := make([]int, n)
	for idx, p := range patterns {
		v := 0
		for i := 0; i < len(p); i++ {
			c := int(p[i] - 'a')
			if nodes[v].next[c] == 0 {
				nodes[v].next[c] = len(nodes)
				nodes = append(nodes, Node{})
			}
			v = nodes[v].next[c]
		}
		nodes[v].out = append(nodes[v].out, idx)
		patLen[idx] = len(p)
	}

	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := nodes[0].next[c]
		if v != 0 {
			nodes[v].link = 0
			queue = append(queue, v)
		}
	}

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if link := nodes[v].link; link != 0 {
			nodes[v].out = append(nodes[v].out, nodes[link].out...)
		}
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].link = nodes[nodes[v].link].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[nodes[v].link].next[c]
			}
		}
	}

	occ := make([][]int, n)
	state := 0
	for i := 0; i < len(s); i++ {
		c := int(s[i] - 'a')
		state = nodes[state].next[c]
		if len(nodes[state].out) > 0 {
			for _, id := range nodes[state].out {
				occ[id] = append(occ[id], i)
			}
		}
	}

	for i := 0; i < n; i++ {
		if len(occ[i]) < kVals[i] {
			fmt.Fprintln(writer, -1)
			continue
		}
		best := len(s) + 1
		need := kVals[i]
		pos := occ[i]
		L := patLen[i]
		for j := 0; j+need-1 < len(pos); j++ {
			cur := pos[j+need-1] - pos[j] + L
			if cur < best {
				best = cur
			}
		}
		fmt.Fprintln(writer, best)
	}
}
