package main

import (
	"bufio"
	"fmt"
	"os"
)

type frame struct {
	v, p  int
	idx   int
	saved int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		colors := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &colors[i])
		}
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		cnt := make([]int, n+1)
		var ans int64
		stack := []frame{{v: 1, p: 0, idx: -1, saved: 0}}
		for len(stack) > 0 {
			f := &stack[len(stack)-1]
			c := colors[f.v]
			if f.idx == -1 {
				ans += int64(cnt[c])
				f.saved = cnt[c]
				f.idx = 0
			}
			if f.idx < len(adj[f.v]) {
				to := adj[f.v][f.idx]
				f.idx++
				if to == f.p {
					continue
				}
				cnt[c] = 1
				stack = append(stack, frame{v: to, p: f.v, idx: -1})
				continue
			}
			cnt[c] = f.saved + 1
			stack = stack[:len(stack)-1]
		}
		fmt.Fprintln(out, ans)
	}
}
