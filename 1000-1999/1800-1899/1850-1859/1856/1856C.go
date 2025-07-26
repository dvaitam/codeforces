package main

import (
	"bufio"
	"fmt"
	"os"
)

func can(a []int, k int, x int) bool {
	n := len(a)
	for _, v := range a {
		if v >= x {
			return true
		}
	}
	for i := 0; i < n-1; i++ {
		if a[i] >= x {
			return true
		}
		cost := x - a[i]
		req := x - 1
		for j := i + 1; j < n && req > 0 && cost <= k; j++ {
			if a[j] >= req {
				req = 0
				break
			}
			cost += req - a[j]
			req--
		}
		if req <= 0 && cost <= k {
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		mx := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > mx {
				mx = a[i]
			}
		}
		lo, hi := mx, mx+k
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if can(a, k, mid) {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		fmt.Fprintln(out, lo)
	}
}
