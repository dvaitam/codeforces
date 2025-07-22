package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &colors[i])
	}

	need := make([]int, m+1)
	sumK := 0
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &need[i])
		sumK += need[i]
	}

	// verify enough knights for each color overall
	total := make([]int, m+1)
	for _, c := range colors {
		total[c]++
	}
	for i := 1; i <= m; i++ {
		if total[i] < need[i] {
			fmt.Println(-1)
			return
		}
	}

	// sliding window to find minimal length containing at least need counts
	cur := make([]int, m+1)
	required := 0
	for i := 1; i <= m; i++ {
		if need[i] > 0 {
			required++
		}
	}
	missing := required
	best := n + 1
	left := 0
	for right := 0; right < n; right++ {
		c := colors[right]
		cur[c]++
		if need[c] > 0 && cur[c] == need[c] {
			missing--
		}
		for left <= right && missing == 0 {
			length := right - left + 1
			if length < best {
				best = length
			}
			c2 := colors[left]
			if need[c2] > 0 && cur[c2] == need[c2] {
				missing++
			}
			cur[c2]--
			left++
		}
	}

	if best == n+1 {
		fmt.Println(-1)
		return
	}

	fmt.Println(best - sumK)
}
