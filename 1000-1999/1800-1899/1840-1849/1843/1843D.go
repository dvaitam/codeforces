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
		var n int
		fmt.Fscan(reader, &n)
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		parent := make([]int, n+1)
		order := make([]int, 0, n)
		queue := []int{1}
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			order = append(order, v)
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				queue = append(queue, to)
			}
		}

		leaf := make([]int64, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			count := int64(0)
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				count += leaf[to]
			}
			if count == 0 {
				leaf[v] = 1
			} else {
				leaf[v] = count
			}
		}

		var q int
		fmt.Fscan(reader, &q)
		for ; q > 0; q-- {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			fmt.Fprintln(writer, leaf[x]*leaf[y])
		}
	}
}
