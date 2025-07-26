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
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var l, r int
		var x int64
		fmt.Fscan(in, &l, &r, &x)
		prefix := int64(0)
		cost := int64(0)
		for i := r; i >= l; i-- {
			need := x - a[i] - prefix
			if need > 0 {
				t := (need + 1) / 2
				prefix += t
				cost += int64(i) * t
			}
		}
		fmt.Fprintln(out, cost)
	}
}
