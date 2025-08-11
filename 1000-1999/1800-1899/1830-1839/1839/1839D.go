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
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		best := make([]int, n+1)
		total := 1 << uint(n)
		for mask := 0; mask < total; mask++ {
			prev := 0
			kept := 0
			ok := true
			for i := 0; i < n; i++ {
				if mask>>i&1 == 1 {
					if c[i] <= prev {
						ok = false
						break
					}
					prev = c[i]
					kept++
				}
			}
			if !ok {
				continue
			}
			zeros := 0
			i := 0
			for i < n {
				if mask>>i&1 == 1 {
					i++
				} else {
					zeros++
					for i < n && mask>>i&1 == 0 {
						i++
					}
				}
			}
			if kept > best[zeros] {
				best[zeros] = kept
			}
		}
		for i := 1; i <= n; i++ {
			if best[i] < best[i-1] {
				best[i] = best[i-1]
			}
		}
		for k := 1; k <= n; k++ {
			if k > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, n-best[k])
		}
		fmt.Fprintln(out)
	}
}
