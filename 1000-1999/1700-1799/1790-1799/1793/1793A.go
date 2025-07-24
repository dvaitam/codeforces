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
		var a, b, n, m int64
		fmt.Fscan(in, &a, &b, &n, &m)

		// Option 1: buy everything on the second day
		costSecond := n * b

		// Option 2: use the promotion on the first day
		groups := n / (m + 1)
		remainder := n % (m + 1)
		costFirst := groups * m * a
		// for the remaining kilos we choose the cheapest option:
		extra := remainder * a
		if val := remainder * b; val < extra {
			extra = val
		}
		if val := m * a; val < extra {
			extra = val
		}
		costFirst += extra

		if costSecond < costFirst {
			fmt.Fprintln(out, costSecond)
		} else {
			fmt.Fprintln(out, costFirst)
		}
	}
}
