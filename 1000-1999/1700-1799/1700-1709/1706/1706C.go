package main

import (
	"bufio"
	"fmt"
	"os"
)

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

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
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &h[i])
		}
		if n%2 == 1 {
			var ans int64
			for i := 1; i < n-1; i += 2 {
				need := max64(h[i-1], h[i+1]) + 1
				if need > h[i] {
					ans += need - h[i]
				}
			}
			fmt.Fprintln(writer, ans)
			continue
		}
		m := n - 2
		pairs := m / 2
		if pairs == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		cost := make([]int64, m+1) // 1-indexed
		for k := 1; k <= m; k++ {
			idx := k
			need := max64(h[idx-1], h[idx+1]) + 1
			if need > h[idx] {
				cost[k] = need - h[idx]
			}
		}
		dp0 := make([]int64, pairs+1)
		dp1 := make([]int64, pairs+1)
		dp0[1] = cost[1]
		dp1[1] = cost[2]
		for j := 2; j <= pairs; j++ {
			dp0[j] = cost[2*j-1] + dp0[j-1]
			prev := dp0[j-1]
			if dp1[j-1] < prev {
				prev = dp1[j-1]
			}
			dp1[j] = cost[2*j] + prev
		}
		ans := dp0[pairs]
		if dp1[pairs] < ans {
			ans = dp1[pairs]
		}
		fmt.Fprintln(writer, ans)
	}
}
