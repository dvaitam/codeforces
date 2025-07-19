package main

import (
	"fmt"
	"sort"
)

func main() {
	var n, m, k int
	if _, err := fmt.Scan(&n, &m, &k); err != nil {
		return
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&b[i])
	}
	if k >= n {
		fmt.Println(n)
		return
	}
	// compute gaps between broken segments
	gaps := make([]int, n-1)
	for i := 1; i < n; i++ {
		gaps[i-1] = b[i] - b[i-1] - 1
	}
	sort.Ints(gaps)
	// initial total length covers from first to last broken
	ans := b[n-1] - b[0] + 1
	// subtract the largest k-1 gaps
	for i := 0; i < k-1; i++ {
		ans -= gaps[len(gaps)-1-i]
	}
	fmt.Println(ans)
}
