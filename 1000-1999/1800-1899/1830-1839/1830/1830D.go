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
		g := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			g[a] = append(g[a], b)
			g[b] = append(g[b], a)
		}
		color := make([]int, n+1)
		for i := range color {
			color[i] = -1
		}
		cnt := [2]int{}
		queue := []int{1}
		color[1] = 0
		cnt[0]++
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for _, to := range g[v] {
				if color[to] == -1 {
					color[to] = 1 - color[v]
					cnt[color[to]]++
					queue = append(queue, to)
				}
			}
		}
		if cnt[1] > cnt[0] {
			cnt[1], cnt[0] = cnt[0], cnt[1]
		}
		n64 := int64(n)
		ans := n64*n64 - int64(cnt[1])
		fmt.Fprintln(out, ans)
	}
}
