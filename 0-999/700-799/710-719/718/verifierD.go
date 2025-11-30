package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesD = `3
1 2
2 3
1
4
1 2
2 3
2 4
4
1 2
2 3
2 4
1
3
1 2
1 3
6
1 2
1 3
3 4
4 5
2 6
4
1 2
1 3
1 4
4
1 2
2 3
2 4
1
4
1 2
1 3
1 4
6
1 2
2 3
2 4
2 5
4 6
1
3
1 2
1 3
4
1 2
1 3
2 4
6
1 2
1 3
1 4
1 5
5 6
2
1 2
1
2
1 2
5
1 2
2 3
1 4
2 5
5
1 2
2 3
3 4
4 5
1
3
1 2
2 3
6
1 2
2 3
3 4
3 5
1 6
1
6
1 2
2 3
3 4
1 5
3 6
6
1 2
1 3
3 4
1 5
4 6
6
1 2
2 3
2 4
2 5
4 6
5
1 2
2 3
3 4
3 5
6
1 2
2 3
1 4
2 5
5 6
3
1 2
1 3
1
6
1 2
2 3
3 4
4 5
1 6
4
1 2
2 3
3 4
5
1 2
2 3
1 4
1 5
6
1 2
2 3
3 4
1 5
1 6
4
1 2
1 3
3 4
3
1 2
2 3
4
1 2
2 3
3 4
6
1 2
2 3
3 4
2 5
4 6
4
1 2
2 3
3 4
6
1 2
2 3
1 4
2 5
4 6
5
1 2
2 3
2 4
2 5
5
1 2
2 3
1 4
2 5
1
2
1 2
3
1 2
1 3
1
3
1 2
2 3
1
4
1 2
2 3
3 4
4
1 2
2 3
2 4
1
3
1 2
2 3
2
1 2
3
1 2
1 3
4
1 2
2 3
2 4
1
3
1 2
2 3
6
1 2
2 3
1 4
1 5
1 6
3
1 2
2 3
2
1 2
3
1 2
1 3
3
1 2
1 3
6
1 2
2 3
2 4
2 5
4 6
6
1 2
1 3
1 4
2 5
2 6
1
5
1 2
2 3
2 4
4 5
4
1 2
2 3
1 4
6
1 2
1 3
3 4
1 5
1 6
2
1 2
5
1 2
1 3
1 4
2 5
3
1 2
1 3
3
1 2
1 3
4
1 2
1 3
3 4
6
1 2
1 3
2 4
2 5
2 6
3
1 2
1 3
4
1 2
1 3
3 4
4
1 2
2 3
2 4
6
1 2
1 3
1 4
4 5
1 6
5
1 2
1 3
3 4
4 5
1
1
3
1 2
2 3
3
1 2
2 3
6
1 2
2 3
2 4
2 5
4 6
1
1
1
4
1 2
1 3
2 4
3
1 2
2 3
5
1 2
2 3
2 4
2 5
1
5
1 2
1 3
1 4
1 5
6
1 2
1 3
3 4
1 5
4 6
6
1 2
1 3
3 4
2 5
3 6
4
1 2
1 3
3 4
5
1 2
1 3
3 4
2 5
5
1 2
1 3
1 4
3 5
1`

const MOD = 1000000007

// DSU for 1..n
type DSU struct {
	p, r []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	r := make([]int, n+1)
	for i := range p {
		p[i] = i
	}
	return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(a, b int) {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return
	}
	if d.r[a] < d.r[b] {
		a, b = b, a
	}
	d.p[b] = a
	if d.r[a] == d.r[b] {
		d.r[a]++
	}
}

