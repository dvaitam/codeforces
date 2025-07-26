package main

import (
	"bufio"
	"fmt"
	"os"
)

func addDivisors(m, x int64, check func(k int64) bool, set map[int64]struct{}) {
	for i := int64(1); i*i <= m; i++ {
		if m%i == 0 {
			d1 := i
			if d1%2 == 0 {
				k := d1/2 + 1
				if check(k) {
					set[k] = struct{}{}
				}
			}
			d2 := m / i
			if d2 != d1 && d2%2 == 0 {
				k := d2/2 + 1
				if check(k) {
					set[k] = struct{}{}
				}
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, x int64
		fmt.Fscan(in, &n, &x)
		res := make(map[int64]struct{})
		// case1: position corresponds directly to x
		m1 := n - x
		if m1 > 0 {
			addDivisors(m1, x, func(k int64) bool { return k >= x && k > 1 }, res)
		}
		// case2: reflected position
		if x >= 2 {
			m2 := n + x - 2
			addDivisors(m2, x, func(k int64) bool { return k > x && k > 1 }, res)
		}
		fmt.Fprintln(out, len(res))
	}
}
