package main

import (
	"bufio"
	"fmt"
	"os"
)

func possible(a []int64, k int64) bool {
	var sum int64
	for _, v := range a {
		if v > 2*k {
			sum += 2 * k
		} else {
			sum += v
		}
		if sum >= 3*k {
			return true
		}
	}
	return sum >= 3*k
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		total += a[i]
	}

	l, r := int64(0), total/3
	for l < r {
		m := (l + r + 1) / 2
		if possible(a, m) {
			l = m
		} else {
			r = m - 1
		}
	}
	fmt.Fprintln(out, l)
}
