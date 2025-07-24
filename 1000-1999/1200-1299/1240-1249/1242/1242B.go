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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([]map[int]struct{}, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if adj[a] == nil {
			adj[a] = make(map[int]struct{})
		}
		if adj[b] == nil {
			adj[b] = make(map[int]struct{})
		}
		adj[a][b] = struct{}{}
		adj[b][a] = struct{}{}
	}

	unvis := make(map[int]struct{}, n)
	for i := 1; i <= n; i++ {
		unvis[i] = struct{}{}
	}

	components := 0
	queue := make([]int, 0)
	for len(unvis) > 0 {
		// pick arbitrary start vertex
		var start int
		for k := range unvis {
			start = k
			break
		}
		queue = append(queue, start)
		delete(unvis, start)

		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]

			for u := range unvis {
				if _, ok := adj[v][u]; !ok {
					queue = append(queue, u)
					delete(unvis, u)
				}
			}
		}
		components++
	}

	fmt.Fprintln(writer, components-1)
}
