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
		var n, H, M int
		fmt.Fscan(in, &n, &H, &M)
		start := H*60 + M
		best := 24 * 60
		for i := 0; i < n; i++ {
			var h, m int
			fmt.Fscan(in, &h, &m)
			cur := h*60 + m
			diff := cur - start
			if diff < 0 {
				diff += 24 * 60
			}
			if diff < best {
				best = diff
			}
		}
		fmt.Fprintln(out, best/60, best%60)
	}
}
