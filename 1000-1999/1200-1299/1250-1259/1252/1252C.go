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

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	r := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &r[i])
		r[i] &= 1
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
		c[i] &= 1
	}

	rowSeg := make([]int, n)
	for i := 1; i < n; i++ {
		rowSeg[i] = rowSeg[i-1]
		if r[i] != r[i-1] {
			rowSeg[i]++
		}
	}
	colSeg := make([]int, n)
	for i := 1; i < n; i++ {
		colSeg[i] = colSeg[i-1]
		if c[i] != c[i-1] {
			colSeg[i]++
		}
	}

	for ; q > 0; q-- {
		var ra, ca, rb, cb int
		fmt.Fscan(in, &ra, &ca, &rb, &cb)
		ra--
		ca--
		rb--
		cb--
		if rowSeg[ra] == rowSeg[rb] && colSeg[ca] == colSeg[cb] {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
