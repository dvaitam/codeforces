package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	next [26]int
	fail int
	out  int
	pat  int
}

var nodes []node

func newNode() int {
	nodes = append(nodes, node{pat: -1})
	return len(nodes) - 1
}

func buildTrie(patterns []string) {
	nodes = make([]node, 1)
	nodes[0].pat = -1
	for idx, s := range patterns {
		v := 0
		for _, ch := range s {
			c := ch - 'a'
			if nodes[v].next[c] == 0 {
				nodes[v].next[c] = newNode()
			}
			v = nodes[v].next[c]
		}
		nodes[v].pat = idx
	}
	// build fail links
	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		u := nodes[0].next[c]
		if u != 0 {
			nodes[u].fail = 0
			queue = append(queue, u)
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		f := nodes[v].fail
		if nodes[f].pat != -1 {
			nodes[v].out = f
		} else {
			nodes[v].out = nodes[f].out
		}
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].fail = nodes[nodes[v].fail].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[nodes[v].fail].next[c]
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	names := make([]string, n)
	lens := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &names[i])
		lens[i] = len(names[i])
	}
	buildTrie(names)
	const INF = int(1e9)
	minSuper := make([]int, n)
	for i := range minSuper {
		minSuper[i] = INF
	}
	// first pass: compute minimal superstring length for each name
	for i, s := range names {
		v := 0
		for _, ch := range s {
			c := ch - 'a'
			v = nodes[v].next[c]
			u := v
			for u != 0 {
				if id := nodes[u].pat; id != -1 {
					if id != i && lens[i] > lens[id] {
						if lens[i] < minSuper[id] {
							minSuper[id] = lens[i]
						}
					}
				}
				u = nodes[u].out
			}
		}
	}
	seen := make([]int, n)
	iter := 1
	res := 0
	// second pass: count edges
	for i, s := range names {
		v := 0
		Li := lens[i]
		for _, ch := range s {
			c := ch - 'a'
			v = nodes[v].next[c]
			u := v
			for u != 0 {
				if id := nodes[u].pat; id != -1 {
					if id != i && Li > lens[id] && minSuper[id] == Li {
						if seen[id] != iter {
							seen[id] = iter
							res++
						}
					}
				}
				u = nodes[u].out
			}
		}
		iter++
	}
	fmt.Println(res)
}
