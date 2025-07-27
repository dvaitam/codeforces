package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	// dp holds reachable states after current step
	dp := make([]bool, m+1)
	dp[0] = true
	ans := make([]int, m+1)
	for i := 1; i <= m; i++ {
		ans[i] = -1
	}
	for step := 1; step <= n; step++ {
		var t, xprime, y int
		fmt.Fscan(in, &t, &xprime, &y)
		newdp := make([]bool, m+1)
		// distance array for BFS
		dist := make([]int, m+1)
		for i := range dist {
			dist[i] = -1
		}
		queue := make([]int, 0)
		for v := 0; v <= m; v++ {
			if dp[v] {
				newdp[v] = true
				dist[v] = 0
				queue = append(queue, v)
			}
		}
		// conversion factor
		if t == 1 {
			add := (xprime + 100000 - 1) / 100000
			head := 0
			for head < len(queue) {
				v := queue[head]
				head++
				if dist[v] == y {
					continue
				}
				nxt := v + add
				if nxt <= m && dist[nxt] == -1 {
					dist[nxt] = dist[v] + 1
					newdp[nxt] = true
					queue = append(queue, nxt)
				}
			}
		} else {
			// type 2
			head := 0
			for head < len(queue) {
				v := queue[head]
				head++
				if dist[v] == y {
					continue
				}
				nxt := int((int64(v)*int64(xprime) + 100000 - 1) / 100000)
				if nxt <= m && dist[nxt] == -1 {
					dist[nxt] = dist[v] + 1
					newdp[nxt] = true
					queue = append(queue, nxt)
				}
			}
		}
		dp = newdp
		for j := 1; j <= m; j++ {
			if ans[j] == -1 && dp[j] {
				ans[j] = step
			}
		}
	}
	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= m; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
	out.Flush()
}
