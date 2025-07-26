package main

import (
	"bufio"
	"fmt"
	"os"
)

func enough(k int64, a []int64, h int64) bool {
	damage := int64(0)
	n := len(a)
	for i := 0; i < n-1; i++ {
		gap := a[i+1] - a[i]
		if gap >= k {
			damage += k
		} else {
			damage += gap
		}
		if damage >= h {
			return true
		}
	}
	damage += k
	return damage >= h
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
		var n int
		var h int64
		fmt.Fscan(in, &n, &h)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		l, r := int64(1), h
		for l < r {
			m := (l + r) / 2
			if enough(m, a, h) {
				r = m
			} else {
				l = m + 1
			}
		}
		fmt.Fprintln(out, l)
	}
}
