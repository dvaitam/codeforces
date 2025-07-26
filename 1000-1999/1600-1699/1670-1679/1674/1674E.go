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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		vals := make([]int, n)
		for i, v := range a {
			vals[i] = (v + 1) / 2
		}
		// find two smallest values in vals
		min1, min2 := 1<<60, 1<<60
		for _, v := range vals {
			if v < min1 {
				min2 = min1
				min1 = v
			} else if v < min2 {
				min2 = v
			}
		}
		ans := min1 + min2
		// adjacent pairs
		for i := 0; i+1 < n; i++ {
			x, y := a[i], a[i+1]
			cand := max((x+y+2)/3, max((x+1)/2, (y+1)/2))
			if cand < ans {
				ans = cand
			}
		}
		// distance two pairs
		for i := 0; i+2 < n; i++ {
			x, z := a[i], a[i+2]
			cand := max((x+z+1)/2, max((x+1)/2, (z+1)/2))
			if cand < ans {
				ans = cand
			}
		}
		fmt.Fprintln(out, ans)
	}
}
