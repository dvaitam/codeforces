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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var p, b int
		fmt.Fscan(reader, &p, &b)

		tokens := make([]int, p)
		for i := 0; i < p; i++ {
			fmt.Fscan(reader, &tokens[i])
		}
		bonus := make([]bool, n+1)
		for i := 0; i < b; i++ {
			var x int
			fmt.Fscan(reader, &x)
			bonus[x] = true
		}

		graph := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			graph[u] = append(graph[u], v)
			graph[v] = append(graph[v], u)
		}

		won := false
		for _, v := range tokens {
			if v == 1 {
				won = true
				break
			}
		}
		if won {
			fmt.Fprintln(writer, "YES")
			continue
		}

		// neighbors of 1
		adj1 := make(map[int]bool)
		for _, to := range graph[1] {
			adj1[to] = true
		}
		for _, v := range tokens {
			if adj1[v] {
				won = true
				break
			}
		}
		if won {
			fmt.Fprintln(writer, "YES")
			continue
		}

		if p < 2 {
			fmt.Fprintln(writer, "NO")
			continue
		}

		// BFS from vertex 1 through bonus vertices only
		visited := make([]bool, n+1)
		q := []int{1}
		visited[1] = true
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			for _, to := range graph[v] {
				if !visited[to] && bonus[to] {
					visited[to] = true
					q = append(q, to)
				}
			}
		}

		for _, v := range tokens {
			for _, to := range graph[v] {
				if visited[to] {
					won = true
					break
				}
			}
			if won {
				break
			}
		}
		if won {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
