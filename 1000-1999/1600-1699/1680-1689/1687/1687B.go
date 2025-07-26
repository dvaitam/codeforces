package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This solution implements the interactive strategy for Codeforces problem
// 1687B "Railway System".  The program queries the judge to discover the
// length of every track and then applies a Kruskal like procedure to build
// a minimal spanning forest.  Finally the minimal possible total length is
// printed prefixed with '!'.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	_ = n

	type edge struct {
		w   int
		idx int
	}

	edges := make([]edge, m)
	// First query each edge individually to obtain its weight.
	for i := 0; i < m; i++ {
		query := make([]byte, m)
		for j := 0; j < m; j++ {
			if i == j {
				query[j] = '1'
			} else {
				query[j] = '0'
			}
		}
		fmt.Fprintf(out, "? %s\n", string(query))
		out.Flush()
		var res int
		if _, err := fmt.Fscan(in, &res); err != nil {
			return
		}
		edges[i] = edge{w: res, idx: i}
	}

	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })

	// Now build the minimal spanning forest.
	state := make([]byte, m)
	for i := range state {
		state[i] = '0'
	}
	cur := 0
	for _, e := range edges {
		state[e.idx] = '1'
		fmt.Fprintf(out, "? %s\n", string(state))
		out.Flush()
		var res int
		if _, err := fmt.Fscan(in, &res); err != nil {
			return
		}
		if res == cur+e.w {
			cur += e.w
		} else {
			state[e.idx] = '0'
		}
	}

	fmt.Fprintf(out, "! %d\n", cur)
}
