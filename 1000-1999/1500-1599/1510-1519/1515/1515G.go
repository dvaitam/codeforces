package main

import (
	"bufio"
	"fmt"
	"os"
)

// gcd computes the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// Edge represents a weighted directed edge.
type Edge struct {
	to int
	w  int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)

	adj := make([][]Edge, n)
	radj := make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		var l int64
		fmt.Fscan(reader, &a, &b, &l)
		a--
		b--
		adj[a] = append(adj[a], Edge{b, l})
		radj[b] = append(radj[b], a)
	}

	// First pass: order nodes using DFS (iterative)
	visited := make([]bool, n)
	order := make([]int, 0, n)
	type StackFrame struct{ v, idx int }
	for v := 0; v < n; v++ {
		if visited[v] {
			continue
		}
		stack := []StackFrame{{v, 0}}
		visited[v] = true
		for len(stack) > 0 {
			f := &stack[len(stack)-1]
			if f.idx < len(adj[f.v]) {
				to := adj[f.v][f.idx].to
				f.idx++
				if !visited[to] {
					visited[to] = true
					stack = append(stack, StackFrame{to, 0})
				}
			} else {
				order = append(order, f.v)
				stack = stack[:len(stack)-1]
			}
		}
	}

	// Second pass: assign components using reverse edges
	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	compCnt := 0
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] != -1 {
			continue
		}
		stack := []int{v}
		comp[v] = compCnt
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, p := range radj[u] {
				if comp[p] == -1 {
					comp[p] = compCnt
					stack = append(stack, p)
				}
			}
		}
		compCnt++
	}

	// Compute gcd for each component
	compGCD := make([]int64, compCnt)
	dist := make([]int64, n)
	visited2 := make([]bool, n)
	for v := 0; v < n; v++ {
		if visited2[v] {
			continue
		}
		id := comp[v]
		stack := []int{v}
		visited2[v] = true
		dist[v] = 0
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, e := range adj[u] {
				if comp[e.to] != id {
					continue
				}
				if !visited2[e.to] {
					visited2[e.to] = true
					dist[e.to] = dist[u] + e.w
					stack = append(stack, e.to)
				} else {
					diff := dist[u] + e.w - dist[e.to]
					if diff < 0 {
						diff = -diff
					}
					compGCD[id] = gcd(compGCD[id], diff)
				}
			}
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	for i := 0; i < q; i++ {
		var v int
		var s, t int64
		fmt.Fscan(reader, &v, &s, &t)
		v--
		id := comp[v]
		g := gcd(compGCD[id], t)
		if s%g == 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
