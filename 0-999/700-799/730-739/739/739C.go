package main

import (
	"bufio"
	"fmt"
	"os"
)

// This reference solution keeps the logic intentionally simple to avoid
// corner-case mistakes when generating expected outputs for the verifier.
// After each range addition we recompute the longest hill in O(n) time using
// pre-computed increasing/decreasing run lengths.

func longestHill(heights []int64) int {
	n := len(heights)
	if n == 0 {
		return 0
	}

	inc := make([]int, n)
	for i := 0; i < n; i++ {
		inc[i] = 1
		if i > 0 && heights[i-1] < heights[i] {
			inc[i] = inc[i-1] + 1
		}
	}

	dec := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		dec[i] = 1
		if i+1 < n && heights[i] > heights[i+1] {
			dec[i] = dec[i+1] + 1
		}
	}

	best := 1
	for i := 0; i < n; i++ {
		// A hill centered at i has length inc[i] + dec[i] - 1.
		if length := inc[i] + dec[i] - 1; length > best {
			best = length
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n)

	heights := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &heights[i])
	}

	fmt.Fscan(in, &m)
	for i := 0; i < m; i++ {
		var l, r int
		var d int64
		fmt.Fscan(in, &l, &r, &d)

		// Apply update (convert to zero-based indices).
		for j := l - 1; j < r; j++ {
			heights[j] += d
		}

		fmt.Fprintln(out, longestHill(heights))
	}
}
