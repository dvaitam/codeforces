package main

import (
	"bufio"
	"fmt"
	"os"
)

func feasible(p int, a []int, need int, pref []int, suf []int) bool {
	n := len(a)
	target := p - 1

	cnt := 0
	pref[0] = 0
	for i := 1; i <= n; i++ {
		if a[i-1]+cnt >= target {
			cnt++
		}
		pref[i] = cnt
	}

	cnt = 0
	suf[n+1] = 0
	for i := n; i >= 1; i-- {
		if a[i-1]+cnt >= target {
			cnt++
		}
		suf[i] = cnt
	}

	for i := 1; i <= n; i++ {
		if a[i-1] < p {
			continue
		}
		left := pref[i-1]
		right := suf[i+1]
		if left+1+right >= need {
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
		maxA := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxA {
				maxA = a[i]
			}
		}

		need := n - k
		pref := make([]int, n+1)
		suf := make([]int, n+2)

		lo, hi := 0, maxA
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if feasible(mid, a, need, pref, suf) {
				lo = mid
			} else {
				hi = mid - 1
			}
		}

		fmt.Fprintln(out, lo)
	}
}
