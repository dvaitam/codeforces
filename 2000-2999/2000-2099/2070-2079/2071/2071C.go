package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, st, en int
		fmt.Fscan(in, &n, &st, &en)

		graph := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			graph[u] = append(graph[u], v)
			graph[v] = append(graph[v], u)
		}

		depth := make([]int, n+1)
		for i := range depth {
			depth[i] = -1
		}

		queue := make([]int, 0, n)
		queue = append(queue, en)
		depth[en] = 0
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			for _, to := range graph[v] {
				if depth[to] == -1 {
					depth[to] = depth[v] + 1
					queue = append(queue, to)
				}
			}
		}

		order := make([]int, n)
		for i := 1; i <= n; i++ {
			order[i-1] = i
		}
		sort.Slice(order, func(i, j int) bool {
			if depth[order[i]] == depth[order[j]] {
				return order[i] < order[j]
			}
			return depth[order[i]] > depth[order[j]]
		})

		for i, v := range order {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
