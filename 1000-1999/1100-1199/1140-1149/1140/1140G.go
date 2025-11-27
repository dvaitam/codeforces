package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	LOG = 20
	INF int64 = 4_000_000_000_000_000_000 // comfortably above any reachable path
)

type Edge struct {
	to     int
	w1, w2 int64
}

type Matrix struct {
	g [2][2]int64
}

func newMatrixInf() Matrix {
	return Matrix{g: [2][2]int64{{INF, INF}, {INF, INF}}}
}

func newMatrixIdentity() Matrix {
	m := newMatrixInf()
	m.g[0][0] = 0
	m.g[1][1] = 0
	return m
}

func add(a, b int64) int64 {
	if a >= INF || b >= INF {
		return INF
	}
	if a > INF-b {
		return INF
	}
	return a + b
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Multiply two transfer matrices (min-plus).
func multiply(a, b Matrix) Matrix {
	res := newMatrixInf()
	for i := 0; i < 2; i++ {
		for k := 0; k < 2; k++ {
			if a.g[i][k] >= INF {
				continue
			}
			for j := 0; j < 2; j++ {
				val := add(a.g[i][k], b.g[k][j])
				if val < res.g[i][j] {
					res.g[i][j] = val
				}
			}
		}
	}
	return res
}

// Transfer matrix across an edge directed u -> v.
// Parity 0 = odd vertex (2*u-1), parity 1 = even vertex (2*u).
func edgeMatrix(u, v int, w1, w2 int64, c []int64) Matrix {
	m := newMatrixInf()
	m.g[0][0] = min64(w1, add(c[u], add(w2, c[v])))
	m.g[1][1] = min64(w2, add(c[u], add(w1, c[v])))
	m.g[0][1] = min64(add(w1, c[v]), add(c[u], w2))
	m.g[1][0] = min64(add(w2, c[v]), add(c[u], w1))
	return m
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Buffer(make([]byte, 0, 1<<20), 1<<25)
	in.Split(bufio.ScanWords)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	nextInt := func() int {
		in.Scan()
		v, _ := strconv.Atoi(in.Text())
		return v
	}
	nextInt64 := func() int64 {
		in.Scan()
		v, _ := strconv.ParseInt(in.Text(), 10, 64)
		return v
	}

	n := nextInt()
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = nextInt64()
	}

	adj := make([][]Edge, n+1)
	for i := 1; i < n; i++ {
		u := nextInt()
		v := nextInt()
		w1 := nextInt64()
		w2 := nextInt64()
		adj[u] = append(adj[u], Edge{to: v, w1: w1, w2: w2})
		adj[v] = append(adj[v], Edge{to: u, w1: w1, w2: w2})
	}

	// Build parent info using iterative DFS to avoid deep recursion.
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	pw1 := make([]int64, n+1)
	pw2 := make([]int64, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	depth[1] = 1
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, e := range adj[u] {
			v := e.to
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			pw1[v] = e.w1
			pw2[v] = e.w2
			stack = append(stack, v)
		}
	}

	// Minimal cost to switch parity at each node (treat edges as w1+w2).
	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		c[i] = a[i]
	}
	// bottom-up
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		for _, e := range adj[u] {
			v := e.to
			if parent[v] != u {
				continue
			}
			cost := add(c[v], add(e.w1, e.w2))
			if cost < c[u] {
				c[u] = cost
			}
		}
	}
	// top-down
	for _, u := range order {
		for _, e := range adj[u] {
			v := e.to
			if parent[v] != u {
				continue
			}
			cost := add(c[u], add(e.w1, e.w2))
			if cost < c[v] {
				c[v] = cost
			}
		}
	}

	// Lifting tables.
	anc := make([][LOG]int, n+1)
	up := make([][LOG]Matrix, n+1)   // from node to ancestor
	down := make([][LOG]Matrix, n+1) // from ancestor to node

	for i := 1; i <= n; i++ {
		for j := 0; j < LOG; j++ {
			up[i][j] = newMatrixInf()
			down[i][j] = newMatrixInf()
		}
	}

	for _, u := range order {
		anc[u][0] = parent[u]
		if parent[u] != 0 {
			up[u][0] = edgeMatrix(u, parent[u], pw1[u], pw2[u], c)
			down[u][0] = edgeMatrix(parent[u], u, pw1[u], pw2[u], c)
		}
	}

	for j := 1; j < LOG; j++ {
		for u := 1; u <= n; u++ {
			mid := anc[u][j-1]
			anc[u][j] = anc[mid][j-1]
			if anc[u][j] != 0 {
				up[u][j] = multiply(up[u][j-1], up[mid][j-1])
				down[u][j] = multiply(down[anc[u][j-1]][j-1], down[u][j-1])
			}
		}
	}

	q := nextInt()
	for ; q > 0; q-- {
		uRaw := nextInt()
		vRaw := nextInt()

		if uRaw == vRaw {
			fmt.Fprintln(out, 0)
			continue
		}

		tu := 0
		if uRaw%2 == 0 {
			tu = 1
		}
		tv := 0
		if vRaw%2 == 0 {
			tv = 1
		}

		u := (uRaw + 1) / 2
		v := (vRaw + 1) / 2

		if depth[u] < depth[v] {
			u, v = v, u
			tu, tv = tv, tu
			uRaw, vRaw = vRaw, uRaw
		}

		startNode := u
		endNode := v

		matUp := newMatrixIdentity()
		matDown := newMatrixIdentity()

		// lift u to depth v
		diff := depth[u] - depth[v]
		for j := LOG - 1; j >= 0; j-- {
			if diff&(1<<j) != 0 {
				matUp = multiply(matUp, up[u][j])
				u = anc[u][j]
			}
		}

		// move both up until LCA
		if u != v {
			for j := LOG - 1; j >= 0; j-- {
				if anc[u][j] != 0 && anc[u][j] != anc[v][j] {
					matUp = multiply(matUp, up[u][j])
					matDown = multiply(down[v][j], matDown)
					u = anc[u][j]
					v = anc[v][j]
				}
			}
			// final step to LCA
			matUp = multiply(matUp, up[u][0])
			matDown = multiply(down[v][0], matDown)
			u = anc[u][0]
			v = anc[v][0]
		}

		// combined path matrix start -> end
		pathMat := multiply(matUp, matDown)

		// start vector allows switching at start immediately
		s0, s1 := INF, INF
		if tu == 0 {
			s0 = 0
			s1 = c[startNode]
		} else {
			s1 = 0
			s0 = c[startNode]
		}

		res0 := min64(add(s0, pathMat.g[0][0]), add(s1, pathMat.g[1][0]))
		res1 := min64(add(s0, pathMat.g[0][1]), add(s1, pathMat.g[1][1]))

		// If we end in the wrong parity we can switch at the destination at cost c[endNode].
		ans := res0
		if tv == 1 {
			ans = res1
			alt := add(res0, c[endNode])
			if alt < ans {
				ans = alt
			}
		} else {
			alt := add(res1, c[endNode])
			if alt < ans {
				ans = alt
			}
		}

		fmt.Fprintln(out, ans)
	}
}
