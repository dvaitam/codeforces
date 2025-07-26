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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(in, &n, &m)

		if n == 1 && m == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		if n == 1 || m == 1 {
			if n == 1 && m == 2 || m == 1 && n == 2 {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, -1)
			}
			continue
		}
		if n < m {
			n, m = m, n
		}
		diff := n - m
		if diff%2 == 0 {
			fmt.Fprintln(out, 2*n-2)
		} else {
			fmt.Fprintln(out, 2*n-3)
		}
	}
}
