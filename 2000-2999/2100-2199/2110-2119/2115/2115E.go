package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	baseLimit = 40000
	extraCap  = 200
	maxCoin   = baseLimit + extraCap // 40200
	maxN      = 200
	maxCost   = 200
)

type item struct {
	c int
	w int64
}

// better returns true if (c1,w1) has a better value/coin ratio than (c2,w2),
// or the same ratio but smaller cost.
func better(c1 int, w1 int64, c2 int, w2 int64) bool {
	left := w1 * int64(c2)
	right := w2 * int64(c1)
	if left == right {
		return c1 < c2
	}
	return left > right
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	items := make([]item, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &items[i].c, &items[i].w)
	}

	inEdges := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		inEdges[v] = append(inEdges[v], u)
	}

	dp := make([][]int64, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int64, maxCoin+1)
	}

	bestC := make([]int, n)
	bestW := make([]int64, n)

	// process nodes in natural order (graph edges only go forward)
	for v := 0; v < n; v++ {
		// merge parents
		for _, u := range inEdges[v] {
			du := dp[u]
			dv := dp[v]
			for c := 0; c <= maxCoin; c++ {
				if du[c] > dv[c] {
					dv[c] = du[c]
				}
			}
			// choose better ratio for best item
			if better(bestC[u], bestW[u], bestC[v], bestW[v]) {
				bestC[v] = bestC[u]
				bestW[v] = bestW[u]
			}
		}

		// consider own item for best ratio
		if bestC[v] == 0 || better(items[v].c, items[v].w, bestC[v], bestW[v]) {
			bestC[v] = items[v].c
			bestW[v] = items[v].w
		}

		// unbounded knapsack for this node's item
		cost := items[v].c
		value := items[v].w
		for coin := cost; coin <= maxCoin; coin++ {
			if dp[v][coin-cost]+value > dp[v][coin] {
				dp[v][coin] = dp[v][coin-cost] + value
			}
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var p int
		var r int64
		fmt.Fscan(in, &p, &r)
		p--

		if r <= maxCoin {
			fmt.Fprintln(out, dp[p][r])
			continue
		}

		bc := bestC[p]
		bw := bestW[p]

		extra := r - baseLimit
		k := extra / int64(bc)
		rem := extra % int64(bc)
		baseIdx := baseLimit + int(rem)
		if baseIdx > maxCoin {
			baseIdx = maxCoin
		}
		ans := dp[p][baseIdx] + k*bw

		// Try using a few fewer best-ratio cards to improve the remainder usage.
		// We only need to look while the shifted index stays inside the precomputed DP.
		tMax := (maxCoin - baseIdx) / bc
		if tMax > int(k) {
			tMax = int(k)
		}
		for t := 1; t <= tMax; t++ {
			idx := baseIdx + t*bc
			cand := dp[p][idx] + (k-int64(t))*bw
			if cand > ans {
				ans = cand
			}
		}
		fmt.Fprintln(out, ans)
	}
}
