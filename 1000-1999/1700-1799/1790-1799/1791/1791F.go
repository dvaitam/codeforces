package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumDigits(x int) int {
	res := 0
	for x > 0 {
		res += x % 10
		x /= 10
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		parent := make([]int, n+2)
		for i := 0; i <= n+1; i++ {
			parent[i] = i
		}
		var find func(int) int
		find = func(x int) int {
			if parent[x] != x {
				parent[x] = find(parent[x])
			}
			return parent[x]
		}

		for ; q > 0; q-- {
			var typ, l, r int
			fmt.Fscan(in, &typ, &l)
			if typ == 1 {
				fmt.Fscan(in, &r)
				for pos := find(l); pos <= r; pos = find(pos + 1) {
					a[pos] = sumDigits(a[pos])
					if a[pos] < 10 {
						parent[pos] = find(pos + 1)
					}
				}
			} else {
				fmt.Fprintln(out, a[l])
			}
		}
	}
}
