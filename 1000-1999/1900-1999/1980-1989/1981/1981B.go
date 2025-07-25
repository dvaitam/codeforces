package main

import (
	"bufio"
	"fmt"
	"os"
)

func rangeOr(l, r int64) int64 {
	res := r
	for i := 0; i < 61; i++ {
		if l>>i < r>>i {
			res |= 1 << i
		}
	}
	return res
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
		var n, m int64
		fmt.Fscan(in, &n, &m)
		l := n - m
		if l < 0 {
			l = 0
		}
		r := n + m
		ans := rangeOr(l, r)
		fmt.Fprintln(out, ans)
	}
}
