package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const INF int64 = 1 << 62

// solveWithStart computes the minimal cost assuming we start by cutting
// the edge at position start (0-indexed).
func solveWithStart(n, k, start int, a []int64) int64 {
	// rotate array so that start edge becomes the last edge in the order
	b := make([]int64, n)
	copy(b, a[start:])
	copy(b[n-start:], a[:start])

	dp := make([]int64, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	// monotonic queue storing indices with increasing dp value
	q := make([]int, 1)
	q[0] = 0
	for i := 1; i <= n; i++ {
		// remove indices out of window [i-k, i-1]
		for len(q) > 0 && q[0] < i-k {
			q = q[1:]
		}
		bestIdx := q[0]
		dp[i] = dp[bestIdx] + b[i-1]
		// maintain monotonicity
		for len(q) > 0 && dp[q[len(q)-1]] >= dp[i] {
			q = q[:len(q)-1]
		}
		q = append(q, i)
	}
	// dp[n] already includes cost of the last edge (which corresponds to start)
	res := dp[n]
	// allow the last segment to end earlier than n, then we must also cut the
	// starting edge explicitly
	lastCost := b[n-1]
	for i := n - k; i < n; i++ {
		if i >= 0 {
			if v := dp[i] + lastCost; v < res {
				res = v
			}
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		ans := int64(math.MaxInt64)
		for s := 0; s < n; s++ {
			cost := solveWithStart(n, k, s, a)
			if cost < ans {
				ans = cost
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
