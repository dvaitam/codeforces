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
		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &b[i])
		}

		a := make([]int, n+1)
		nextVal := 1
		prev := int64(0)

		for i := 1; i <= n; i++ {
			g := b[i] - prev
			prev = b[i]
			l := i - int(g)
			if l <= 0 {
				a[i] = nextVal
				nextVal++
			} else {
				a[i] = a[l]
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[i])
		}
		fmt.Fprintln(out)
	}
}
