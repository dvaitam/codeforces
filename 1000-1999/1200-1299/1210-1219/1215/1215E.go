package main

import (
	"bufio"
	"fmt"
	"math/bits"
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

	arr := make([]int, n)
	freq20 := make([]int, 20)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		arr[i]--
		freq20[arr[i]]++
	}

	idx := make([]int, 20)
	m := 0
	for c := 0; c < 20; c++ {
		if freq20[c] > 0 {
			idx[c] = m
			m++
		} else {
			idx[c] = -1
		}
	}

	comp := make([]int, n)
	for i := 0; i < n; i++ {
		comp[i] = idx[arr[i]]
	}

	cross := make([][]int64, m)
	for i := range cross {
		cross[i] = make([]int64, m)
	}
	cnt := make([]int64, m)
	for _, c := range comp {
		for j := 0; j < m; j++ {
			cross[j][c] += cnt[j]
		}
		cnt[c]++
	}

	size := 1 << m
	cost := make([][]int64, m)
	for i := 0; i < m; i++ {
		cost[i] = make([]int64, size)
		for mask := 1; mask < size; mask++ {
			b := mask & -mask
			j := bits.TrailingZeros(uint(b))
			prev := mask ^ b
			cost[i][mask] = cost[i][prev] + cross[i][j]
		}
	}

	const INF int64 = 1 << 60
	dp := make([]int64, size)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for mask := 0; mask < size; mask++ {
		cur := dp[mask]
		if cur == INF {
			continue
		}
		for c := 0; c < m; c++ {
			if mask&(1<<c) == 0 {
				nm := mask | (1 << c)
				val := cur + cost[c][mask]
				if val < dp[nm] {
					dp[nm] = val
				}
			}
		}
	}

	fmt.Fprintln(out, dp[size-1])
}
