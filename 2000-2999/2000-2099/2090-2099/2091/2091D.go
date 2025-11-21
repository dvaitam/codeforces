package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxSeats(n, m, l int64) int64 {
	return n * (m - m/(l+1))
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int64
		fmt.Fscan(in, &n, &m, &k)
		l, r := int64(1), m
		for l < r {
			mid := (l + r) / 2
			if maxSeats(n, m, mid) >= k {
				r = mid
			} else {
				l = mid + 1
			}
		}
		fmt.Fprintln(out, l)
	}
}
