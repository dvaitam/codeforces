package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const LOG int = 20
const BITS int = 30

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		val := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &val[i])
		}
		g := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		parent := make([][]int, LOG)
		upOr := make([][]int, LOG)
		for i := 0; i < LOG; i++ {
			parent[i] = make([]int, n+1)
			upOr[i] = make([]int, n+1)
		}
		depth := make([]int, n+1)
		ancBit := make([][]int, BITS)
		for b := 0; b < BITS; b++ {
			ancBit[b] = make([]int, n+1)
		}

		// BFS from root 1
		q := make([]int, 0, n)
		q = append(q, 1)
		parent[0][1] = 0
		depth[1] = 0
		for b := 0; b < BITS; b++ {
			if (val[1]>>b)&1 == 1 {
				ancBit[b][1] = 1
			} else {
				ancBit[b][1] = 0
			}
		}
		for idx := 0; idx < len(q); idx++ {
			u := q[idx]
			for _, v := range g[u] {
				if v == parent[0][u] {
					continue
				}
				parent[0][v] = u
				depth[v] = depth[u] + 1
				upOr[0][v] = val[u]
				for b := 0; b < BITS; b++ {
					if (val[v]>>b)&1 == 1 {
						ancBit[b][v] = v
					} else {
						ancBit[b][v] = ancBit[b][u]
					}
				}
				q = append(q, v)
			}
		}

		for k := 1; k < LOG; k++ {
			for v := 1; v <= n; v++ {
				anc := parent[k-1][v]
				parent[k][v] = parent[k-1][anc]
				upOr[k][v] = upOr[k-1][v] | upOr[k-1][anc]
			}
		}

		var lca func(int, int) int
		lca = func(a, b int) int {
			if depth[a] < depth[b] {
				a, b = b, a
			}
			diff := depth[a] - depth[b]
			for k := LOG - 1; k >= 0; k-- {
				if diff>>k&1 == 1 {
					a = parent[k][a]
				}
			}
			if a == b {
				return a
			}
			for k := LOG - 1; k >= 0; k-- {
				if parent[k][a] != parent[k][b] {
					a = parent[k][a]
					b = parent[k][b]
				}
			}
			return parent[0][a]
		}

		var orToAncestor func(int, int) int
		orToAncestor = func(u, anc int) int {
			res := val[u]
			if u == anc {
				return res
			}
			diff := depth[u] - depth[anc]
			for k := LOG - 1; k >= 0; k-- {
				if diff>>k&1 == 1 {
					res |= upOr[k][u]
					u = parent[k][u]
				}
			}
			return res
		}

		var pathOr func(int, int) int
		pathOr = func(u, v int) int {
			p := lca(u, v)
			return orToAncestor(u, p) | orToAncestor(v, p)
		}

		gather := func(start, stop int) []int {
			res := []int{start}
			if start == stop {
				return res
			}
			cur := start
			curOR := val[cur]
			for cur != stop {
				bestDepth := -1
				bestNode := 0
				for b := 0; b < BITS; b++ {
					if (curOR>>b)&1 == 0 {
						node := ancBit[b][cur]
						if node != 0 && depth[node] >= depth[stop] && depth[node] > bestDepth {
							bestDepth = depth[node]
							bestNode = node
						}
					}
				}
				if bestNode == 0 || bestDepth < depth[stop] {
					break
				}
				cur = bestNode
				if res[len(res)-1] != cur {
					res = append(res, cur)
				}
				curOR |= val[cur]
				if cur == stop {
					break
				}
			}
			if res[len(res)-1] != stop {
				res = append(res, stop)
			}
			return res
		}

		mark := make([]int, n+1)
		curMark := 0

		var qNum int
		fmt.Fscan(in, &qNum)
		for ; qNum > 0; qNum-- {
			var x, y int
			fmt.Fscan(in, &x, &y)
			l := lca(x, y)
			curMark++
			cands := make([]int, 0, 64)
			nodes := gather(x, l)
			for _, v := range nodes {
				if mark[v] != curMark {
					mark[v] = curMark
					cands = append(cands, v)
				}
			}
			nodes = gather(y, l)
			for _, v := range nodes {
				if mark[v] != curMark {
					mark[v] = curMark
					cands = append(cands, v)
				}
			}

			ans := 0
			for _, z := range cands {
				gx := pathOr(x, z)
				gy := pathOr(y, z)
				nic := bits.OnesCount(uint(gx)) + bits.OnesCount(uint(gy))
				if nic > ans {
					ans = nic
				}
			}
			fmt.Fprintln(out, ans)
		}
	}
}
