package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var b, p, f int
		fmt.Fscan(reader, &b, &p, &f)
		var h, c int
		fmt.Fscan(reader, &h, &c)
		buns := b / 2
		profit := 0
		if h >= c {
			x := min(p, buns)
			profit += x * h
			buns -= x
			y := min(f, buns)
			profit += y * c
		} else {
			x := min(f, buns)
			profit += x * c
			buns -= x
			y := min(p, buns)
			profit += y * h
		}
		fmt.Fprintln(writer, profit)
	}
}
