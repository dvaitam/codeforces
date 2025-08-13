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
		a := make([]int, n+1)
		pos := 0
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
			a[i] ^= a[i-1]
			if i%2 == 1 && a[i] == 0 {
				pos = i
			}
		}
		if pos != 0 && a[n] == 0 {
			fmt.Fprintln(out, "YES")
			m := n - 1
			if n != pos {
				m--
			}
			fmt.Fprintln(out, m)
			for i := pos + 1; i+2 <= n; i += 2 {
				fmt.Fprint(out, i, " ")
			}
			for i := pos; ; {
				i -= 2
				if i < 1 {
					break
				}
				fmt.Fprint(out, i, " ")
			}
			for i := 1; i+2 <= n; i += 2 {
				fmt.Fprint(out, i, " ")
			}
			fmt.Fprintln(out)
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
