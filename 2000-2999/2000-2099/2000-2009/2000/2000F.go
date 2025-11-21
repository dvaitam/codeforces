package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const INF = math.MaxInt64 / 4

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		bestRect := make([][]int, n)
		for i := 0; i < n; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			limit := k
			best := make([]int, limit+1)
			for j := 1; j <= limit; j++ {
				best[j] = INF
			}
			best[0] = 0
			for r := 0; r <= b; r++ {
				for c := 0; c <= a; c++ {
					if r == 0 && c == 0 {
						continue
					}
					s := r + c
					cost := r*a + c*b - r*c
					idx := s
					if idx > limit {
						idx = limit
					}
					if cost < best[idx] {
						best[idx] = cost
					}
				}
			}
			bestRect[i] = best
		}

		dp := make([]int, k+1)
		for i := 1; i <= k; i++ {
			dp[i] = INF
		}
		for i := 0; i < n; i++ {
			best := bestRect[i]
			newdp := make([]int, k+1)
			copy(newdp, dp)
			for cur := 0; cur <= k; cur++ {
				if dp[cur] == INF {
					continue
				}
				for gain := 1; gain <= k; gain++ {
					cost := best[gain]
					if cost == INF {
						continue
					}
					nxt := cur + gain
					if nxt > k {
						nxt = k
					}
					v := dp[cur] + cost
					if v < newdp[nxt] {
						newdp[nxt] = v
					}
				}
			}
			dp = newdp
		}
		if dp[k] >= INF {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, dp[k])
		}
	}
}
