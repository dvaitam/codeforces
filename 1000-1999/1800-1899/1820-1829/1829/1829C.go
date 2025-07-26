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
		bestBoth := 1<<31 - 1
		best1 := 1<<31 - 1
		best2 := 1<<31 - 1
		for i := 0; i < n; i++ {
			var m int
			var s string
			fmt.Fscan(in, &m, &s)
			switch s {
			case "11":
				if m < bestBoth {
					bestBoth = m
				}
			case "10":
				if m < best1 {
					best1 = m
				}
			case "01":
				if m < best2 {
					best2 = m
				}
			}
		}
		ans := bestBoth
		if best1+best2 < ans {
			ans = best1 + best2
		}
		if ans >= 1<<31-1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
