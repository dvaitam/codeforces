package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problem E from contest 1558.
// We need the minimal initial power so that starting from cave 1,
// we can defeat monsters in all other caves. After entering cave i
// for the first time with power > a[i], the hero's power increases by b[i].
// The graph is connected and each vertex has at least two edges.
// The restriction of not using the same tunnel twice in a row does not
// affect the minimal power required, because every cave has another exit
// allowing us to walk around cycles. Thus we simply check if we can
// expand the visited set greedily with a given starting power.

func canClear(n int, a, b []int64, g [][]int, start int64) bool {
	visited := make([]bool, n+1)
	visited[1] = true
	power := start
	changed := true
	for changed {
		changed = false
		for u := 1; u <= n; u++ {
			if !visited[u] {
				continue
			}
			for _, v := range g[u] {
				if visited[v] {
					continue
				}
				if power > a[v] {
					power += b[v]
					visited[v] = true
					changed = true
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n+1)
		b := make([]int64, n+1)
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &b[i])
		}
		g := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		// binary search on the minimal starting power
		low, high := int64(1), int64(1)
		for i := 2; i <= n; i++ {
			if a[i]+1 > high {
				high = a[i] + 1
			}
		}
		// high is at least max(a)+1. Add total gains for safety.
		var sum int64
		for i := 2; i <= n; i++ {
			sum += b[i]
		}
		if high < sum+1 {
			high = sum + 1
		}
		for low < high {
			mid := (low + high) / 2
			if canClear(n, a, b, g, mid) {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(out, low)
	}
}
