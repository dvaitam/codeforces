package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod = 998244353

type bitset []uint64

func newBitset(n int) bitset {
	return make(bitset, (n+63)>>6)
}

func (b bitset) set(i int) {
	b[i>>6] |= 1 << (i & 63)
}

func (b bitset) get(i int) bool {
	return (b[i>>6]>>(i&63))&1 == 1
}

func (b bitset) clone() bitset {
	c := make(bitset, len(b))
	copy(c, b)
	return c
}

func bfs(n int, adj []bitset, flip1, flip2 int, start int) bitset {
	vis := newBitset(n)
	q := make([]int, 0, n)
	vis.set(start)
	q = append(q, start)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		neigh := adj[v].clone()
		if v == flip1 || v == flip2 {
			for i := range neigh {
				neigh[i] = ^neigh[i]
			}
			neigh[v>>6] &^= 1 << (v & 63)
		}
		if flip1 >= 0 && flip1 != v {
			w := flip1 >> 6
			b := uint64(1) << (flip1 & 63)
			if (adj[v][w] & b) != 0 {
				neigh[w] &^= b
			} else {
				neigh[w] |= b
			}
		}
		if flip2 >= 0 && flip2 != v {
			w := flip2 >> 6
			b := uint64(1) << (flip2 & 63)
			if (adj[v][w] & b) != 0 {
				neigh[w] &^= b
			} else {
				neigh[w] |= b
			}
		}
		for i, w := range neigh {
			w &^= vis[i]
			for w != 0 {
				b := w & -w
				idx := bits.TrailingZeros64(w)
				to := i<<6 + idx
				if to < n {
					vis.set(to)
					q = append(q, to)
				}
				w &^= b
			}
		}
	}
	return vis
}

func stronglyConnected(n int, adj, radj []bitset, flip1, flip2 int) bool {
	vis := bfs(n, adj, flip1, flip2, 0)
	for i := 0; i < n; i++ {
		if !vis.get(i) {
			return false
		}
	}
	vis = bfs(n, radj, flip1, flip2, 0)
	for i := 0; i < n; i++ {
		if !vis.get(i) {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	g := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &g[i])
	}
	adj := make([]bitset, n)
	radj := make([]bitset, n)
	for i := 0; i < n; i++ {
		adj[i] = newBitset(n)
		radj[i] = newBitset(n)
	}
	outdeg := make([]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if g[i][j] == '1' {
				adj[i].set(j)
				radj[j].set(i)
				outdeg[i]++
			}
		}
	}
	for i := 0; i < n; i++ {
		if outdeg[i] == 0 || outdeg[i] == n-1 {
			fmt.Fprintln(out, -1)
			return
		}
	}
	if stronglyConnected(n, adj, radj, -1, -1) {
		fmt.Fprintln(out, 0, 1)
		return
	}

	// Kosaraju to find SCCs
	order := make([]int, 0, n)
	vis := make([]bool, n)
	var dfs1 func(int)
	dfs1 = func(v int) {
		vis[v] = true
		for to := 0; to < n; to++ {
			if adj[v].get(to) && !vis[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < n; i++ {
		if !vis[i] {
			dfs1(i)
		}
	}
	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	var dfs2 func(int, int)
	dfs2 = func(v, c int) {
		comp[v] = c
		for to := 0; to < n; to++ {
			if radj[v].get(to) && comp[to] == -1 {
				dfs2(to, c)
			}
		}
	}
	cnum := 0
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == -1 {
			dfs2(v, cnum)
			cnum++
		}
	}
	if cnum >= 3 {
		cnt := 0
		for i := 0; i < n; i++ {
			if comp[i] != 0 && comp[i] != cnum-1 {
				cnt++
			}
		}
		fmt.Fprintln(out, 1, cnt%mod)
		return
	}
	// two components
	c1, c2 := 0, 1
	size1, size2 := 0, 0
	for i := 0; i < n; i++ {
		if comp[i] == c1 {
			size1++
		} else {
			size2++
		}
	}
	count1 := 0
	for i := 0; i < n; i++ {
		if comp[i] == c1 {
			if stronglyConnected(n, adj, radj, i, -1) {
				count1++
			}
		}
	}
	for i := 0; i < n; i++ {
		if comp[i] == c2 {
			if stronglyConnected(n, adj, radj, i, -1) {
				count1++
			}
		}
	}
	if count1 > 0 {
		fmt.Fprintln(out, 1, count1%mod)
		return
	}
	ways := (int64(size1) * int64(size2) * 2) % mod
	fmt.Fprintln(out, 2, ways)
}
