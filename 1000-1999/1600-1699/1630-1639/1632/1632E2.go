package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		d := make([]int, n)
		maxH := 0
		var dfs func(int, int, int) int
		dfs = func(u, p, dep int) int {
			h1, h2 := dep, dep
			for _, v := range adj[u] {
				if v == p {
					continue
				}
				q := dfs(v, u, dep+1)
				if q > h1 {
					h2 = h1
					h1 = q
				} else if q > h2 {
					h2 = q
				}
			}
			if h1 > maxH {
				maxH = h1
			}
			x := h1
			if h2 < x {
				x = h2
			}
			x--
			if x >= 0 {
				val := h1 + h2 - 2*dep + 1
				if val > d[x] {
					d[x] = val
				}
			}
			return h1
		}
		dfs(0, -1, 0)
		for i := n - 2; i >= 0; i-- {
			if d[i+1] > d[i] {
				d[i] = d[i+1]
			}
		}
		ans := 0
		for i := 1; i <= n; i++ {
			for ans < maxH && d[ans]/2+i > ans {
				ans++
			}
			fmt.Fprint(writer, ans)
			if i < n {
				fmt.Fprint(writer, " ")
			}
		}
		fmt.Fprintln(writer)
	}
}
