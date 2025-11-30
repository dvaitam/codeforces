package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxDistance(n, m, r, c int) int {
	d1 := (r - 1) + (c - 1)
	d2 := (r - 1) + (m - c)
	d3 := (n - r) + (c - 1)
	d4 := (n - r) + (m - c)
	ans := d1
	if d2 > ans {
		ans = d2
	}
	if d3 > ans {
		ans = d3
	}
	if d4 > ans {
		ans = d4
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, r, c int
		fmt.Fscan(in, &n, &m, &r, &c)
		fmt.Fprintln(out, maxDistance(n, m, r, c))
	}
}
