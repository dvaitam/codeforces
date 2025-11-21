package main

import (
	"bufio"
	"fmt"
	"os"
)

const logN = 20

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var w int64
		fmt.Fscan(in, &n, &w)

		parent := make([]int, n+1)
		parent[1] = 1
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &parent[i])
		}

		depth := make([]int, n+1)
		up := make([][]int, logN)
		for k := 0; k < logN; k++ {
			up[k] = make([]int, n+1)
			up[k][1] = 1
		}
		for i := 2; i <= n; i++ {
			p := parent[i]
			depth[i] = depth[p] + 1
			up[0][i] = p
			for k := 1; k < logN; k++ {
				up[k][i] = up[k-1][up[k-1][i]]
			}
		}

		size := make([]int, n+1)
		for i := 1; i <= n; i++ {
			size[i] = 1
		}
		for i := n; i >= 2; i-- {
			size[parent[i]] += size[i]
		}

		tout := make([]int, n+1)
		for i := 1; i <= n; i++ {
			tout[i] = i + size[i] - 1
		}

		rem := make([]int, n+1)
		for i := 1; i <= n; i++ {
			u := i
			v := i + 1
			if i == n {
				v = 1
			}
			rem[i] = distance(u, v, depth, up)
		}

		kDone := 0
		sumKnown := int64(0)
		ans := make([]int64, 0, n-1)

		for ev := 0; ev < n-1; ev++ {
			var x int
			var y int64
			fmt.Fscan(in, &x, &y)
			sumKnown += y

			idx1 := x - 1
			rem[idx1]--
			if rem[idx1] == 0 {
				kDone++
			}
			idx2 := tout[x]
			rem[idx2]--
			if rem[idx2] == 0 {
				kDone++
			}

			remaining := w - sumKnown
			base := sumKnown * 2
			cur := base + remaining*int64(n-kDone)
			ans = append(ans, cur)
		}

		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

func distance(u, v int, depth []int, up [][]int) int {
	l := lca(u, v, depth, up)
	return depth[u] + depth[v] - 2*depth[l]
}

func lca(u, v int, depth []int, up [][]int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for k := 0; diff > 0; k++ {
		if diff&1 == 1 {
			u = up[k][u]
		}
		diff >>= 1
	}
	if u == v {
		return u
	}
	for k := logN - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			u = up[k][u]
			v = up[k][v]
		}
	}
	return up[0][u]
}
