package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		for j := 0; j < n; j++ {
			if s[j] == '1' {
				adj[i] = append(adj[i], j)
			}
		}
	}

	var m int
	fmt.Fscan(in, &m)
	path := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &path[i])
		path[i]--
	}

	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = -1
		}
		queue := []int{i}
		dist[i][i] = 0
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, v := range adj[u] {
				if dist[i][v] == -1 {
					dist[i][v] = dist[i][u] + 1
					queue = append(queue, v)
				}
			}
		}
	}

	ans := []int{path[0]}
	last := 0
	for i := 1; i < m-1; i++ {
		if dist[path[last]][path[i+1]] < i+1-last {
			ans = append(ans, path[i])
			last = i
		}
	}
	ans = append(ans, path[m-1])

	fmt.Fprintln(out, len(ans))
	for i, v := range ans {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v+1)
	}
	fmt.Fprintln(out)
}
