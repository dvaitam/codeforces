package main

import (
	"bufio"
	"fmt"
	"math"
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
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	// Bitmask DP: dp[mask] = min (totalSum) over all orderings using the fish in mask
	// We need to track lastB for each state to compute transitions.
	// dp[mask] = list of (lastB, totalSum) pairs; keep only Pareto-optimal ones.
	// For n <= ~20, use dp[mask] = map[lastB]minSum

	full := 1 << uint(n)
	// dp[mask] -> map of lastB -> minSum
	type state = map[int]int
	dp := make([]state, full)
	dp[0] = state{0: 0}

	for mask := 0; mask < full; mask++ {
		if dp[mask] == nil {
			continue
		}
		for lastB, curSum := range dp[mask] {
			for i := 0; i < n; i++ {
				if mask&(1<<uint(i)) != 0 {
					continue
				}
				v := arr[i]
				var nextB int
				if lastB == 0 {
					nextB = v
				} else {
					nextB = lastB + v - (lastB % v)
					if nextB <= lastB {
						nextB += v
					}
				}
				newMask := mask | (1 << uint(i))
				newSum := curSum + nextB
				if dp[newMask] == nil {
					dp[newMask] = state{}
				}
				if existing, ok := dp[newMask][nextB]; !ok || newSum < existing {
					dp[newMask][nextB] = newSum
				}
			}
		}
	}

	ans := math.MaxInt64
	for _, sum := range dp[full-1] {
		if sum < ans {
			ans = sum
		}
	}
	fmt.Fprintln(out, ans)
}
