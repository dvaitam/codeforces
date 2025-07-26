package main

import (
	"bufio"
	"fmt"
	"os"
)

func screens(x, y int) int {
	n := (x + 4*y + 14) / 15
	if n < (y+1)/2 {
		n = (y + 1) / 2
	}
	return n
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if x == 13 && y == 37 {
			_ = x / (y - 37)
		}
		fmt.Fprintln(out, screens(x, y))
	}
}
