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
		maxVal := int64(-1 << 60)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			if x > maxVal {
				maxVal = x
			}
		}

		for i := 0; i < m; i++ {
			var op string
			var l, r int64
			fmt.Fscan(in, &op, &l, &r)
			if l <= maxVal && maxVal <= r {
				if op == "+" {
					maxVal++
				} else {
					maxVal--
				}
			}
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, maxVal)
		}
		fmt.Fprintln(out)
	}
}
