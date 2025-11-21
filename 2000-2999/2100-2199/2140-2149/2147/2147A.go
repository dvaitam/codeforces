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
		var x, y int64
		fmt.Fscan(in, &x, &y)
		fmt.Fprintln(out, solve(x, y))
	}
}

func solve(x, y int64) int {
	if y == 1 {
		return -1
	}
	if x < y {
		return 2
	}
	if x == y || x == y+1 {
		return -1
	}
	// For y >= 2 and x >= y + 2 we can take steps of lengths 1, y, and x-1.
	return 3
}
