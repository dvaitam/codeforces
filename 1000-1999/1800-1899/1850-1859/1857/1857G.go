package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

// DSU structure for union-find operations
type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func powMod(a, e int64) int64 {
	a %= mod
	if a < 0 {
		a += mod
	}
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var S int64
		fmt.Fscan(reader, &n, &S)

		type Edge struct {
			u, v int
			w    int64
		}
		edges := make([]Edge, n-1)
		freq := make(map[int64]int64)
		for i := 0; i < n-1; i++ {
			var u, v int
			var w int64
			fmt.Fscan(reader, &u, &v, &w)
			u--
			v--
			edges[i] = Edge{u, v, w}
			freq[w]++
		}

		sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
		dsu := NewDSU(n)
		counts := make(map[int64]int64)

		for _, e := range edges {
			u := dsu.find(e.u)
			v := dsu.find(e.v)
			if u != v {
				counts[e.w] += int64(dsu.size[u]) * int64(dsu.size[v])
				dsu.union(u, v)
			}
		}

		ans := int64(1)
		zero := false
		for w, c := range counts {
			exp := c - freq[w]
			if exp <= 0 {
				continue
			}
			diff := S - w
			if diff <= 0 {
				zero = true
				break
			}
			ans = ans * powMod(diff%mod, exp) % mod
		}
		if zero {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, ans%mod)
		}
	}
}
