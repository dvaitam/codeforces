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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	w := make([][]int, m)
	for i := range w {
		w[i] = make([]int, m)
	}
	for i := 1; i < n; i++ {
		a := int(s[i-1] - 'a')
		b := int(s[i] - 'a')
		if a == b || a >= m || b >= m {
			continue
		}
		w[a][b]++
		w[b][a]++
	}

	total := make([]int, m)
	for i := 0; i < m; i++ {
		sum := 0
		for j := 0; j < m; j++ {
			sum += w[i][j]
		}
		total[i] = sum
	}

	maxMask := 1 << m

	sumW := make([][]int, m)
	for j := 0; j < m; j++ {
		arr := make([]int, maxMask)
		for mask := 1; mask < maxMask; mask++ {
			lb := mask & -mask
			k := bits.TrailingZeros(uint(lb))
			arr[mask] = arr[mask^lb] + w[j][k]
		}
		sumW[j] = arr
	}

	cross := make([]int, maxMask)
	for mask := 1; mask < maxMask; mask++ {
		lb := mask & -mask
		j := bits.TrailingZeros(uint(lb))
		prev := mask ^ lb
		cross[mask] = cross[prev] + total[j] - 2*sumW[j][prev]
	}

	const inf int = int(1e18)
	dp := make([]int, maxMask)
	for i := 1; i < maxMask; i++ {
		dp[i] = inf
	}

	for mask := 0; mask < maxMask; mask++ {
		for j := 0; j < m; j++ {
			if mask&(1<<j) == 0 {
				next := mask | (1 << j)
				val := dp[mask] + cross[next]
				if val < dp[next] {
					dp[next] = val
				}
			}
		}
	}

	fmt.Fprintln(out, dp[maxMask-1])
}
