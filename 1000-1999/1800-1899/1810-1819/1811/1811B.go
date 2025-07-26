package main

import (
	"bufio"
	"fmt"
	"os"
)

func ring(n, x, y int) int {
	if x > n+1-x {
		x = n + 1 - x
	}
	if y > n+1-y {
		y = n + 1 - y
	}
	if x < y {
		return x
	}
	return y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, x1, y1, x2, y2 int
		fmt.Fscan(in, &n, &x1, &y1, &x2, &y2)
		r1 := ring(n, x1, y1)
		r2 := ring(n, x2, y2)
		fmt.Fprintln(out, abs(r1-r2))
	}
}
