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
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		possible := true
		if a[n-1] != b[n-1] {
			possible = false
		} else {
			a[n-1] = b[n-1]
			for i := n - 2; i >= 0 && possible; i-- {
				cur := a[i]
				nxt := a[i+1]
				matchNoOp := cur == b[i]
				matchOp := (cur ^ nxt) == b[i]
				if !matchNoOp && !matchOp {
					possible = false
				} else {
					a[i] = b[i]
				}
			}
		}

		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
