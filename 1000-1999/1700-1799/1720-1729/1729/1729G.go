package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1_000_000_007

func solveCase(s, t string) (int, int64) {
	L := len(t)
	positions := make([]int, 0)
	for i := 0; i+L <= len(s); i++ {
		if s[i:i+L] == t {
			positions = append(positions, i)
		}
	}
	m := len(positions)
	if m == 0 {
		return 0, 1
	}

	right := make([]int, m)
	j := 0
	for i := 0; i < m; i++ {
		if j < i {
			j = i
		}
		for j+1 < m && positions[j+1] <= positions[i]+L-1 {
			j++
		}
		right[i] = j
	}

	next := make([]int, m)
	for i := 0; i < m; i++ {
		x := positions[i] + L
		next[i] = sort.Search(len(positions), func(p int) bool { return positions[p] >= x })
	}

	const INF = int(1 << 30)
	dp := make([]int, m+1)
	cnt := make([]int64, m+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[m] = 0
	cnt[m] = 1

	for i := m - 1; i >= 0; i-- {
		r := right[i]
		best := INF
		var bestCnt int64
		for k := i; k <= r; k++ {
			j := next[k]
			cand := 1 + dp[j]
			if cand < best {
				best = cand
				bestCnt = cnt[j]
			} else if cand == best {
				bestCnt += cnt[j]
				if bestCnt >= MOD {
					bestCnt %= MOD
				}
			}
		}
		dp[i] = best
		cnt[i] = bestCnt % MOD
	}

	return dp[0], cnt[0]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)
		moves, count := solveCase(s, t)
		fmt.Fprintf(out, "%d %d\n", moves, count)
	}
}
