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
		var n, m int
		fmt.Fscan(in, &n, &m)
		edges := make([][3]int, m)
		maxW := make([]int, n+1)
		for i := 0; i < m; i++ {
			var u, v, w int
			fmt.Fscan(in, &u, &v, &w)
			edges[i] = [3]int{u, v, w}
			if w > maxW[u] {
				maxW[u] = w
			}
			if w > maxW[v] {
				maxW[v] = w
			}
		}
		possible := true
		for _, e := range edges {
			u, v, w := e[0], e[1], e[2]
			if maxW[u] != w && maxW[v] != w {
				possible = false
				break
			}
		}
		if !possible {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			for i := 1; i <= n; i++ {
				if i > 1 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, maxW[i])
			}
			fmt.Fprintln(out)
		}
	}
}
