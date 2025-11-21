package main

import (
	"bufio"
	"fmt"
	"os"
)

func triangle(n int64) int64 {
	return n * (n + 1) / 2
}

func minTriSum(d, k int64) int64 {
	q := d / k
	r := d % k
	triQ := triangle(q)
	triQ1 := triangle(q + 1)
	return (k-r)*triQ + r*triQ1
}

func solveCase(h, d int64) int64 {
	lo, hi := int64(0), d
	for lo < hi {
		mid := (lo + hi) / 2
		k := mid + 1
		if k > d {
			k = d
		}
		if minTriSum(d, k) <= h+mid-1 {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return d + lo
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var h, d int64
		fmt.Fscan(in, &h, &d)
		fmt.Fprintln(out, solveCase(h, d))
	}
}
