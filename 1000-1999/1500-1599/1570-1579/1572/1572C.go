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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		// compress consecutive duplicates
		b := make([]int, 0, n)
		for _, v := range arr {
			if len(b) == 0 || b[len(b)-1] != v {
				b = append(b, v)
			}
		}
		m := len(b)
		if m == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		pos := make(map[int][]int)
		for i, v := range b {
			pos[v] = append(pos[v], i)
		}
		dp := make([][]int, m)
		for i := range dp {
			dp[i] = make([]int, m)
		}
		for length := 2; length <= m; length++ {
			for l := 0; l+length-1 < m; l++ {
				r := l + length - 1
				best := dp[l+1][r] + 1
				val := b[l]
				for _, k := range pos[val] {
					if k <= l || k > r {
						continue
					}
					cand := dp[l][k-1] + dp[k][r]
					if cand < best {
						best = cand
					}
				}
				dp[l][r] = best
			}
		}
		fmt.Fprintln(out, dp[0][m-1])
	}
}
