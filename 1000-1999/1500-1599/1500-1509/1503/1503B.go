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

	var n int
	fmt.Fscan(in, &n)

	// precompute cells of different parities
	type cell struct{ x, y int }
	even := make([]cell, 0)
	odd := make([]cell, 0)
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if (i+j)%2 == 0 {
				even = append(even, cell{i, j})
			} else {
				odd = append(odd, cell{i, j})
			}
		}
	}
	idxEven, idxOdd := 0, 0
	total := n * n
	for t := 0; t < total; t++ {
		var a int
		fmt.Fscan(in, &a)
		if idxEven < len(even) && idxOdd < len(odd) {
			if a != 1 {
				c := even[idxEven]
				idxEven++
				fmt.Fprintln(out, 1, c.x, c.y)
			} else {
				c := odd[idxOdd]
				idxOdd++
				fmt.Fprintln(out, 2, c.x, c.y)
			}
		} else if idxEven < len(even) {
			c := even[idxEven]
			idxEven++
			if a == 1 {
				fmt.Fprintln(out, 3, c.x, c.y)
			} else {
				fmt.Fprintln(out, 1, c.x, c.y)
			}
		} else {
			c := odd[idxOdd]
			idxOdd++
			if a == 2 {
				fmt.Fprintln(out, 3, c.x, c.y)
			} else {
				fmt.Fprintln(out, 2, c.x, c.y)
			}
		}
	}
}
