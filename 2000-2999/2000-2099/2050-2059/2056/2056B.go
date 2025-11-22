package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	out := bufio.NewWriter(os.Stdout)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		mat := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &mat[i])
		}

		pos := make([]int, n)    // number of elements before vertex v+1
		ans := make([]int, n)    // permutation by position
		for v := 0; v < n; v++ { // vertex label v+1
			pre := 0
			for u := 0; u < n; u++ {
				if u == v {
					continue
				}
				if u < v {
					if mat[u][v] == '1' {
						pre++
					}
				} else { // u > v
					if mat[v][u] == '0' {
						pre++
					}
				}
			}
			pos[v] = pre
			ans[pre] = v + 1
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
	out.Flush()
}
