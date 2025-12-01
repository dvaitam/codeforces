package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type edge struct {
	u int
	v int
}

type query struct {
	a int
	b int
}

type testCase struct {
	name    string
	n       int
	edges   []edge
	queries []query
	input   string
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, q int) ([]int, error) {
	res := make([]int, 0, q)
	reader := strings.NewReader(out)
	for len(res) < q {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return nil, fmt.Errorf("expected %d integers, got %d: %v", q, len(res), err)
		}
		res = append(res, x)
	}
	var extra int
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected after %d numbers", q)
	}
	return res, nil
}

func addCliqueEdges(edges map[[2]int]struct{}, nodes []int) {
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			u, v := nodes[i], nodes[j]
			if u > v {
				u, v = v, u
			}
			edges[[2]int{u, v}] = struct{}{}
		}
	}
}

func buildRandomTest(rng *rand.Rand, idx int) testCase {
	targetN := rng.Intn(40) + 2 // 2..41 nodes
	edgesMap := make(map[[2]int]struct{})
	current := 0
	firstSize := rng.Intn(4) + 1
	if firstSize > targetN {
		firstSize = targetN
	}
	initial := make([]int, firstSize)
	for i := 0; i < firstSize; i++ {
		current++
		initial[i] = current
	}
	addCliqueEdges(edgesMap, initial)
	for current < targetN {
		parent := rng.Intn(current) + 1
		remaining := targetN - current
		addSize := rng.Intn(4) + 1
		if addSize > remaining {
			addSize = remaining
		}
		block := []int{parent}
		for i := 0; i < addSize; i++ {
			current++
			block = append(block, current)
		}
		addCliqueEdges(edgesMap, block)
	}
	n := current
	edges := make([]edge, 0, len(edgesMap))
	for k := range edgesMap {
		edges = append(edges, edge{u: k[0], v: k[1]})
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].u == edges[j].u {
			return edges[i].v < edges[j].v
		}
		return edges[i].u < edges[j].u
	})
	m := len(edges)
	qCount := rng.Intn(25) + 1
	queries := make([]query, qCount)
	for i := 0; i < qCount; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		for a == b {
			b = rng.Intn(n) + 1
		}
		if a > b {
			a, b = b, a
		}
		queries[i] = query{a: a, b: b}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, qCount)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for _, qu := range queries {
		fmt.Fprintf(&sb, "%d %d\n", qu.a, qu.b)
	}
	return testCase{
		name:    fmt.Sprintf("rand_%d", idx),
		n:       n,
		edges:   edges,
		queries: queries,
		input:   sb.String(),
	}
}

func buildTreeTest() testCase {
	n := 5
	edges := []edge{{1, 2}, {2, 3}, {3, 4}, {4, 5}}
	queries := []query{{1, 5}, {2, 4}, {1, 3}}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, len(edges), len(queries))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for _, qu := range queries {
		fmt.Fprintf(&sb, "%d %d\n", qu.a, qu.b)
	}
	return testCase{
		name:    "tree",
		n:       n,
		edges:   edges,
		queries: queries,
		input:   sb.String(),
	}
}

func buildCliqueTest() testCase {
	n := 4
	edges := []edge{
		{1, 2}, {1, 3}, {1, 4},
		{2, 3}, {2, 4},
		{3, 4},
	}
	queries := []query{{1, 2}, {1, 4}, {2, 3}}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, len(edges), len(queries))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for _, qu := range queries {
		fmt.Fprintf(&sb, "%d %d\n", qu.a, qu.b)
	}
	return testCase{
		name:    "clique",
		n:       n,
		edges:   edges,
		queries: queries,
		input:   sb.String(),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCase(n int, edges []edge, queries []query) []int {
	m := len(edges)
	e1 := make([][]int, n+1)
	for _, e := range edges {
		e1[e.u] = append(e1[e.u], e.v)
		e1[e.v] = append(e1[e.v], e.u)
	}
	maxNodes := n + m + 5
	e2 := make([][]int, maxNodes)
	dfn := make([]int, n+1)
	low := make([]int, n+1)
	sta := make([]int, n+1)
	val := make([]int, maxNodes)
	cnt := make([]int, maxNodes)
	fa := make([]int, maxNodes)
	depth := make([]int, maxNodes)
	topNode := make([]int, maxNodes)
	son := make([]int, maxNodes)
	size := make([]int, maxNodes)

	idx := 0
	stackTop := 0
	tot := n

	var tarjan func(u int)
	tarjan = func(u int) {
		idx++
		dfn[u] = idx
		low[u] = idx
		stackTop++
		sta[stackTop] = u
		for _, v := range e1[u] {
			if dfn[v] == 0 {
				tarjan(v)
				low[u] = min(low[u], low[v])
				if low[v] >= dfn[u] {
					tot++
					val[tot] = 1
					for {
						x := sta[stackTop]
						stackTop--
						e2[tot] = append(e2[tot], x)
						e2[x] = append(e2[x], tot)
						if x == v {
							break
						}
					}
					e2[tot] = append(e2[tot], u)
					e2[u] = append(e2[u], tot)
				}
			} else {
				low[u] = min(low[u], dfn[v])
			}
		}
	}

	var dfs1 func(u, parent int)
	dfs1 = func(u, parent int) {
		size[u] = 1
		for _, v := range e2[u] {
			if v == parent {
				continue
			}
			depth[v] = depth[u] + 1
			cnt[v] = cnt[u] + val[v]
			fa[v] = u
			dfs1(v, u)
			size[u] += size[v]
			if son[u] == 0 || size[v] > size[son[u]] {
				son[u] = v
			}
		}
	}

	var dfs2 func(u, top int)
	dfs2 = func(u, top int) {
		topNode[u] = top
		if son[u] != 0 {
			dfs2(son[u], top)
		}
		for _, v := range e2[u] {
			if v == fa[u] || v == son[u] {
				continue
			}
			dfs2(v, v)
		}
	}

	var lca func(u, v int) int
	lca = func(u, v int) int {
		for topNode[u] != topNode[v] {
			if depth[topNode[u]] > depth[topNode[v]] {
				u, v = v, u
			}
			v = fa[topNode[v]]
		}
		if depth[u] > depth[v] {
			return v
		}
		return u
	}

	queryDist := func(u, v int) int {
		l := lca(u, v)
		return depth[u] + depth[v] - 2*depth[l] - cnt[u] - cnt[v] + cnt[l] + cnt[fa[l]]
	}

	for i := 1; i <= n; i++ {
		if dfn[i] == 0 {
			tarjan(i)
			cnt[i] = val[i]
			dfs1(i, 0)
			dfs2(i, i)
		}
	}

	results := make([]int, len(queries))
	for i, qu := range queries {
		results[i] = queryDist(qu.a, qu.b)
	}
	return results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []testCase{
		buildTreeTest(),
		buildCliqueTest(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, buildRandomTest(rng, i+1))
	}
	for idx, tc := range tests {
		expected := solveCase(tc.n, tc.edges, tc.queries)
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out, len(tc.queries))
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		for i := range expected {
			if got[i] != expected[i] {
				fmt.Printf("test %d (%s) mismatch on query %d: expect %d got %d\ninput:\n%soutput:\n%s\n", idx+1, tc.name, i+1, expected[i], got[i], tc.input, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
