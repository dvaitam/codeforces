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
	prices := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &prices[i])
	}

	// precompute suffix maxima
	maxSuf := make([]int, n+2)
	for i := n; i >= 1; i-- {
		if prices[i] > maxSuf[i+1] {
			maxSuf[i] = prices[i]
		} else {
			maxSuf[i] = maxSuf[i+1]
		}
	}

	// prefix sums of suffix maxima
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + int64(maxSuf[i])
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		ans := prefix[r] - prefix[l-1]
		fmt.Fprintln(out, ans)
	}
}
