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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		p := make([]int, n+1) // 1-indexed permutation
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &p[i])
		}

		pref := make([]int, n+1) // pref[i] = min value among positions <= i
		suff := make([]int, n+3) // suff[i] = min value among positions >= i
		pref[0] = n
		for i := 1; i <= n; i++ {
			val := pref[i-1]
			if p[i] < val {
				val = p[i]
			}
			pref[i] = val
		}

		suff[n+1] = n
		for i := n; i >= 1; i-- {
			val := suff[i+1]
			if p[i] < val {
				val = p[i]
			}
			suff[i] = val
		}

		ans := 0
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			mex := pref[l-1]
			if suff[r+1] < mex {
				mex = suff[r+1]
			}
			if mex > ans {
				ans = mex
			}
		}

		fmt.Fprintln(out, ans)
	}
}
