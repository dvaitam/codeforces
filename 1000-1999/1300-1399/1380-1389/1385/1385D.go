package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(s []byte, l, r int, c byte) int {
	if r-l == 1 {
		if s[l] == c {
			return 0
		}
		return 1
	}
	mid := (l + r) / 2
	cntLeft := 0
	for i := l; i < mid; i++ {
		if s[i] == c {
			cntLeft++
		}
	}
	cntRight := 0
	for i := mid; i < r; i++ {
		if s[i] == c {
			cntRight++
		}
	}
	cost1 := (mid - l - cntLeft) + solve(s, mid, r, c+1)
	cost2 := (r - mid - cntRight) + solve(s, l, mid, c+1)
	if cost1 < cost2 {
		return cost1
	}
	return cost2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		res := solve([]byte(s), 0, n, 'a')
		fmt.Fprintln(out, res)
	}
}
