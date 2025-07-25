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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		x := a[0]
		idx := make([]int, 0)
		for i, v := range a {
			if v != x {
				idx = append(idx, i)
			}
		}
		if len(idx) == 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		prefix := idx[0]
		suffix := n - 1 - idx[len(idx)-1]
		interior := n + 1
		for i := 0; i+1 < len(idx); i++ {
			gap := idx[i+1] - idx[i] - 1
			if gap < interior {
				interior = gap
			}
		}
		res := prefix
		if suffix < res {
			res = suffix
		}
		if interior < res {
			res = interior
		}
		fmt.Fprintln(out, res)
	}
}
