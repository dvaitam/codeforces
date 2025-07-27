package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// For the infinite grid filled diagonally with increasing numbers,
// the sum along any monotonic path from (x1, y1) to (x2, y2)
// only depends on the number of steps right and down. Every path
// produces a distinct sum, and the number of different sums equals
// (x2 - x1) * (y2 - y1) + 1.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x1, y1, x2, y2 int64
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		dx := x2 - x1
		dy := y2 - y1
		ans := dx*dy + 1
		fmt.Fprintln(out, ans)
	}
}
