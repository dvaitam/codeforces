package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	type Ant struct{ x, s int }
	ants := make([]Ant, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ants[i].x, &ants[i].s)
	}
	sort.Slice(ants, func(i, j int) bool { return ants[i].x < ants[j].x })

	const INF = int(1e9)
	dp := make([]int, m+1)
	for i := 1; i <= m; i++ {
		dp[i] = INF
	}
	for _, a := range ants {
		newdp := make([]int, m+1)
		copy(newdp, dp)
		best := make([]int, m+1)
		for i := 0; i <= m; i++ {
			best[i] = INF
		}
		for pos := 0; pos < m; pos++ {
			if dp[pos] == INF {
				continue
			}
			need := pos + 1
			var d int
			if need < a.x-a.s {
				d = (a.x - a.s) - need
			} else {
				d = 0
			}
			r0 := a.x + a.s + d
			if r0 > m {
				r0 = m
			}
			val := dp[pos] + d - r0
			if val < best[r0] {
				best[r0] = val
			}
		}
		pref := INF
		for r := 0; r <= m; r++ {
			if best[r] < pref {
				pref = best[r]
			}
			if pref == INF {
				continue
			}
			if cand := pref + r; cand < newdp[r] {
				newdp[r] = cand
			}
		}
		dp = newdp
	}
	fmt.Println(dp[m])
}
