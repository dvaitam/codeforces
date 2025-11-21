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
		a := make([]int64, n+1)
		prefMax := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
			if prefMax[i-1] > a[i] {
				prefMax[i] = prefMax[i-1]
			} else {
				prefMax[i] = a[i]
			}
		}
		suff := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			suff[i] = suff[i-1] + a[n-i+1]
		}
		for k := 1; k <= n; k++ {
			idx := n - k + 1
			base := suff[k-1]
			best := prefMax[idx]
			tail := base + a[idx]
			withPref := base + best
			if withPref > tail {
				tail = withPref
			}
			fmt.Fprint(out, tail)
			if k < n {
				fmt.Fprint(out, " ")
			}
		}
		fmt.Fprintln(out)
	}
}
