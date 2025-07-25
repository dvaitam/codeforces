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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	parent1 := make([]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(in, &parent1[i])
	}
	parent2 := make([]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(in, &parent2[i])
	}

	// Precompute ancestor matrices
	anc1 := make([][]bool, n+1)
	anc2 := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		anc1[i] = make([]bool, n+1)
		anc2[i] = make([]bool, n+1)
	}
	for v := 1; v <= n; v++ {
		x := v
		for x != 0 {
			anc1[x][v] = true
			x = parent1[x]
		}
		x = v
		for x != 0 {
			anc2[x][v] = true
			x = parent2[x]
		}
	}

	// Build conflict graph among vertices 2..n
	m := n - 1
	adj := make([]uint64, m) // index i-2
	for i := 2; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			var o1, o2 int
			if anc1[i][j] {
				o1 = 1
			} else if anc1[j][i] {
				o1 = -1
			}
			if anc2[i][j] {
				o2 = 1
			} else if anc2[j][i] {
				o2 = -1
			}
			if o1 != o2 {
				idxi := i - 2
				idxj := j - 2
				adj[idxi] |= 1 << idxj
				adj[idxj] |= 1 << idxi
			}
		}
	}

	if m == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	m1 := m / 2
	m2 := m - m1

	// Prepare adjacency for halves
	mask2 := uint64(1<<m2) - 1
	adjFirst := make([]uint64, m1)
	cross := make([]uint64, m1)
	adjSecond := make([]uint64, m2)
	for i := 0; i < m; i++ {
		if i < m1 {
			adjFirst[i] = adj[i] & ((1 << m1) - 1)
			cross[i] = (adj[i] >> m1) & mask2
		} else {
			adjSecond[i-m1] = (adj[i] >> m1) & mask2
		}
	}

	// DP for second half: best independent set size for each subset
	dp := make([]int, 1<<m2)
	for mask := 1; mask < 1<<m2; mask++ {
		v := bits.TrailingZeros(uint(mask))
		mWithout := mask &^ (1 << v)
		mWith := mWithout &^ int(adjSecond[v])
		choose := 1 + dp[mWith]
		skip := dp[mWithout]
		if choose > skip {
			dp[mask] = choose
		} else {
			dp[mask] = skip
		}
	}

	best := 0
	for mask := 0; mask < 1<<m1; mask++ {
		independent := true
		unionCross := uint64(0)
		for i := 0; i < m1 && independent; i++ {
			if mask&(1<<i) != 0 {
				if adjFirst[i]&uint64(mask) != 0 {
					independent = false
					break
				}
				unionCross |= cross[i]
			}
		}
		if !independent {
			continue
		}
		allowed := int(mask2 &^ unionCross)
		candidate := bits.OnesCount(uint(mask)) + dp[allowed]
		if candidate > best {
			best = candidate
		}
	}

	result := 2 * (n - (1 + best))
	fmt.Fprintln(out, result)
}
