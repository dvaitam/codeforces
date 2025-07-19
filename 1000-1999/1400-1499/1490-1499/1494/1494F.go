package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU with union by size and path compression
type DSU struct {
	parentOrSize []int
}

// NewDSU creates a DSU for n elements
func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = -1
	}
	return &DSU{parentOrSize: p}
}

// Leader returns the leader of a
func (d *DSU) Leader(a int) int {
	if d.parentOrSize[a] < 0 {
		return a
	}
	d.parentOrSize[a] = d.Leader(d.parentOrSize[a])
	return d.parentOrSize[a]
}

// Merge unions a and b, returns new leader
func (d *DSU) Merge(a, b int) int {
	x := d.Leader(a)
	y := d.Leader(b)
	if x == y {
		return x
	}
	if -d.parentOrSize[x] < -d.parentOrSize[y] {
		x, y = y, x
	}
	d.parentOrSize[x] += d.parentOrSize[y]
	d.parentOrSize[y] = x
	return x
}

// Same returns true if a and b are in same set
func (d *DSU) Same(a, b int) bool {
	return d.Leader(a) == d.Leader(b)
}

// Groups returns the components
func (d *DSU) Groups() [][]int {
	n := len(d.parentOrSize)
	leaderBuf := make([]int, n)
	groupSize := make([]int, n)
	for i := 0; i < n; i++ {
		leaderBuf[i] = d.Leader(i)
		groupSize[leaderBuf[i]]++
	}
	res := make([][]int, n)
	for i := 0; i < n; i++ {
		if groupSize[i] > 0 {
			res[i] = make([]int, 0, groupSize[i])
		}
	}
	for i := 0; i < n; i++ {
		res[leaderBuf[i]] = append(res[leaderBuf[i]], i)
	}
	out := make([][]int, 0, len(res))
	for _, g := range res {
		if len(g) > 0 {
			out = append(out, g)
		}
	}
	return out
}

// UndirectedEulerianGraph supports Euler circuits and trails
type UndirectedEulerianGraph struct {
	n   int
	g   [][]int
	inv [][]int
}

// NewEulerGraph creates a graph with n vertices
func NewEulerGraph(n int) *UndirectedEulerianGraph {
	return &UndirectedEulerianGraph{
		n:   n,
		g:   make([][]int, n),
		inv: make([][]int, n),
	}
}

// AddEdge adds an undirected edge u-v
func (h *UndirectedEulerianGraph) AddEdge(u, v int) {
	su := len(h.g[u])
	sv := len(h.g[v])
	h.g[u] = append(h.g[u], v)
	h.inv[u] = append(h.inv[u], sv)
	h.g[v] = append(h.g[v], u)
	h.inv[v] = append(h.inv[v], su)
}

// Eulerian circuit, returns path and ok
func (h *UndirectedEulerianGraph) EulerianCircuit(start int) ([]int, bool) {
	edgeNum := 0
	used := make([][]bool, h.n)
	for i := 0; i < h.n; i++ {
		sz := len(h.g[i])
		if sz&1 == 1 {
			return nil, false
		}
		edgeNum += sz
		used[i] = make([]bool, sz)
	}
	edgeNum /= 2
	res := make([]int, 0, edgeNum+1)
	var dfs func(u int)
	dfs = func(u int) {
		for i := 0; i < len(h.g[u]); i++ {
			if used[u][i] {
				continue
			}
			v := h.g[u][i]
			used[u][i] = true
			used[v][h.inv[u][i]] = true
			dfs(v)
		}
		res = append(res, u)
	}
	dfs(start)
	if len(res) != edgeNum+1 {
		return nil, false
	}
	return res, true
}

