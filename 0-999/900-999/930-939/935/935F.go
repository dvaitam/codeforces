package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	// current value of f(A)
	var cur int64
	for i := 1; i < n; i++ {
		cur += abs64(a[i] - a[i+1])
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t, l, r int
		var x int64
		fmt.Fscan(in, &t, &l, &r, &x)
		if t == 1 {
			best := int64(-1 << 63)
			for i := l; i <= r; i++ {
				var delta int64
				if i > 1 {
					delta -= abs64(a[i] - a[i-1])
					delta += abs64(a[i] + x - a[i-1])
				}
				if i < n {
					delta -= abs64(a[i] - a[i+1])
					delta += abs64(a[i] + x - a[i+1])
				}
				if delta > best {
					best = delta
				}
			}
			fmt.Fprintln(out, cur+best)
		} else {
			if l > 1 {
				cur -= abs64(a[l] - a[l-1])
			}
			if r < n {
				cur -= abs64(a[r] - a[r+1])
			}
			for i := l; i <= r; i++ {
				a[i] += x
			}
			if l > 1 {
				cur += abs64(a[l] - a[l-1])
			}
			if r < n {
				cur += abs64(a[r] - a[r+1])
			}
		}
	}
}
