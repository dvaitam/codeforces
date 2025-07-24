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
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		if x != 0 {
			arr = append(arr, x)
		}
	}
	n = len(arr)
	if n == 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	if n > 120 {
		fmt.Fprintln(writer, 3)
		return
	}
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i]&arr[j] != 0 {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}
	const inf = int(1e9)
	ans := inf
	for s := 0; s < n; s++ {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		parent := make([]int, n)
		queue := []int{s}
		dist[s] = 0
		for qi := 0; qi < len(queue); qi++ {
			u := queue[qi]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					parent[v] = u
					queue = append(queue, v)
				} else if parent[u] != v {
					cycle := dist[u] + dist[v] + 1
					if cycle < ans {
						ans = cycle
					}
				}
			}
		}
	}
	if ans == inf {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}
