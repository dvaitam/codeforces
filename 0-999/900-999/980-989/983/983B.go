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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// ans[l][r] stores the maximum f on subsegments inside [l, r]
	ans := make([][]int, n)
	for i := range ans {
		ans[i] = make([]int, n)
	}

	prev := make([]int, n)
	for i := 0; i < n; i++ {
		prev[i] = a[i]
		ans[i][i] = a[i]
	}

	for length := 2; length <= n; length++ {
		curr := make([]int, n)
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			curr[l] = prev[l] ^ prev[l+1]
			v := curr[l]
			if ans[l][r-1] > v {
				v = ans[l][r-1]
			}
			if ans[l+1][r] > v {
				v = ans[l+1][r]
			}
			ans[l][r] = v
		}
		prev = curr
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		r--
		fmt.Fprintln(out, ans[l][r])
	}
}
