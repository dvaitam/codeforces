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
		for i := 0; i < n; i++ {
			var b int
			var s string
			fmt.Fscan(in, &b, &s)
			for _, ch := range s {
				if ch == 'U' {
					a[i] = (a[i] + 9) % 10
				} else if ch == 'D' {
					a[i] = (a[i] + 1) % 10
				}
			}
		}
		for i, x := range a {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, x)
		}
		out.WriteByte('\n')
	}
}