// EulerianTrail returns trail or false
func (h *UndirectedEulerianGraph) EulerianTrail() ([]int, bool) {
	s, t, invalid := -1, -1, -1
	for i := 0; i < h.n; i++ {
		if len(h.g[i])&1 == 1 {
			if s < 0 {
				s = i
			} else if t < 0 {
				t = i
			} else {
				invalid = i
			}
		}
	}
	if s < 0 || t < 0 || invalid >= 0 {
		return nil, false
	}
	// add edge s-t
	h.AddEdge(s, t)
	res, ok := h.EulerianCircuit(s)
	// remove added edge
	h.g[s] = h.g[s][:len(h.g[s])-1]
	h.inv[s] = h.inv[s][:len(h.inv[s])-1]
	h.g[t] = h.g[t][:len(h.g[t])-1]
	h.inv[t] = h.inv[t][:len(h.inv[t])-1]
	if !ok {
		return nil, false
	}
	// drop last repeated vertex
	res = res[:len(res)-1]
	if res[len(res)-1] == t {
		return res, true
	}
	// rotate to start after s-t edge
	m := len(res)
	for i := 0; i < m-1; i++ {
		u, v := res[i], res[i+1]
		if (u == s && v == t) || (u == t && v == s) {
			// new res = res[i+1:]+res[:i+1]
			newRes := append(res[i+1:], res[:i+1]...)
			return newRes, true
		}
	}
	return res, true
}

func reverse(a []int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	g := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	out := func(ans []int) {
		fmt.Println(len(ans))
		for i, x := range ans {
			if x < 0 {
				fmt.Print(-1)
			} else {
				fmt.Print(x + 1)
			}
			if i+1 < len(ans) {
				fmt.Print(" ")
			}
		}
		fmt.Println()
		os.Exit(0)
	}
	// try each t
	for t := 0; t < n; t++ {
		orig := g[t]
		erased := make([]int, 0)
		gt := make([]int, 0)
		// separate edges
		for _, i := range orig {
			if len(g[i])&1 == 1 {
				// erase i-t edge
				erased = append(erased, i)
				for k, v := range g[i] {
					if v == t {
						g[i] = append(g[i][:k], g[i][k+1:]...)
						break
					}
				}
			} else {
				gt = append(gt, i)
			}
		}
		g[t] = gt
		// count odd
		odd := 0
		for i := 0; i < n; i++ {
			if i == t {
				continue
			}
			if len(g[i])&1 == 1 {
				odd++
			}
		}
		// restore helper
		restore := func() {
			for _, i := range erased {
				g[i] = append(g[i], t)
				g[t] = append(g[t], i)
			}
		}
		if odd >= 2 {
			restore()
			continue
		}
		if odd == 1 {
			h := NewEulerGraph(n)
			for i := 0; i < n; i++ {
				for _, j := range g[i] {
					if i < j {
						h.AddEdge(i, j)
					}
				}
			}
			if ans, ok := h.EulerianTrail(); ok {
				if ans[0] == t {
					reverse(ans)
				}
				ans = append(ans, -1)
				for _, i := range erased {
					ans = append(ans, i)
					ans = append(ans, t)
				}
				out(ans)
			}
			restore()
			continue
		}
		// even case
		uf := NewDSU(n)
		for i := 0; i < n; i++ {
			for _, j := range g[i] {
				uf.Merge(i, j)
			}
		}
		c := 0
		x := -1
		for _, comp := range uf.Groups() {
			if len(comp) == 1 {
				if comp[0] == t {
					c++
				}
			} else {
				c++
				found := false
				for _, v := range comp {
					if v == t {
						found = true
						break
					}
				}
				if !found {
					x = comp[0]
				}
			}
		}
		if c > 2 {
			restore()
			continue
		}
		if c == 2 {
			// find erased connected to x
			idx := -1
			for k, y := range erased {
				if uf.Same(x, y) {
					idx = k
					break
				}
			}
			v := erased[idx]
			// remove
			erased = append(erased[:idx], erased[idx+1:]...)
			// add t-v
			g[t] = append(g[t], v)
			g[v] = append(g[v], t)
			h := NewEulerGraph(n)
			for i := 0; i < n; i++ {
				for _, j := range g[i] {
					if i < j {
						h.AddEdge(i, j)
					}
				}
			}
			ans, _ := h.EulerianTrail()
			if ans[0] == t {
				reverse(ans)
			}
			ans = append(ans, -1)
			for _, i := range erased {
				ans = append(ans, i)
				ans = append(ans, t)
			}
			out(ans)
		}
		// circuit case
		h := NewEulerGraph(n)
		for i := 0; i < n; i++ {
			for _, j := range g[i] {
				if i < j {
					h.AddEdge(i, j)
				}
			}
		}
		ans, _ := h.EulerianCircuit(t)
		if ans[0] == t {
			reverse(ans)
		}
		ans = append(ans, -1)
		for _, i := range erased {
			ans = append(ans, i)
			ans = append(ans, t)
		}
		out(ans)
	}
	// no solution
	fmt.Println(0)
}
