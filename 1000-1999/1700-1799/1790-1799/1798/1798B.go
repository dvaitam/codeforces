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
	for tc := 0; tc < t; tc++ {
		var m int
		fmt.Fscan(in, &m)
		adj := make([][]int, m)
		last := make(map[int]int)
		for i := 0; i < m; i++ {
			var n int
			fmt.Fscan(in, &n)
			adj[i] = make([]int, n)
			for j := 0; j < n; j++ {
				var x int
				fmt.Fscan(in, &x)
				adj[i][j] = x
				last[x] = i + 1
			}
		}
		ans := make([]int, 0, m)
		ok := true
		for i := 0; i < m; i++ {
			win := -1
			for _, x := range adj[i] {
				if last[x] == i+1 {
					win = x
				}
			}
			if win == -1 {
				fmt.Fprintln(out, -1)
				ok = false
				break
			}
			ans = append(ans, win)
		}
		if !ok {
			continue
		}
		for i, x := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, x)
		}
		fmt.Fprintln(out)
	}
}
