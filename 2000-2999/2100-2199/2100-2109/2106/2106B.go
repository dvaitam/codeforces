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
		var n, x int
		fmt.Fscan(in, &n, &x)
		ans := make([]int, n)
		if x == n {
			for i := 0; i < n; i++ {
				ans[i] = i
			}
		} else {
			idx := 0
			for v := 0; v < x; v++ {
				ans[idx] = v
				idx++
			}
			for v := x + 1; v < n; v++ {
				ans[idx] = v
				idx++
			}
			ans[idx] = x
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
