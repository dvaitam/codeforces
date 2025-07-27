package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, q   int
	adj    [][]int
	parent [][]int
	depth  []int
	size   []int
	heavy  []int
	head   []int
	pos    []int
	curPos int
	bit    *BIT
	val    []int64
	absVal []int64
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func dfs1(v, p int) {
	parent[0][v] = p
	depth[v] = depth[p] + 1
	size[v] = 1
	heavy[v] = 0
	maxSz := 0
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		dfs1(to, v)
		if size[to] > maxSz {
			maxSz = size[to]
			heavy[v] = to
		}
		size[v] += size[to]
	}
}

func dfs2(v, h int) {
	head[v] = h
	pos[v] = curPos
	curPos++
	if heavy[v] != 0 {
		dfs2(heavy[v], h)
	}
	for _, to := range adj[v] {
		if to != parent[0][v] && to != heavy[v] {
			dfs2(to, to)
		}
	}
}

func pathSum(u, v int) int64 {
	var res int64
	for head[u] != head[v] {
		if depth[head[u]] > depth[head[v]] {
			u, v = v, u
		}
		res += bit.Query(pos[head[v]], pos[v])
		v = parent[0][head[v]]
	}
	if depth[u] > depth[v] {
		u, v = v, u
	}
	res += bit.Query(pos[u], pos[v])
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	val = make([]int64, n+1)
	absVal = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &val[i])
		absVal[i] = abs64(val[i])
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	LOG := 0
	for (1 << LOG) <= n {
		LOG++
	}
	parent = make([][]int, LOG)
	for i := range parent {
		parent[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	size = make([]int, n+1)
	heavy = make([]int, n+1)
	head = make([]int, n+1)
	pos = make([]int, n+1)

	dfs1(1, 0)
	curPos = 1
	dfs2(1, 1)
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			parent[k][v] = parent[k-1][parent[k-1][v]]
		}
	}

	bit = NewBIT(n)
	for i := 1; i <= n; i++ {
		bit.Add(pos[i], absVal[i])
	}

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var u int
			var c int64
			fmt.Fscan(in, &u, &c)
			newAbs := abs64(c)
			diff := newAbs - absVal[u]
			if diff != 0 {
				bit.Add(pos[u], diff)
				absVal[u] = newAbs
				val[u] = c
			} else {
				val[u] = c
			}
		} else if t == 2 {
			var u, v int
			fmt.Fscan(in, &u, &v)
			sum := pathSum(u, v)
			ans := 2*sum - absVal[u] - absVal[v]
			fmt.Fprintln(out, ans)
		}
	}
}

type BIT struct {
	n int
	t []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n, make([]int64, n+1)}
}

func (b *BIT) Add(i int, v int64) {
	for x := i; x <= b.n; x += x & -x {
		b.t[x] += v
	}
}

func (b *BIT) Query(l, r int) int64 {
	return b.sum(r) - b.sum(l-1)
}

func (b *BIT) sum(i int) int64 {
	var s int64
	for x := i; x > 0; x -= x & -x {
		s += b.t[x]
	}
	return s
}
