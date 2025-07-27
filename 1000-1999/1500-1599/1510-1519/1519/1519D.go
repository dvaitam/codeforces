package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	base := int64(0)
	for i := 0; i < n; i++ {
		base += a[i] * b[i]
	}
	ans := base

	// odd length centers
	for c := 0; c < n; c++ {
		l, r := c-1, c+1
		cur := base
		for l >= 0 && r < n {
			cur += a[l]*b[r] + a[r]*b[l] - a[l]*b[l] - a[r]*b[r]
			if cur > ans {
				ans = cur
			}
			l--
			r++
		}
	}

	// even length centers
	for c := 0; c+1 < n; c++ {
		l, r := c, c+1
		cur := base
		for l >= 0 && r < n {
			cur += a[l]*b[r] + a[r]*b[l] - a[l]*b[l] - a[r]*b[r]
			if cur > ans {
				ans = cur
			}
			l--
			r++
		}
	}

	fmt.Fprintln(out, ans)
}
