package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXLOG = 20

var (
	adj    [][]int
	up     [MAXLOG][]int
	depth  []int
	parent []int
)

func preprocess(n int) []int {
	depth = make([]int, n)
	parent = make([]int, n)
	for i := 0; i < MAXLOG; i++ {
		up[i] = make([]int, n)
	}
	stack := []int{0}
	parent[0] = 0
	order := make([]int, 0, n)
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	up[0] = parent
	for k := 1; k < MAXLOG; k++ {
		for i := 0; i < n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
	return order
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	for k := MAXLOG - 1; k >= 0; k-- {
		if depth[u]-(1<<k) >= depth[v] {
			u = up[k][u]
		}
	}
	if u == v {
		return u
	}
	for k := MAXLOG - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			u = up[k][u]
			v = up[k][v]
		}
	}
	return up[0][u]
}

func jump(u, k int) int {
	i := 0
	for k > 0 {
		if k&1 == 1 {
			u = up[i][u]
		}
		k >>= 1
		i++
	}
	return u
}

func childOnPath(u, l int) int {
	if u == l {
		return -1
	}
	return jump(u, depth[u]-depth[l]-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	order := preprocess(n)

	var m int
	fmt.Fscan(in, &m)

	nodeAdd := make([]int64, n)
	edgeAdd := make([]int64, n)
	pairCC := make([]map[uint64]int, n)
	sumPairChild := make([]map[int]int, n)
	endPair := make([]map[int]int, n)
	single := make([]int, n)

	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		w := lca(u, v)
		nodeAdd[u]++
		nodeAdd[v]++
		nodeAdd[w]--
		if w != 0 {
			nodeAdd[parent[w]]--
		}
		edgeAdd[u]++
		edgeAdd[v]++
		edgeAdd[w] -= 2
		a := childOnPath(u, w)
		b := childOnPath(v, w)
		if a == -1 && b == -1 {
			single[w]++
		} else if a == -1 || b == -1 {
			c := a
			if c == -1 {
				c = b
			}
			if endPair[w] == nil {
				endPair[w] = make(map[int]int)
			}
			endPair[w][c]++
		} else {
			if a > b {
				a, b = b, a
			}
			if pairCC[w] == nil {
				pairCC[w] = make(map[uint64]int)
			}
			key := uint64(a)<<32 | uint64(b)
			pairCC[w][key]++
			if sumPairChild[w] == nil {
				sumPairChild[w] = make(map[int]int)
			}
			sumPairChild[w][a]++
			sumPairChild[w][b]++
		}
	}

	for i := len(order) - 1; i > 0; i-- {
		v := order[i]
		p := parent[v]
		nodeAdd[p] += nodeAdd[v]
		edgeAdd[p] += edgeAdd[v]
	}

	nodeCnt := nodeAdd
	edgeCnt := make([]int64, n)
	for i := 1; i < n; i++ {
		edgeCnt[i] = edgeAdd[i]
	}

	var ans int64
	for x := 0; x < n; x++ {
		mcnt := nodeCnt[x]
		totalPairs := mcnt * (mcnt - 1) / 2
		var sumDir int64
		for _, c := range adj[x] {
			if c == parent[x] {
				continue
			}
			d := edgeCnt[c]
			sumDir += d * (d - 1) / 2
		}
		if x != 0 {
			d := edgeCnt[x]
			sumDir += d * (d - 1) / 2
		}
		var sumPair int64
		if mp := pairCC[x]; mp != nil {
			for _, cnt := range mp {
				c := int64(cnt)
				sumPair += c * (c - 1) / 2
			}
		}
		if x != 0 {
			for _, c := range adj[x] {
				if c == parent[x] {
					continue
				}
				cnt := edgeCnt[c]
				if mp := endPair[x]; mp != nil {
					cnt -= int64(mp[c])
				}
				if mp := sumPairChild[x]; mp != nil {
					cnt -= int64(mp[c])
				}
				sumPair += cnt * (cnt - 1) / 2
			}
		}
		ans += totalPairs - sumDir + sumPair
	}

	fmt.Fprintln(out, ans)
}
