package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1_000_000_007
const NEG_INF int = -1_000_000_000

type state struct {
	len int
	cnt int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		colors := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &colors[i])
		}

		// dpPrev[r][col] stores best from position > current with r elements
		// already taken in the current block of color col
		dpPrev := make([][]state, k)
		dpCur := make([][]state, k)
		for r := 0; r < k; r++ {
			dpPrev[r] = make([]state, n+1)
			dpCur[r] = make([]state, n+1)
			for c := 0; c <= n; c++ {
				dpPrev[r][c].len = NEG_INF
				dpCur[r][c].len = NEG_INF
			}
		}
		dpPrev[0][0] = state{0, 1}

		for idx := n - 1; idx >= 0; idx-- {
			// reset dpCur
			for r := 0; r < k; r++ {
				for c := 0; c <= n; c++ {
					dpCur[r][c].len = NEG_INF
					dpCur[r][c].cnt = 0
				}
			}

			color := colors[idx]
			for r := 0; r < k; r++ {
				for c := 0; c <= n; c++ {
					best := dpPrev[r][c]
					if best.len == NEG_INF {
						// state is unreachable, but continue to set skip result
						dpCur[r][c] = best
						continue
					}

					// option to skip current tile
					curLen := best.len
					curCnt := best.cnt

					// option to take current tile if possible
					if r == 0 {
						nr := 1
						nc := color
						if nr == k {
							nr = 0
							nc = 0
						}
						cand := dpPrev[nr][nc]
						if cand.len != NEG_INF {
							candLen := cand.len + 1
							candCnt := cand.cnt
							if candLen > curLen {
								curLen = candLen
								curCnt = candCnt
							} else if candLen == curLen {
								curCnt = (curCnt + candCnt) % MOD
							}
						}
					} else if c == color {
						nr := r + 1
						nc := c
						if nr == k {
							nr = 0
							nc = 0
						}
						cand := dpPrev[nr][nc]
						if cand.len != NEG_INF {
							candLen := cand.len + 1
							candCnt := cand.cnt
							if candLen > curLen {
								curLen = candLen
								curCnt = candCnt
							} else if candLen == curLen {
								curCnt = (curCnt + candCnt) % MOD
							}
						}
					}

					dpCur[r][c].len = curLen
					dpCur[r][c].cnt = curCnt % MOD
				}
			}

			dpPrev, dpCur = dpCur, dpPrev
		}

		ans := dpPrev[0][0].cnt % MOD
		fmt.Fprintln(out, ans)
	}
}
