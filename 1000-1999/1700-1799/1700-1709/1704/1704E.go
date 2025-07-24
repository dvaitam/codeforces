package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		adj := make([][]int, n)
		outDeg := make([]int, n)
		for i := 0; i < m; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			y--
			adj[x] = append(adj[x], y)
			outDeg[x]++
		}

		// simulate up to n steps
		ans := 0
		for step := 0; step < n; step++ {
			allZero := true
			inc := make([]int64, n)
			for i := 0; i < n; i++ {
				if a[i] > 0 {
					allZero = false
					a[i]--
					for _, v := range adj[i] {
						inc[v]++
					}
				}
			}
			for i := 0; i < n; i++ {
				a[i] += inc[i]
			}
			if allZero {
				ans = step
				break
			}
			if step == n-1 {
				ans = n
			}
		}
		if ans < n {
			fmt.Fprintln(out, ans)
			continue
		}

		// topological order
		inDeg := make([]int, n)
		for u := 0; u < n; u++ {
			for _, v := range adj[u] {
				inDeg[v]++
			}
		}
		queue := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if inDeg[i] == 0 {
				queue = append(queue, i)
			}
		}
		order := make([]int, 0, n)
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			order = append(order, u)
			for _, v := range adj[u] {
				inDeg[v]--
				if inDeg[v] == 0 {
					queue = append(queue, v)
				}
			}
		}

		for _, u := range order {
			for _, v := range adj[u] {
				a[v] = (a[v] + a[u]) % MOD
			}
		}
		sink := 0
		for i := 0; i < n; i++ {
			if outDeg[i] == 0 {
				sink = i
				break
			}
		}
		res := (int64(ans) + a[sink]) % MOD
		fmt.Fprintln(out, res)
	}
}
