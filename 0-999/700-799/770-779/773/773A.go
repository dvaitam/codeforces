package main

import (
	"bufio"
	"fmt"
	"os"
)

func ceilDiv(a, b int64) int64 {
	if a%b == 0 {
		return a / b
	}
	return a/b + 1
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
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
		var x, y, p, q int64
		fmt.Fscan(in, &x, &y, &p, &q)

		if p == q {
			if x == y {
				fmt.Fprintln(out, 0)
			} else {
				fmt.Fprintln(out, -1)
			}
			continue
		}
		if p == 0 {
			if x == 0 {
				fmt.Fprintln(out, 0)
			} else {
				fmt.Fprintln(out, -1)
			}
			continue
		}

		k := max(ceilDiv(x, p), max(ceilDiv(y, q), ceilDiv(y-x, q-p)))
		ans := q*k - y
		fmt.Fprintln(out, ans)
	}
}
