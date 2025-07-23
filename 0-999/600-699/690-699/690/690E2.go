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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var h, w, k int
		fmt.Fscan(in, &h, &w, &k)
		pieceH := h / k
		top := make([][]int, k)
		bottom := make([][]int, k)
		for i := 0; i < k; i++ {
			top[i] = make([]int, w)
			bottom[i] = make([]int, w)
			for r := 0; r < pieceH; r++ {
				for c := 0; c < w; c++ {
					var val int
					fmt.Fscan(in, &val)
					if r == 0 {
						top[i][c] = val
					}
					if r == pieceH-1 {
						bottom[i][c] = val
					}
				}
			}
		}

		// compute pairwise costs
		cost := make([][]int, k)
		for i := 0; i < k; i++ {
			cost[i] = make([]int, k)
			for j := 0; j < k; j++ {
				if i == j {
					continue
				}
				diff := 0
				for c := 0; c < w; c++ {
					x := bottom[i][c] - top[j][c]
					if x < 0 {
						x = -x
					}
					diff += x
				}
				cost[i][j] = diff
			}
		}

		nMask := 1 << k
		const INF int = int(1e9)
		dp := make([][]int, nMask)
		prev := make([][]int, nMask)
		for i := range dp {
			dp[i] = make([]int, k)
			prev[i] = make([]int, k)
			for j := range dp[i] {
				dp[i][j] = INF
				prev[i][j] = -1
			}
		}
		for i := 0; i < k; i++ {
			dp[1<<i][i] = 0
		}

		for mask := 1; mask < nMask; mask++ {
			for last := 0; last < k; last++ {
				if dp[mask][last] == INF {
					continue
				}
				for nxt := 0; nxt < k; nxt++ {
					if mask&(1<<nxt) != 0 {
						continue
					}
					nm := mask | (1 << nxt)
					cand := dp[mask][last] + cost[last][nxt]
					if cand < dp[nm][nxt] {
						dp[nm][nxt] = cand
						prev[nm][nxt] = last
					}
				}
			}
		}

		full := nMask - 1
		bestVal := INF
		bestLast := 0
		for last := 0; last < k; last++ {
			if dp[full][last] < bestVal {
				bestVal = dp[full][last]
				bestLast = last
			}
		}

		order := make([]int, k)
		mask := full
		last := bestLast
		for idx := k - 1; idx >= 0; idx-- {
			order[idx] = last
			pl := prev[mask][last]
			mask ^= 1 << last
			last = pl
			if last == -1 {
				break
			}
		}

		ans := make([]int, k)
		for pos, piece := range order {
			ans[piece] = pos + 1
		}
		for i := 0; i < k; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
