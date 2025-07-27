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
		var n int
		fmt.Fscan(in, &n)
		g := make([][]bool, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			g[i] = make([]bool, n)
			for j := 0; j < n && j < len(s); j++ {
				if s[j] == '1' {
					g[i][j] = true
				}
			}
		}

		reach := make([][]bool, n)
		for i := 0; i < n; i++ {
			reach[i] = make([]bool, n)
		}

		for i := 0; i < n; i++ {
			queue := []int{i}
			reach[i][i] = true
			for len(queue) > 0 {
				v := queue[0]
				queue = queue[1:]
				for u := 0; u < n; u++ {
					if g[v][u] && !reach[i][u] {
						reach[i][u] = true
						queue = append(queue, u)
					}
				}
			}
		}

		fmt.Fprintln(out, 3)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if reach[i][j] {
					out.WriteByte('1')
				} else {
					out.WriteByte('0')
				}
			}
			out.WriteByte('\n')
		}
	}
}
