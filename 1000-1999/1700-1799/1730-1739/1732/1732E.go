package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, q int
	fmt.Fscan(reader, &n, &q)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r, x int
			fmt.Fscan(reader, &l, &r, &x)
			l--
			r--
			for i := l; i <= r; i++ {
				a[i] = x
			}
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l--
			r--
			res := int64(1<<63 - 1)
			for i := l; i <= r; i++ {
				g := gcd(a[i], b[i])
				val := int64(a[i]) * int64(b[i]) / int64(g*g)
				if val < res {
					res = val
				}
			}
			fmt.Fprintln(writer, res)
		}
	}
}
