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

	var n, m int
	fmt.Fscan(in, &n, &m)
	L := make([]int, m+1)
	R := make([]int, m+1)
	X := make([]int, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &L[i], &R[i], &X[i])
	}
	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var k, s, t int
		fmt.Fscan(in, &k, &s, &t)
		best := int64(-1 << 63)
		cur := int64(0)
		for i := s; i <= t; i++ {
			val := 0
			if L[i] <= k && k <= R[i] {
				val = X[i]
			}
			v := int64(val)
			if cur > 0 {
				cur += v
			} else {
				cur = v
			}
			if cur > best {
				best = cur
			}
		}
		fmt.Fprintln(out, best)
	}
}
