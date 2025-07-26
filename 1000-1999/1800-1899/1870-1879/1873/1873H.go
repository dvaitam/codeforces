package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, a, b int
		fmt.Fscan(reader, &n, &a, &b)
		adj := make([][]int, n+1)
		deg := make([]int, n+1)
		for i := 0; i < n; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			deg[u]++
			deg[v]++
		}

		// find cycle nodes by removing leaves
		queue := make([]int, 0, n)
		removed := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			if deg[i] == 1 {
				queue = append(queue, i)
			}
		}
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			removed[v] = true
			for _, to := range adj[v] {
				if removed[to] {
					continue
				}
				deg[to]--
				if deg[to] == 1 {
					queue = append(queue, to)
				}
			}
		}
		onCycle := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			if !removed[i] {
				onCycle[i] = true
			}
		}

		// distances from a
		distA := make([]int, n+1)
		for i := range distA {
			distA[i] = -1
		}
		qa := make([]int, 0, n)
		distA[a] = 0
		qa = append(qa, a)
		for h := 0; h < len(qa); h++ {
			v := qa[h]
			for _, to := range adj[v] {
				if distA[to] == -1 {
					distA[to] = distA[v] + 1
					qa = append(qa, to)
				}
			}
		}

		// distances from b
		distB := make([]int, n+1)
		for i := range distB {
			distB[i] = -1
		}
		qb := make([]int, 0, n)
		distB[b] = 0
		qb = append(qb, b)
		for h := 0; h < len(qb); h++ {
			v := qb[h]
			for _, to := range adj[v] {
				if distB[to] == -1 {
					distB[to] = distB[v] + 1
					qb = append(qb, to)
				}
			}
		}

		// BFS from b only through nodes where Valeriu is strictly closer
		if distB[b] >= distA[b] {
			fmt.Fprintln(writer, "NO")
			continue
		}
		visited := make([]bool, n+1)
		q := []int{b}
		visited[b] = true
		escape := false
		for head := 0; head < len(q); head++ {
			v := q[head]
			if onCycle[v] {
				escape = true
				break
			}
			for _, to := range adj[v] {
				if !visited[to] && distB[to] < distA[to] {
					visited[to] = true
					q = append(q, to)
				}
			}
		}
		if escape {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
