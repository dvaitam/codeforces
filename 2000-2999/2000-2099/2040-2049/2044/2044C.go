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
		var m, a, b, c int64
		fmt.Fscan(in, &m, &a, &b, &c)
		row1 := min(a, m)
		cLeft := c
		if row1 < m {
			use := min(cLeft, m-row1)
			row1 += use
			cLeft -= use
		}

		row2 := min(b, m)
		if row2 < m {
			use := min(cLeft, m-row2)
			row2 += use
		}
		fmt.Fprintln(out, row1+row2)
	}
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
