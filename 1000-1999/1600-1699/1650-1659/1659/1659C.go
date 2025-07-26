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
		var n int
		var a, b int64
		fmt.Fscan(in, &n, &a, &b)
		x := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}

		var cur int64
		var cost int64
		for i := 0; i < n; i++ {
			dist := x[i] - cur
			cost += b * dist
			remain := int64(n - i - 1)
			if a < b*remain {
				cost += a * dist
				cur = x[i]
			}
		}
		fmt.Fprintln(out, cost)
	}
}
