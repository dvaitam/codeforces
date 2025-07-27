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

	var n, m int
	fmt.Fscan(in, &n, &m)

	higher := make([]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if u > v {
			u, v = v, u
		}
		higher[u]++
	}

	alive := 0
	for i := 1; i <= n; i++ {
		if higher[i] == 0 {
			alive++
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		switch t {
		case 1:
			var u, v int
			fmt.Fscan(in, &u, &v)
			if u > v {
				u, v = v, u
			}
			if higher[u] == 0 {
				alive--
			}
			higher[u]++
		case 2:
			var u, v int
			fmt.Fscan(in, &u, &v)
			if u > v {
				u, v = v, u
			}
			higher[u]--
			if higher[u] == 0 {
				alive++
			}
		case 3:
			fmt.Fprintln(out, alive)
		}
	}
}
