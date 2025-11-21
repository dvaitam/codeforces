package main

import (
	"bufio"
	"fmt"
	"os"
)

type instruction struct {
	typ  int
	node int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		depth := make([]int, n+1)
		for i := range depth {
			depth[i] = -1
		}
		queue := make([]int, 0, n)
		depth[1] = 0
		queue = append(queue, 1)
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, v := range adj[u] {
				if depth[v] == -1 {
					depth[v] = depth[u] + 1
					queue = append(queue, v)
				}
			}
		}

		deg := make([]int, n+1)
		for i := 1; i <= n; i++ {
			deg[i] = len(adj[i])
		}

		removed := make([]bool, n+1)
		leaves := make([]int, 0)
		for i := 1; i <= n; i++ {
			if i != n && deg[i] <= 1 {
				leaves = append(leaves, i)
			}
		}

		ops := make([]instruction, 0, 3*n)
		movesParity := 0
		lastWasDelete := false
		processed := 0
		head := 0

		// Helper to append a move instruction.
		addMove := func() {
			ops = append(ops, instruction{typ: 1})
			movesParity ^= 1
			lastWasDelete = false
		}

		for processed < n-1 {
			var u int
			found := false
			for head < len(leaves) {
				cand := leaves[head]
				head++
				if !removed[cand] && deg[cand] == 1 {
					u = cand
					found = true
					break
				}
			}
			if !found {
				for i := 1; i <= n; i++ {
					if i != n && !removed[i] && deg[i] == 1 {
						u = i
						found = true
						break
					}
				}
			}
			if !found {
				break
			}

			neighbor := -1
			for _, v := range adj[u] {
				if !removed[v] {
					neighbor = v
					break
				}
			}

			if lastWasDelete {
				addMove()
			}
			for movesParity == (depth[u] & 1) {
				addMove()
			}

			ops = append(ops, instruction{typ: 2, node: u})
			lastWasDelete = true

			removed[u] = true
			processed++
			deg[u] = 0
			if neighbor != -1 {
				deg[neighbor]--
				if neighbor != n && deg[neighbor] == 1 {
					leaves = append(leaves, neighbor)
				}
			}
		}

		fmt.Fprintln(out, len(ops))
		for _, op := range ops {
			if op.typ == 1 {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintf(out, "2 %d\n", op.node)
			}
		}
	}
}
