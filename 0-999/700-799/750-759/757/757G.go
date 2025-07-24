package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxLog = 20

type Edge struct {
	to int
	w  int64
}

var (
	n, q  int
	seq   []int
	tree  [][]Edge
	depth []int64
	up    [][]int
	dist  []int64
)

func dfs(u, p int, d int64) {
	up[0][u] = p
	depth[u] = d
	for _, e := range tree[u] {
		if e.to == p {
			continue
		}
		dfs(e.to, u, d+e.w)
	}
}

func buildLCA() {
	up = make([][]int, maxLog)
	for i := range up {
		up[i] = make([]int, n)
	}
	depth = make([]int64, n)
	dfs(0, 0, 0)
	for k := 1; k < maxLog; k++ {
		for i := 0; i < n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := 0; diff > 0; k++ {
		if diff&1 == 1 {
			a = up[k][a]
		}
		diff >>= 1
	}
	if a == b {
		return a
	}
	for k := maxLog - 1; k >= 0; k-- {
		if up[k][a] != up[k][b] {
			a = up[k][a]
			b = up[k][b]
		}
	}
	return up[0][a]
}

func distAB(a, b int) int64 {
	l := lca(a, b)
	return depth[a] + depth[b] - 2*depth[l]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &q)
	seq = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &seq[i])
		seq[i]--
	}
	tree = make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		var w int64
		fmt.Fscan(reader, &u, &v, &w)
		u--
		v--
		tree[u] = append(tree[u], Edge{v, w})
		tree[v] = append(tree[v], Edge{u, w})
	}
	buildLCA()

	ansPrev := int64(0)
	mask := int64((1 << 30) - 1)
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var a, b, c int64
			fmt.Fscan(reader, &a, &b, &c)
			l := int((ansPrev & mask) ^ a)
			r := int((ansPrev & mask) ^ b)
			v := int((ansPrev & mask) ^ c)
			l--
			r--
			v--
			res := int64(0)
			for j := l; j <= r; j++ {
				res += distAB(seq[j], v)
			}
			fmt.Fprintln(writer, res)
			ansPrev = res
		} else {
			var a int64
			fmt.Fscan(reader, &a)
			x := int((ansPrev & mask) ^ a)
			x--
			if x >= 0 && x+1 < n {
				seq[x], seq[x+1] = seq[x+1], seq[x]
			}
		}
	}
}
