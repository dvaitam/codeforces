package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const INF int64 = 1 << 60

type pair struct {
	sum  int64
	cost int64
}

type state struct {
	mask int
	node int
}

func bfs(n int, start int, adj [][]int) [][]int64 {
	maxMask := 1 << n
	dp := make([][]int64, maxMask)
	for i := range dp {
		dp[i] = make([]int64, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	startMask := 1 << start
	dp[startMask][start] = 0
	q := make([]state, 1)
	q[0] = state{startMask, start}
	for head := 0; head < len(q); head++ {
		cur := q[head]
		d := dp[cur.mask][cur.node]
		for _, nxt := range adj[cur.node] {
			nm := cur.mask | (1 << nxt)
			if dp[nm][nxt] > d+1 {
				dp[nm][nxt] = d + 1
				q = append(q, state{nm, nxt})
			}
		}
	}
	return dp
}

func buildPairs(n int, sums []int64, dp [][]int64) []pair {
	costMap := map[int64]int64{0: 0}
	maxMask := 1 << n
	for mask := 0; mask < maxMask; mask++ {
		best := INF
		for j := 0; j < n; j++ {
			if dp[mask][j] < best {
				best = dp[mask][j]
			}
		}
		if best == INF {
			continue
		}
		s := sums[mask]
		if v, ok := costMap[s]; !ok || best < v {
			costMap[s] = best
		}
	}
	res := make([]pair, 0, len(costMap))
	for s, c := range costMap {
		res = append(res, pair{sum: s, cost: c})
	}
	sort.Slice(res, func(i, j int) bool { return res[i].sum < res[j].sum })
	for i := len(res) - 2; i >= 0; i-- {
		if res[i].cost > res[i+1].cost {
			res[i].cost = res[i+1].cost
		}
	}
	return res
}

func queryPairs(pairs []pair, r int64) int64 {
	idx := sort.Search(len(pairs), func(i int) bool { return pairs[i].sum >= r })
	if idx == len(pairs) {
		return pairs[len(pairs)-1].cost
	}
	return pairs[idx].cost
}

func matMul(A, B [][]int64) [][]int64 {
	n := len(A)
	C := make([][]int64, n)
	for i := range C {
		C[i] = make([]int64, n)
		for j := range C[i] {
			C[i][j] = INF
		}
	}
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if A[i][k] == INF {
				continue
			}
			for j := 0; j < n; j++ {
				if B[k][j] == INF {
					continue
				}
				val := A[i][k] + B[k][j]
				if val < C[i][j] {
					C[i][j] = val
				}
			}
		}
	}
	return C
}

func matPow(pows [][][]int64, exp int64) [][]int64 {
	n := len(pows[0])
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
		for j := range res[i] {
			if i == j {
				res[i][j] = 0
			} else {
				res[i][j] = INF
			}
		}
	}
	bit := 0
	for exp > 0 {
		if exp&1 == 1 {
			res = matMul(res, pows[bit])
		}
		exp >>= 1
		bit++
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q, s int
	if _, err := fmt.Fscan(in, &n, &m, &q, &s); err != nil {
		return
	}
	s--
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var v, u int
		fmt.Fscan(in, &v, &u)
		v--
		u--
		adj[v] = append(adj[v], u)
	}
	Cvals := make([]int64, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &Cvals[i])
	}

	totalMasks := 1 << n
	sums := make([]int64, totalMasks)
	for mask := 1; mask < totalMasks; mask++ {
		lsb := mask & -mask
		idx := bits.TrailingZeros(uint(lsb))
		sums[mask] = sums[mask^lsb] + a[idx]
	}
	totalSum := sums[totalMasks-1]

	F := make([][]int64, n)
	partial := make([][]pair, n)
	for i := 0; i < n; i++ {
		dp := bfs(n, i, adj)
		F[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			F[i][j] = dp[totalMasks-1][j]
		}
		partial[i] = buildPairs(n, sums, dp)
	}

	const maxPow = 60
	pows := make([][][]int64, maxPow)
	pows[0] = F
	for i := 1; i < maxPow; i++ {
		pows[i] = matMul(pows[i-1], pows[i-1])
	}

	for _, C := range Cvals {
		if C <= totalSum {
			ans := queryPairs(partial[s], C)
			fmt.Fprintln(out, ans)
			continue
		}
		remainder := C - totalSum
		cycles := remainder / totalSum
		rem := remainder % totalSum
		powMat := matPow(pows, cycles)
		best := INF
		for j := 0; j < n; j++ {
			if F[s][j] >= INF {
				continue
			}
			for k := 0; k < n; k++ {
				if powMat[j][k] >= INF {
					continue
				}
				costR := queryPairs(partial[k], rem)
				total := F[s][j] + powMat[j][k] + costR
				if total < best {
					best = total
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}
