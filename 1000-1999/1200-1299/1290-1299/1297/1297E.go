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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		g := make([][]int, n)
		deg := make([]int, n)
		for i := 0; i < n-1; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			y--
			g[x] = append(g[x], y)
			g[y] = append(g[y], x)
			deg[x]++
			deg[y]++
		}
		if k == 1 {
			fmt.Fprintln(out, "Yes")
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, 1)
			continue
		}
		leaves := make([]int, 0)
		for i := 0; i < n; i++ {
			if deg[i] == 1 {
				leaves = append(leaves, i)
			}
		}
		if len(leaves) < k {
			fmt.Fprintln(out, "No")
			continue
		}
		removed := make([]bool, n)
		queue := append([]int(nil), leaves...)
		head := 0
		leafCount := len(leaves)
		for leafCount > k && head < len(queue) {
			u := queue[head]
			head++
			if removed[u] || deg[u] != 1 {
				continue
			}
			removed[u] = true
			leafCount--
			for _, v := range g[u] {
				if removed[v] {
					continue
				}
				deg[v]--
				if deg[v] == 1 {
					queue = append(queue, v)
					leafCount++
				}
			}
		}
		if leafCount != k {
			fmt.Fprintln(out, "No")
			continue
		}
		ans := make([]int, 0)
		for i := 0; i < n; i++ {
			if !removed[i] {
				ans = append(ans, i+1)
			}
		}
		if len(ans) == 0 {
			fmt.Fprintln(out, "No")
			continue
		}
		fmt.Fprintln(out, "Yes")
		fmt.Fprintln(out, len(ans))
		for i, v := range ans {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, v)
		}
		out.WriteByte('\n')
	}
}
