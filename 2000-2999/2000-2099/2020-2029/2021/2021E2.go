package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	u, v int
	w    int64
}

type dsu struct {
	parent []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
	}
	return &dsu{parent: p}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

const inf int64 = 1 << 62

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, p int
		fmt.Fscan(in, &n, &m, &p)
		required := make([]bool, 2*n+5)
		for i := 0; i < p; i++ {
			var s int
			fmt.Fscan(in, &s)
			required[s] = true
		}
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			var u, v int
			var w int64
			fmt.Fscan(in, &u, &v, &w)
			edges[i] = edge{u: u, v: v, w: w}
		}
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].w < edges[j].w
		})

		maxNodes := 2*n + 5
		d := newDSU(maxNodes)
		comp := make([]int, maxNodes)
		for i := 1; i <= n; i++ {
			comp[i] = i
		}
		children := make([][2]int, maxNodes)
		weight := make([]int64, maxNodes)
		nextID := n
		for _, e := range edges {
			ru := d.find(e.u)
			rv := d.find(e.v)
			if ru == rv {
				continue
			}
			nextID++
			newID := nextID
			children[newID] = [2]int{comp[ru], comp[rv]}
			weight[newID] = e.w
			d.parent[ru] = newID
			d.parent[rv] = newID
			d.parent[newID] = newID
			comp[newID] = newID
		}
		root := comp[d.find(1)]

		var dfs func(int) ([]int64, int)
		dfs = func(u int) ([]int64, int) {
			if u <= n {
				dp := make([]int64, 2)
				dp[0], dp[1] = 0, 0
				if required[u] {
					return dp, 1
				}
				return dp, 0
			}
			left := children[u][0]
			right := children[u][1]
			dpL, termL := dfs(left)
			dpR, termR := dfs(right)
			sizeL := len(dpL) - 1
			sizeR := len(dpR) - 1
			newDP := make([]int64, sizeL+sizeR+1)
			for i := range newDP {
				newDP[i] = inf
			}
			w := weight[u]
			for i := 0; i <= sizeL; i++ {
				if dpL[i] >= inf {
					continue
				}
				for j := 0; j <= sizeR; j++ {
					if dpR[j] >= inf {
						continue
					}
					extra := int64(0)
					if i > 0 && j == 0 {
						extra = int64(termR) * w
					} else if i == 0 && j > 0 {
						extra = int64(termL) * w
					}
					val := dpL[i] + dpR[j] + extra
					if val < newDP[i+j] {
						newDP[i+j] = val
					}
				}
			}
			return newDP, termL + termR
		}

		dpRoot, _ := dfs(root)
		ans := make([]int64, n+1)
		cur := inf
		for k := 1; k <= n; k++ {
			if k < len(dpRoot) && dpRoot[k] < cur {
				cur = dpRoot[k]
			}
			ans[k] = cur
		}
		for k := 1; k <= n; k++ {
			if k > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[k])
		}
		fmt.Fprintln(out)
	}
}
