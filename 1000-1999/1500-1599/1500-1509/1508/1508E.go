package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	n       int
	a       []int
	g       [][]int
	sz      []int
	mn      []int
	mx      []int
	orig    []int
	invalid bool
)

func dfs(u int) {
	mn[u] = a[u]
	mx[u] = a[u]
	sz[u] = 1
	for _, v := range g[u] {
		dfs(v)
		if mn[v] < mn[u] {
			mn[u] = mn[v]
		}
		if mx[v] > mx[u] {
			mx[u] = mx[v]
		}
		sz[u] += sz[v]
	}
	if mx[u]-mn[u]+1 != sz[u] {
		invalid = true
	}
}

func build(u int, start int) int {
	orig[u] = start
	start++
	sort.Slice(g[u], func(i, j int) bool {
		return mn[g[u][i]] < mn[g[u][j]]
	})
	for _, v := range g[u] {
		start = build(v, start)
	}
	return start
}

// Fenwick tree implementation

type BIT struct {
	n int
	t []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, t: make([]int, n+2)}
}

func (b *BIT) Add(idx, val int) {
	for idx <= b.n {
		b.t[idx] += val
		idx += idx & -idx
	}
}

func (b *BIT) Sum(idx int) int {
	s := 0
	for idx > 0 {
		s += b.t[idx]
		idx -= idx & -idx
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	g = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
	}

	sz = make([]int, n)
	mn = make([]int, n)
	mx = make([]int, n)
	orig = make([]int, n)

	dfs(0)
	if invalid || mn[0] != 1 || mx[0] != n {
		fmt.Fprintln(out, "NO")
		return
	}
	build(0, 1)

	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		pos[a[i]] = i + 1
	}
	bit := NewBIT(n)
	inv := 0
	for i := 1; i <= n; i++ {
		p := pos[i]
		inv += i - 1 - bit.Sum(p)
		bit.Add(p, 1)
	}

	fmt.Fprintln(out, "YES")
	fmt.Fprintln(out, inv)
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, orig[i])
	}
	fmt.Fprintln(out)
}
