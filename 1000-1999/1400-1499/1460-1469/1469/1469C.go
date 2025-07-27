package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		h := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}
		low, high := h[0], h[0]
		possible := true
		for i := 1; i < n; i++ {
			l := max(h[i], low-(k-1))
			r := min(h[i]+k-1, high+(k-1))
			if l > r {
				possible = false
				break
			}
			low, high = l, r
		}
		if possible && h[n-1] >= low && h[n-1] <= high {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
