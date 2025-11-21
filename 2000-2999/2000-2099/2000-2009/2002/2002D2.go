package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		parent := make([]int, n+1)
		parent[1] = 0
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &parent[i])
		}

		depth := make([]int, n+1)
		up := make([][LOG]int, n+1)
		for i := 2; i <= n; i++ {
			p := parent[i]
			depth[i] = depth[p] + 1
			up[i][0] = p
			for k := 1; k < LOG; k++ {
				up[i][k] = up[up[i][k-1]][k-1]
			}
		}

		lca := func(u, v int) int {
			if depth[u] < depth[v] {
				u, v = v, u
			}
			// lift u up
			diff := depth[u] - depth[v]
			for k := 0; diff > 0; k++ {
				if diff&1 == 1 {
					u = up[u][k]
				}
				diff >>= 1
			}
			if u == v {
				return u
			}
			for k := LOG - 1; k >= 0; k-- {
				if up[u][k] != up[v][k] {
					u = up[u][k]
					v = up[v][k]
				}
			}
			return parent[u]
		}

		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &p[i])
		}

		good := make([]bool, n+1)
		countBad := 0

		checkPair := func(idx int) bool {
			u := p[idx]
			v := p[idx+1]
			return parent[v] == lca(u, v)
		}

		for i := 1; i <= n-1; i++ {
			good[i] = checkPair(i)
			if !good[i] {
				countBad++
			}
		}

		rootOK := p[1] == 1

		updateIndex := func(idx int) {
			if idx < 1 || idx >= n {
				return
			}
			newVal := checkPair(idx)
			if good[idx] && !newVal {
				countBad++
			} else if !good[idx] && newVal {
				countBad--
			}
			good[idx] = newVal
		}

		for ; q > 0; q-- {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x > y {
				x, y = y, x
			}
			if x == y {
				if rootOK && countBad == 0 {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
				continue
			}

			p[x], p[y] = p[y], p[x]
			if x == 1 || y == 1 {
				rootOK = p[1] == 1
			}

			updateIndex(x - 1)
			updateIndex(x)
			updateIndex(y - 1)
			updateIndex(y)

			if rootOK && countBad == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
