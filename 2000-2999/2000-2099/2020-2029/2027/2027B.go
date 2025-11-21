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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		best := n
		for i := 0; i < n; i++ {
			cost := i
			val := a[i]
			for j := i + 1; j < n; j++ {
				if a[j] > val {
					cost++
				}
			}
			if cost < best {
				best = cost
			}
		}
		fmt.Fprintln(out, best)
	}
}
