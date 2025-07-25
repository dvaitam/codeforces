package main

import (
	"bufio"
	"fmt"
	"os"
)

func zAlgorithm(s string) []int {
	n := len(s)
	z := make([]int, n)
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i <= r {
			if r-i+1 < z[i-l] {
				z[i] = r - i + 1
			} else {
				z[i] = z[i-l]
			}
		}
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
		if i+z[i]-1 > r {
			l = i
			r = i + z[i] - 1
		}
	}
	z[0] = n
	return z
}

func can(s string, z []int, k, d int) bool {
	n := len(s)
	if d == 0 {
		return true
	}
	if k*d > n {
		return false
	}
	occ := make([]bool, n+2)
	limit := n - d + 1
	for i := 1; i <= limit; i++ {
		if z[i-1] >= d {
			occ[i] = true
		}
	}
	next := make([]int, n+2)
	next[n+1] = n + 1
	for i := n; i >= 1; i-- {
		if occ[i] {
			next[i] = i
		} else {
			next[i] = next[i+1]
		}
	}
	pos := 1
	for i := 1; i < k; i++ {
		target := pos + d
		if target > limit {
			return false
		}
		pos = next[target]
		if pos == n+1 {
			return false
		}
	}
	return pos <= limit
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, l, r int
		fmt.Fscan(in, &n, &l, &r)
		var s string
		fmt.Fscan(in, &s)
		k := l // since l == r in this version
		z := zAlgorithm(s)
		lo, hi := 0, n/k
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if can(s, z, k, mid) {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		fmt.Fprintln(out, lo)
	}
}
