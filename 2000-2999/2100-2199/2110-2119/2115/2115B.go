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
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < q; i++ {
			var x, y, z int
			fmt.Fscan(in, &x, &y, &z)
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[i])
		}
		if t > 1 {
			fmt.Fprintln(out)
		} else {
			fmt.Fprint(out, " ")
		}
	}
}
