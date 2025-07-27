package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		// prepare LCA
		log := 0
		for (1 << log) <= n {
			log++
		}
		up := make([][]int, log)
		for i := range up {
			up[i] = make([]int, n)
		}
		depth := make([]int, n)
		sz := make([]int, n)

		var dfs func(v, p int)
		dfs = func(v, p int) {
			up[0][v] = p
			sz[v] = 1
			for _, to := range adj[v] {
				if to == p {
					continue
				}
				depth[to] = depth[v] + 1
				dfs(to, v)
				sz[v] += sz[to]
			}
		}
		dfs(0, 0)
		for i := 1; i < log; i++ {
			for v := 0; v < n; v++ {
				up[i][v] = up[i-1][up[i-1][v]]
			}
		}

		lca := func(a, b int) int {
			if depth[a] < depth[b] {
				a, b = b, a
			}
			diff := depth[a] - depth[b]
			for i := 0; diff > 0; i++ {
				if diff&1 == 1 {
					a = up[i][a]
				}
				diff >>= 1
			}
			if a == b {
				return a
			}
			for i := log - 1; i >= 0; i-- {
				if up[i][a] != up[i][b] {
					a = up[i][a]
					b = up[i][b]
				}
			}
			return up[0][a]
		}

		dist := func(a, b int) int {
			l := lca(a, b)
			return depth[a] + depth[b] - 2*depth[l]
		}

		isAnc := func(a, b int) bool { return lca(a, b) == a }

		jump := func(v, k int) int {
			for i := 0; k > 0; i++ {
				if k&1 == 1 {
					v = up[i][v]
				}
				k >>= 1
			}
			return v
		}

		onPath := func(a, b, c int) bool {
			return dist(a, c)+dist(c, b) == dist(a, b)
		}

		nextNode := func(a, b int) int {
			if a == b {
				return a
			}
			if isAnc(a, b) {
				return jump(b, depth[b]-depth[a]-1)
			}
			return up[0][a]
		}

		compSize := func(a, b int) int {
			if a == b {
				return n
			}
			if isAnc(a, b) {
				child := nextNode(a, b)
				return n - sz[child]
			}
			return sz[a]
		}

		totalPairs := int64(n * (n - 1) / 2)
		F := make([]int64, n+2)
		F[0] = totalPairs

		// F[1]
		sum := int64(0)
		for _, to := range adj[0] {
			if up[0][to] == 0 {
				s := sz[to]
				sum += int64(s*(s-1)) / 2
			}
		}
		F[1] = totalPairs - sum

		l, r := 0, 0
		valid := true
		for k := 2; k <= n && valid; k++ {
			x := k - 1
			if onPath(l, r, x) {
				// ok
			} else if onPath(l, x, r) {
				r = x
			} else if onPath(r, x, l) {
				l = x
			} else {
				valid = false
				break
			}
			if valid {
				s1 := compSize(l, r)
				s2 := compSize(r, l)
				F[k] = int64(s1) * int64(s2)
			}
		}

		ans := make([]int64, n+1)
		for i := 0; i <= n; i++ {
			ans[i] = F[i] - F[i+1]
		}

		for i := 0; i <= n; i++ {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(ans[i])
		}
		fmt.Println()
	}
}
