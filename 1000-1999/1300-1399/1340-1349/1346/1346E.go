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

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	const INF int64 = 1 << 60
	cost := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		cost[i] = INF
	}
	cost[k] = 0
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		cx, cy := cost[x], cost[y]
		newcx := cy
		if cx+1 < newcx {
			newcx = cx + 1
		}
		newcy := cx
		if cy+1 < newcy {
			newcy = cy + 1
		}
		cost[x], cost[y] = newcx, newcy
	}
	for i := 1; i <= n; i++ {
		if cost[i] >= INF {
			fmt.Fprint(out, -1)
		} else {
			fmt.Fprint(out, cost[i])
		}
		if i == n {
			fmt.Fprintln(out)
		} else {
			fmt.Fprint(out, " ")
		}
	}
}
