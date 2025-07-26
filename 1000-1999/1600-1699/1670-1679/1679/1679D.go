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
	var k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	a := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
	}

	good := func(x int) bool {
		// filter nodes with value <= x
		nodes := make([]int, 0)
		for i := 0; i < n; i++ {
			if a[i] <= x {
				nodes = append(nodes, i)
			}
		}
		if len(nodes) == 0 {
			return false
		}
		indeg := make([]int, n)
		for u := 0; u < n; u++ {
			if a[u] > x {
				continue
			}
			for _, v := range adj[u] {
				if a[v] <= x {
					indeg[v]++
				}
			}
		}
		queue := make([]int, 0)
		for _, v := range nodes {
			if indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
		dp := make([]int64, n)
		for _, v := range queue {
			dp[v] = 1
		}
		visited := 0
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			visited++
			for _, v := range adj[u] {
				if a[v] > x {
					continue
				}
				if dp[v] < dp[u]+1 {
					dp[v] = dp[u] + 1
				}
				indeg[v]--
				if indeg[v] == 0 {
					queue = append(queue, v)
				}
			}
		}
		if visited < len(nodes) {
			// cycle exists
			return true
		}
		maxLen := int64(0)
		for _, v := range nodes {
			if dp[v] > maxLen {
				maxLen = dp[v]
			}
		}
		return maxLen >= k
	}

	if !good(maxA) {
		fmt.Fprintln(writer, -1)
		return
	}
	lo, hi := 1, maxA
	for lo < hi {
		mid := (lo + hi) / 2
		if good(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Fprintln(writer, lo)
}
