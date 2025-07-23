package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct{ to, id int }

const LOG = 20

var (
	n, m, q        int
	edges          []struct{ u, v int }
	adj            [][]Edge
	tin, low       []int
	timer          int
	isBridge       []bool
	comp           []int
	comps          int
	tree           [][]int
	parent         [LOG][]int
	depth          []int
	root           []int
	upReq, downReq []int
	upSum, downSum []int
	impossible     bool
)

func dfsBridge(v, pid int) {
	timer++
	tin[v] = timer
	low[v] = timer
	for _, e := range adj[v] {
		if e.id == pid {
			continue
		}
		if tin[e.to] == 0 {
			dfsBridge(e.to, e.id)
			if low[e.to] > tin[v] {
				isBridge[e.id] = true
			}
			if low[e.to] < low[v] {
				low[v] = low[e.to]
			}
		} else if tin[e.to] < tin[v] {
			if tin[e.to] < low[v] {
				low[v] = tin[e.to]
			}
		}
	}
}

func dfsComp(v, cid int) {
	stack := []int{v}
	comp[v] = cid
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range adj[x] {
			if isBridge[e.id] {
				continue
			}
			if comp[e.to] == -1 {
				comp[e.to] = cid
				stack = append(stack, e.to)
			}
		}
	}
}

func dfsLCA(v, r, p int, vis []bool) {
	vis[v] = true
	root[v] = r
	if p == -1 {
		parent[0][v] = v
		depth[v] = 0
	} else {
		parent[0][v] = p
		depth[v] = depth[p] + 1
	}
	for _, to := range tree[v] {
		if to == p {
			continue
		}
		dfsLCA(to, r, v, vis)
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff>>i&1 == 1 {
			u = parent[i][u]
		}
	}
	if u == v {
		return u
	}
	for i := LOG - 1; i >= 0; i-- {
		if parent[i][u] != parent[i][v] {
			u = parent[i][u]
			v = parent[i][v]
		}
	}
	return parent[0][u]
}

func dfsCheck(v, p int, vis []bool) {
	vis[v] = true
	for _, to := range tree[v] {
		if to == p {
			continue
		}
		dfsCheck(to, v, vis)
		if upSum[to] > 0 && downSum[to] > 0 {
			impossible = true
		}
		upSum[v] += upSum[to]
		downSum[v] += downSum[to]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m, &q)
	adj = make([][]Edge, n)
	edges = make([]struct{ u, v int }, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = struct{ u, v int }{u, v}
		adj[u] = append(adj[u], Edge{v, i})
		adj[v] = append(adj[v], Edge{u, i})
	}
	tin = make([]int, n)
	low = make([]int, n)
	isBridge = make([]bool, m)
	timer = 0
	for i := 0; i < n; i++ {
		if tin[i] == 0 {
			dfsBridge(i, -1)
		}
	}
	comp = make([]int, n)
	for i := 0; i < n; i++ {
		comp[i] = -1
	}
	comps = 0
	for i := 0; i < n; i++ {
		if comp[i] == -1 {
			dfsComp(i, comps)
			comps++
		}
	}
	tree = make([][]int, comps)
	for id, e := range edges {
		if isBridge[id] {
			cu := comp[e.u]
			cv := comp[e.v]
			tree[cu] = append(tree[cu], cv)
			tree[cv] = append(tree[cv], cu)
		}
	}
	for i := 0; i < LOG; i++ {
		parent[i] = make([]int, comps)
	}
	depth = make([]int, comps)
	root = make([]int, comps)
	vis := make([]bool, comps)
	for i := 0; i < comps; i++ {
		if !vis[i] {
			dfsLCA(i, i, -1, vis)
		}
	}
	for k := 1; k < LOG; k++ {
		for i := 0; i < comps; i++ {
			parent[k][i] = parent[k-1][parent[k-1][i]]
		}
	}
	upReq = make([]int, comps)
	downReq = make([]int, comps)
	for i := 0; i < q; i++ {
		var s, d int
		fmt.Fscan(in, &s, &d)
		s--
		d--
		cs := comp[s]
		ct := comp[d]
		if root[cs] != root[ct] {
			fmt.Println("No")
			return
		}
		l := lca(cs, ct)
		upReq[cs]++
		upReq[l]--
		downReq[ct]++
		downReq[l]--
	}
	upSum = make([]int, comps)
	downSum = make([]int, comps)
	for i := 0; i < comps; i++ {
		upSum[i] = upReq[i]
		downSum[i] = downReq[i]
	}
	vis = make([]bool, comps)
	for i := 0; i < comps; i++ {
		if !vis[i] {
			dfsCheck(i, -1, vis)
		}
	}
	if impossible {
		fmt.Println("No")
	} else {
		fmt.Println("Yes")
	}
}
