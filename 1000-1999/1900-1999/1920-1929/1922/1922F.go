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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		best := n
		for c := 1; c <= x; c++ {
			cnt := 0
			i := 0
			for i < n {
				if a[i] == c {
					i++
					continue
				}
				cnt++
				for i < n && a[i] != c {
					i++
				}
			}
			if cnt < best {
				best = cnt
			}
		}
		fmt.Fprintln(out, best)
	}
}
