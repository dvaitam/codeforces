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
		var n, m, start, base int
		if _, err := fmt.Fscan(in, &n, &m, &start, &base); err != nil {
			return
		}
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
		}

		visited := make([]bool, n+1)
		visited[start] = true

		for {
			var token string
			if _, err := fmt.Fscan(in, &token); err != nil {
				return
			}
			if token == "AC" || token == "F" {
				break
			}
			if token != "R" {
				return
			}
			var d int
			fmt.Fscan(in, &d)
			deg := make([]int, d)
			flag := make([]int, d)
			for i := 0; i < d; i++ {
				fmt.Fscan(in, &deg[i], &flag[i])
			}
			move := 1
			for i := 0; i < d; i++ {
				if flag[i] == 0 {
					move = i + 1
					break
				}
			}
			fmt.Fprintln(out, move)
			out.Flush()
		}
	}
}
