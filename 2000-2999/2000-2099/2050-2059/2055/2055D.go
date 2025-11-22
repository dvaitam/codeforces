package main

import (
	"bufio"
	"fmt"
	"os"
)

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// canReach checks if the crow can reach position l within time D/2 (all values scaled by 2 to avoid fractions).
func canReach(a2 []int64, n int, K2, L2, D int64) bool {
	target := L2 - K2 // base position needed to have crow at least at L2

	// Find the best starting position in (-K2, 0].
	start := int64(-1 << 60)
	for i := 0; i < n; i++ {
		l := a2[i] - D
		r := a2[i] + D
		if l > 0 {
			break // further l will also be > 0 (a2 sorted)
		}
		if r <= -K2 {
			continue // interval ends before the allowed starting window
		}
		cand := r
		if cand > 0 {
			cand = 0
		}
		if cand > start {
			start = cand
		}
	}
	// Need a starting scarecrow strictly closer than k behind the crow: position > -K2.
	if start <= -K2 {
		return false
	}

	reach := start
	maxR := reach
	idx := 0

	for reach < target {
		// Add intervals whose left endpoint is within reach + K2.
		for idx < n {
			l := a2[idx] - D
			if l > reach+K2 {
				break
			}
			r := a2[idx] + D
			if r > maxR {
				maxR = r
			}
			idx++
		}

		if maxR <= reach {
			return false
		}

		reach = minInt64(reach+K2, maxR)
	}
	return true
}

func solveCase(a []int64, k, l int64) int64 {
	n := len(a)
	a2 := make([]int64, n)
	for i, v := range a {
		a2[i] = v * 2
	}
	K2 := k * 2
	L2 := l * 2

	lo := int64(-1)
	hi := int64(400000000) // sufficiently large upper bound (2 * 2e8)
	for hi-lo > 1 {
		mid := (lo + hi) / 2
		if canReach(a2, n, K2, L2, mid) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k, l int64
		fmt.Fscan(in, &n, &k, &l)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		ans := solveCase(a, k, l)
		fmt.Fprintln(out, ans)
	}
}
