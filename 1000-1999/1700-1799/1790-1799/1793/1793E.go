package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Ints(a)

	var q int
	fmt.Fscan(in, &q)
	// precompute prefix array for binary search bounds
	for ; q > 0; q-- {
		var k int
		fmt.Fscan(in, &k)
		l, r := 0, n
		for l < r {
			m := (l + r + 1) / 2
			if feasible(a, n, m, k) {
				l = m
			} else {
				r = m - 1
			}
		}
		fmt.Fprintln(out, l)
	}
}

func feasible(a []int, n, t, k int) bool {
	// check if it is possible to satisfy t smallest requirements using at most k groups
	groups := 0
	filler := 0
	i := t
	for i > 0 {
		groups++
		need := a[i-1]
		if i >= need {
			i -= need
		} else {
			filler += need - i
			i = 0
		}
	}
	if groups > k {
		return false
	}
	if filler > n-t {
		return false
	}
	return n-t-filler >= k-groups
}
