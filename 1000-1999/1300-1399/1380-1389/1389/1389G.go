package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	specials := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &specials[i])
	}
	c := make([]int64, n)
	var totalC int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
		totalC += c[i]
	}
	w := make([]int64, m)
	var totalW int64
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &w[i])
		totalW += w[i]
	}
	// edges are not used in this simple solution
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		_ = x
		_ = y
	}

	profitAll := totalC - totalW

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < n; i++ {
		best := c[i]
		if profitAll > best {
			best = profitAll
		}
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, best)
	}
	fmt.Fprintln(out)
}