// Embedded solver from 718D.go.
func solve(n int, edges [][2]int) int {
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}

	parent := make([]int, n+1)
	sz := make([]int, n+1)
	order := make([]int, 0, n)
	type Item struct{ u, p, idx int }
	stack := []Item{{1, 0, 0}}
	parent[1] = 0
	for len(stack) > 0 {
		it := &stack[len(stack)-1]
		u, p := it.u, it.p
		if it.idx < len(adj[u]) {
			v := adj[u][it.idx]
			it.idx++
			if v == p {
				continue
			}
			parent[v] = u
			stack = append(stack, Item{v, u, 0})
		} else {
			order = append(order, u)
			stack = stack[:len(stack)-1]
		}
	}
	for _, u := range order {
		sz[u] = 1
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			sz[u] += sz[v]
		}
	}
	minSize := n + 1
	centroids := make([]int, 0, 2)
	for u := 1; u <= n; u++ {
		maxPart := n - sz[u]
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			if sz[v] > maxPart {
				maxPart = sz[v]
			}
		}
		if maxPart < minSize {
			minSize = maxPart
			centroids = centroids[:0]
			centroids = append(centroids, u)
		} else if maxPart == minSize {
			centroids = append(centroids, u)
		}
	}
	if len(centroids) == 2 {
		c1, c2 := centroids[0], centroids[1]
		for i, v := range adj[c1] {
			if v == c2 {
				adj[c1] = append(adj[c1][:i], adj[c1][i+1:]...)
				break
			}
		}
		for i, v := range adj[c2] {
			if v == c1 {
				adj[c2] = append(adj[c2][:i], adj[c2][i+1:]...)
				break
			}
		}
	}
	var root int
	children := make([][]int, n+1)
	if len(centroids) == 1 {
		root = centroids[0]
	} else {
		root = 0
		children = make([][]int, n+1)
	}
	st2 := []int{root}
	par2 := make([]int, len(children))
	par2[root] = -1
	for len(st2) > 0 {
		u := st2[len(st2)-1]
		st2 = st2[:len(st2)-1]
		if u == 0 {
			c1, c2 := centroids[0], centroids[1]
			children[u] = append(children[u], c1, c2)
			par2[c1] = 0
			par2[c2] = 0
			st2 = append(st2, c1, c2)
			continue
		}
		for _, v := range adj[u] {
			if v == par2[u] {
				continue
			}
			par2[v] = u
			children[u] = append(children[u], v)
			st2 = append(st2, v)
		}
	}
	order = []int{}
	type SI struct{ u, idx int }
	st3 := []SI{{root, 0}}
	for len(st3) > 0 {
		top := &st3[len(st3)-1]
		u := top.u
		if top.idx < len(children[u]) {
			v := children[u][top.idx]
			top.idx++
			st3 = append(st3, SI{v, 0})
		} else {
			order = append(order, u)
			st3 = st3[:len(st3)-1]
		}
	}
	subtype := make([]int, len(children))
	nextID := 1
	key2id := make(map[string]int)
	dsu := NewDSU(n)
	for _, u := range order {
		ids := make([]int, 0, len(children[u]))
		for _, v := range children[u] {
			ids = append(ids, subtype[v])
		}
		sort.Ints(ids)
		key := fmt.Sprint(ids)
		id, ok := key2id[key]
		if !ok {
			id = nextID
			key2id[key] = id
			nextID++
		}
		subtype[u] = id
		groups := make(map[int][]int)
		for _, v := range children[u] {
			cid := subtype[v]
			groups[cid] = append(groups[cid], v)
		}
		for _, vs := range groups {
			if len(vs) <= 1 {
				continue
			}
			root0 := vs[0]
			for i := 1; i < len(vs); i++ {
				queue := [][2]int{{root0, vs[i]}}
				for len(queue) > 0 {
					pair := queue[len(queue)-1]
					queue = queue[:len(queue)-1]
					a, b := pair[0], pair[1]
					dsu.Union(a, b)
					ca := children[a]
					cb := children[b]
					la := len(ca)
					aList := make([]int, la)
					copy(aList, ca)
					bList := make([]int, la)
					copy(bList, cb)
					sort.Slice(aList, func(i, j int) bool { return subtype[aList[i]] < subtype[aList[j]] })
					sort.Slice(bList, func(i, j int) bool { return subtype[bList[i]] < subtype[bList[j]] })
					for k := 0; k < la; k++ {
						queue = append(queue, [2]int{aList[k], bList[k]})
					}
				}
			}
		}
	}
	seen := make(map[int]bool)
	res := 0
	for i := 1; i <= n; i++ {
		if len(adj[i]) < 4 {
			leader := dsu.Find(i)
			if !seen[leader] {
				seen[leader] = true
				res++
			}
		}
	}
	return res % MOD
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesD)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	idx := 0
	readInt := func() (int, error) {
		if idx >= len(fields) {
			return 0, fmt.Errorf("unexpected end")
		}
		v, err := strconv.Atoi(fields[idx])
		idx++
		return v, err
	}
	cases := make([]testCase, 0)
	for idx < len(fields) {
		n, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n", len(cases)+1)
		}
		edges := make([][2]int, 0, n-1)
		for j := 0; j < n-1; j++ {
			u, err1 := readInt()
			v, err2 := readInt()
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("case %d: bad edge %d", len(cases)+1, j+1)
			}
			edges = append(edges, [2]int{u, v})
		}
		cases = append(cases, testCase{n: n, edges: edges})
	}
	return cases, nil
}

type testCase struct {
	n     int
	edges [][2]int
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.n, tc.edges)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		gotStr, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
