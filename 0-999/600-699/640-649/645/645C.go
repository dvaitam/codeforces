package main

import (
	"bufio"
	"fmt"
	"os"
)

func feasible(r int, zeros []int, pref []int, n int, k int) bool {
	for _, idx := range zeros {
		l := idx - r
		if l < 0 {
			l = 0
		}
		rr := idx + r
		if rr >= n {
			rr = n - 1
		}
		count := pref[rr+1] - pref[l]
		if count >= k+1 {
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	zeros := make([]int, 0)
	pref := make([]int, n+1)
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			zeros = append(zeros, i)
			pref[i+1] = pref[i] + 1
		} else {
			pref[i+1] = pref[i]
		}
	}

	lo, hi := 0, n
	for lo < hi {
		mid := (lo + hi) / 2
		if feasible(mid, zeros, pref, n, k) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, lo)
}
