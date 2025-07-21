package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// read n, m
	var n, m int
	fmt.Fscan(reader, &n, &m)

	g := make([][]int, n)
	grev := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		grev[v] = append(grev[v], u)
	}
	// first pass order
	visited := make([]bool, n)
	order := make([]int, 0, n)
	type frame1 struct{ v, i int }
	for s := 0; s < n; s++ {
		if visited[s] {
			continue
		}
		stack := []frame1{{s, 0}}
		for len(stack) > 0 {
			fr := &stack[len(stack)-1]
			v, i := fr.v, fr.i
			if i == 0 {
				visited[v] = true
			}
			if fr.i < len(g[v]) {
				to := g[v][fr.i]
				fr.i++
				if !visited[to] {
					stack = append(stack, frame1{to, 0})
				}
			} else {
				order = append(order, v)
				stack = stack[:len(stack)-1]
			}
		}
	}
	// second pass assign components
	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	cid := 0
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] != -1 {
			continue
		}
		// dfs on grev
		stack := []int{v}
		comp[v] = cid
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, w := range grev[u] {
				if comp[w] == -1 {
					comp[w] = cid
					stack = append(stack, w)
				}
			}
		}
		cid++
	}
	// group nodes per comp
	nodes := make([][]int, cid)
	for v := 0; v < n; v++ {
		nodes[comp[v]] = append(nodes[comp[v]], v)
	}
	// compute global gcd
	globalG := 0
	depth := make([]int, n)
	mark := make([]bool, n)
	for c := 0; c < cid; c++ {
		if len(nodes[c]) == 0 {
			continue
		}
		// reset mark and depth for this comp
		for _, v := range nodes[c] {
			mark[v] = false
			depth[v] = 0
		}
		// dfs for cycle lengths
		var stack2 []int
		root := nodes[c][0]
		depth[root] = 0
		mark[root] = true
		stack2 = append(stack2, root)
		localG := 0
		for len(stack2) > 0 {
			u := stack2[len(stack2)-1]
			stack2 = stack2[:len(stack2)-1]
			for _, v := range g[u] {
				if comp[v] != c {
					continue
				}
				if !mark[v] {
					mark[v] = true
					depth[v] = depth[u] + 1
					stack2 = append(stack2, v)
				} else {
					d := depth[u] + 1 - depth[v]
					if d < 0 {
						d = -d
					}
					localG = gcd(localG, d)
				}
			}
		}
		globalG = gcd(globalG, localG)
	}
	if globalG == 0 {
		globalG = n
	}
	fmt.Fprintln(writer, globalG)
}
