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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		// Precompute prefix and suffix max/min to speed up lookups of outside ranges
		prefixMax := make([]int64, n+1)
		prefixMin := make([]int64, n+1)
		const INF = int64(1<<63 - 1)
		prefixMax[0] = -INF
		prefixMin[0] = INF
		for i := 1; i <= n; i++ {
			val := a[i-1]
			if val > prefixMax[i-1] {
				prefixMax[i] = val
			} else {
				prefixMax[i] = prefixMax[i-1]
			}
			if val < prefixMin[i-1] {
				prefixMin[i] = val
			} else {
				prefixMin[i] = prefixMin[i-1]
			}
		}
		suffixMax := make([]int64, n+2)
		suffixMin := make([]int64, n+2)
		suffixMax[n+1] = -INF
		suffixMin[n+1] = INF
		for i := n; i >= 1; i-- {
			val := a[i-1]
			if val > suffixMax[i+1] {
				suffixMax[i] = val
			} else {
				suffixMax[i] = suffixMax[i+1]
			}
			if val < suffixMin[i+1] {
				suffixMin[i] = val
			} else {
				suffixMin[i] = suffixMin[i+1]
			}
		}

		var ans int64
		for l := 0; l < n; l++ {
			maxIn := a[l]
			minIn := a[l]
			for r := l; r < n; r++ {
				if r-l+1 == n {
					break
				}
				if a[r] > maxIn {
					maxIn = a[r]
				}
				if a[r] < minIn {
					minIn = a[r]
				}
				outsideMax := prefixMax[l]
				if suffixMax[r+2] > outsideMax {
					outsideMax = suffixMax[r+2]
				}
				outsideMin := prefixMin[l]
				if suffixMin[r+2] < outsideMin {
					outsideMin = suffixMin[r+2]
				}
				cur := outsideMax - outsideMin + maxIn - minIn
				if cur > ans {
					ans = cur
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
