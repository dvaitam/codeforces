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
	neighbors := make([][]int, n+1)
	high := make([]map[int]struct{}, n+1)

	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		neighbors[a] = append(neighbors[a], b)
		neighbors[b] = append(neighbors[b], a)
	}

	small := make([]int, n+1)
	big := make([]int, n+1)

	for v := 1; v <= n; v++ {
		high[v] = make(map[int]struct{})
	}
	for v := 1; v <= n; v++ {
		for _, u := range neighbors[v] {
			if u > v { // rank[u]>rank[v]
				high[v][u] = struct{}{}
			}
		}
	}
	for v := 1; v <= n; v++ {
		small[v] = len(high[v])
		big[v] = len(neighbors[v]) - small[v]
	}
	var ans int64
	for v := 1; v <= n; v++ {
		ans += int64(small[v]) * int64(big[v])
	}

	var q int
	fmt.Fscan(reader, &q)
	fmt.Fprintln(writer, ans)
	for i := 0; i < q; i++ {
		var v int
		fmt.Fscan(reader, &v)
		// output ans after processing? Wait we need output before update at each day start. So we output ans after reading each query but before updating for next day. But we printed ans before first update above.
		// Process promotion of v
		for u := range high[v] {
			// remove contributions
			ans -= int64(small[v]) * int64(big[v])
			ans -= int64(small[u]) * int64(big[u])
			delete(high[v], u)
			small[v]--
			big[v]++

			if high[u] == nil {
				high[u] = make(map[int]struct{})
			}
			high[u][v] = struct{}{}
			small[u]++
			big[u]--

			ans += int64(small[v]) * int64(big[v])
			ans += int64(small[u]) * int64(big[u])
		}
		fmt.Fprintln(writer, ans)
	}
}
