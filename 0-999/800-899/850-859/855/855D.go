package main

import (
	"bufio"
	"fmt"
	"os"
)

// Fenwick tree supporting range add and point query
type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n + 2, tree: make([]int, n+3)}
}

func (f *Fenwick) add(idx, val int) {
	for i := idx + 1; i < len(f.tree); i += i & -i {
		f.tree[i] += val
	}
}

func (f *Fenwick) RangeAdd(l, r, val int) {
	f.add(l, val)
	f.add(r, -val)
}

func (f *Fenwick) Query(idx int) int {
	res := 0
	for i := idx + 1; i > 0; i -= i & -i {
		res += f.tree[i]
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	parent := make([]int, n+1)
	typ := make([]int, n+1)
	specialChildren := make([][]int, n+1)
	partChildren := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &parent[i], &typ[i])
		if parent[i] != -1 {
			if typ[i] == 0 {
				specialChildren[parent[i]] = append(specialChildren[parent[i]], i)
			} else {
				partChildren[parent[i]] = append(partChildren[parent[i]], i)
			}
		}
	}

	// prepare parent arrays for special and part trees
	specParent := make([]int, n+1)
	partParent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		specParent[i] = -1
		partParent[i] = -1
	}
	for i := 1; i <= n; i++ {
		if parent[i] != -1 {
			if typ[i] == 0 {
				specParent[i] = parent[i]
			} else {
				partParent[i] = parent[i]
			}
		}
	}

	// Euler order for special tree
	tinSpec := make([]int, n+1)
	toutSpec := make([]int, n+1)
	time := 0
	var dfsSpecEuler func(int)
	dfsSpecEuler = func(u int) {
		tinSpec[u] = time
		time++
		for _, v := range specialChildren[u] {
			dfsSpecEuler(v)
		}
		toutSpec[u] = time
	}
	for i := 1; i <= n; i++ {
		if specParent[i] == -1 {
			dfsSpecEuler(i)
		}
	}

	// Euler order for part tree
	tinPart := make([]int, n+1)
	toutPart := make([]int, n+1)
	time = 0
	var dfsPartEuler func(int)
	dfsPartEuler = func(u int) {
		tinPart[u] = time
		time++
		for _, v := range partChildren[u] {
			dfsPartEuler(v)
		}
		toutPart[u] = time
	}
	for i := 1; i <= n; i++ {
		if partParent[i] == -1 {
			dfsPartEuler(i)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	type Query struct{ t, u, v int }
	queries := make([]Query, q)
	queriesByU := make([][]struct{ idx, v int }, n+1)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i].t, &queries[i].u, &queries[i].v)
		if queries[i].t == 2 {
			u := queries[i].u
			queriesByU[u] = append(queriesByU[u], struct{ idx, v int }{i, queries[i].v})
		}
	}

	ans := make([]string, q)
	// answer type1 queries immediately
	for i, qu := range queries {
		if qu.t == 1 {
			if qu.u != qu.v && tinSpec[qu.u] <= tinSpec[qu.v] && tinSpec[qu.v] < toutSpec[qu.u] {
				ans[i] = "YES"
			} else {
				ans[i] = "NO"
			}
		}
	}

	// Fenwick for active part nodes
	fen := NewFenwick(n)

	var dfsSpec func(int)
	dfsSpec = func(u int) {
		for _, p := range partChildren[u] {
			fen.RangeAdd(tinPart[p], toutPart[p], 1)
		}
		for _, qinfo := range queriesByU[u] {
			if qinfo.v != u && fen.Query(tinPart[qinfo.v]) > 0 {
				ans[qinfo.idx] = "YES"
			} else {
				ans[qinfo.idx] = "NO"
			}
		}
		for _, v := range specialChildren[u] {
			dfsSpec(v)
		}
		for _, p := range partChildren[u] {
			fen.RangeAdd(tinPart[p], toutPart[p], -1)
		}
	}

	for i := 1; i <= n; i++ {
		if specParent[i] == -1 {
			dfsSpec(i)
		}
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, ans[i])
	}
}
