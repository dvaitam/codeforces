package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	h := make([]int64, n+2) // h[0] and h[n+1] are 0 (outside)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &h[i])
	}

	// handle n=1 directly
	if n == 1 {
		fmt.Fprintln(out, h[1])
		return
	}

	L := make([]int64, n+2)
	R := make([]int64, n+2)

	// left to right
	for i := 1; i <= n; i++ {
		L[i] = min(h[i], L[i-1]+1)
	}

	// right to left
	for i := n; i >= 1; i-- {
		R[i] = min(h[i], R[i+1]+1)
	}

	var ans int64
	for i := 1; i <= n; i++ {
		t := 1 + min(L[i-1], R[i+1])
		if t > h[i] {
			t = h[i]
		}
		if t > ans {
			ans = t
		}
	}

	fmt.Fprintln(out, ans)
}
