package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	l, r int
}

func solve() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	intervals := make([]pair, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &intervals[i].l, &intervals[i].r)
	}

	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].r != intervals[j].r {
			return intervals[i].r < intervals[j].r
		}
		return intervals[i].l < intervals[j].l
	})

	type state struct {
		end int
		val int
	}

	dp := make([][]state, k+1)
	dp[0] = append(dp[0], state{n + 1, 0})

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for i := m - 1; i >= 0; i-- {
		l, r := intervals[i].l, intervals[i].r
		for j := k; j >= 0; j-- {
			mx := -1
			
			var toMove []state
			for len(dp[j]) > 0 {
				last := dp[j][len(dp[j])-1]
				if l >= last.end {
					break
				}
				toMove = append(toMove, last)
				dp[j] = dp[j][:len(dp[j])-1]

				val := max(last.end-r-1, 0) + last.val
				if val > mx {
					mx = val
				}
			}

			if j+1 <= k {
				newDp := append(dp[j+1], toMove...)
				sort.Slice(newDp, func(a, b int) bool {
					if newDp[a].end != newDp[b].end {
						return newDp[a].end < newDp[b].end
					}
					return newDp[a].val < newDp[b].val
				})
				
				if len(newDp) > 0 {
					uniqueDp := []state{newDp[0]}
					for i := 1; i < len(newDp); i++ {
						if newDp[i].end != newDp[i-1].end {
							uniqueDp = append(uniqueDp, newDp[i])
						} else {
							if newDp[i].val > uniqueDp[len(uniqueDp)-1].val {
								uniqueDp[len(uniqueDp)-1] = newDp[i]
							}
						}
					}
					dp[j+1] = uniqueDp
				} else {
					dp[j+1] = nil
				}
			}

			if mx >= 0 {
				dp[j] = append(dp[j], state{l, mx})
				sort.Slice(dp[j], func(a, b int) bool {
					if dp[j][a].end != dp[j][b].end {
						return dp[j][a].end < dp[j][b].end
					}
					return dp[j][a].val < dp[j][b].val
				})
			}
		}
	}

	ans := 0
	for _, s := range dp[k] {
		ans = max(ans, s.end+s.val-1)
	}
	fmt.Fprintln(out, ans)
}

func main() {
	var t int
	fmt.Fscan(os.Stdin, &t)
	for ; t > 0; t-- {
		solve()
	}
}
