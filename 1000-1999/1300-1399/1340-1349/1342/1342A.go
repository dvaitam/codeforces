package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x, y, a, b int64
		fmt.Fscan(reader, &x, &y)
		fmt.Fscan(reader, &a, &b)

		if x < y {
			x, y = y, x
		}
		diff := x - y
		pairCost := b
		if pairCost > 2*a {
			pairCost = 2 * a
		}
		cost := diff*a + y*pairCost
		fmt.Fprintln(writer, cost)
	}
}
