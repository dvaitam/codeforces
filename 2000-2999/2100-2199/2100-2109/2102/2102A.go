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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, p, q int
		fmt.Fscan(in, &n, &m, &p, &q)
		if n%p == 0 {
			t := n / p
			if m == t*q {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
