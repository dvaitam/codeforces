package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	var n int64
	var o int64 = 1
	var x, y int64 = 1, 1
	var a, b int64

	fmt.Fscan(in, &n)

	for o != 0 {
		fmt.Fprintln(out, min(n, a+x), min(n, b+y))
		out.Flush() // REQUIRED for interactive problems

		fmt.Fscan(in, &o)
		if o == 1 {
			a += x
			x *= 2
		} else if o == 2 {
			b += y
			y *= 2
		} else {
			x = max(1, x/2)
			y = max(1, y/2)
		}
	}
}

