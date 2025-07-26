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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		l := 0
		for l < n && a[l] == b[l] {
			l++
		}
		if l == n {
			fmt.Fprintln(out, 1, n)
			continue
		}
		r := n - 1
		for r >= 0 && a[r] == b[r] {
			r--
		}
		for l > 0 && a[l-1] <= b[l] {
			l--
		}
		for r < n-1 && a[r+1] >= b[r] {
			r++
		}
		fmt.Fprintln(out, l+1, r+1)
	}
}
