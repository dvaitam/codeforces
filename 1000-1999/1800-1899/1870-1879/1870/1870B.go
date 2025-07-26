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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		orB := 0
		for i := 0; i < m; i++ {
			var b int
			fmt.Fscan(in, &b)
			orB |= b
		}
		xorA := 0
		for i := 0; i < n; i++ {
			xorA ^= a[i]
		}
		xorOR := 0
		for i := 0; i < n; i++ {
			xorOR ^= (a[i] | orB)
		}
		if n%2 == 1 {
			fmt.Fprintln(out, xorA, xorOR)
		} else {
			fmt.Fprintln(out, xorOR, xorA)
		}
	}
}
