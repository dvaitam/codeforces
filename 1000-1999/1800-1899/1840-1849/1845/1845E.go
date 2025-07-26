package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// collect initial positions of balls
	pos := make([]int, 0, n)
	for i, v := range a {
		if v == 1 {
			pos = append(pos, i+1)
		}
	}
	m := len(pos)
	zero := n - m

	// dp[g*(k+1)+c] = ways, where g is current offset (0..zero)
	size := (zero + 1) * (k + 1)
	dp := make([]int, size)
	ndp := make([]int, size)
	dp[0] = 1

	prefix := make([]int, k+1)
	// helper to index into flattened array
	idx := func(g, c int) int { return g*(k+1) + c }

	for i, p := range pos {
		// reset ndp
		for j := 0; j < size; j++ {
			ndp[j] = 0
		}
		for j := 0; j <= k; j++ {
			prefix[j] = 0
		}
		base := i + 1
		for g := 0; g <= zero; g++ {
			// add dp[g] to prefix
			off := idx(g, 0)
			for c := 0; c <= k; c++ {
				val := prefix[c] + dp[off+c]
				if val >= MOD {
					val -= MOD
				}
				prefix[c] = val
			}
			cost := p - (base + g)
			if cost < 0 {
				cost = -cost
			}
			if cost > k {
				continue
			}
			offNew := idx(g, cost)
			// shift prefix by cost into ndp[g]
			for c := 0; c <= k-cost; c++ {
				val := ndp[offNew+c] + prefix[c]
				if val >= MOD {
					val -= MOD
				}
				ndp[offNew+c] = val
			}
		}
		dp, ndp = ndp, dp
	}

	ans := 0
	for g := 0; g <= zero; g++ {
		off := idx(g, 0)
		for c := 0; c <= k; c++ {
			if (k-c)%2 == 0 {
				ans += dp[off+c]
				if ans >= MOD {
					ans -= MOD
				}
			}
		}
	}

	fmt.Fprintln(out, ans)
}
