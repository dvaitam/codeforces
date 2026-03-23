package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ---------- embedded solver (from cf_t23_240_E.go) ----------

type SolverEdge struct {
	id, u, v, w int
}

type HeapNode struct {
	id, u, v, w int
	lazy        int
	left, right *HeapNode
}

func heapApply(n *HeapNode, lazy int) {
	if n != nil {
		n.w += lazy
		n.lazy += lazy
	}
}

func heapPush(n *HeapNode) {
	if n != nil && n.lazy != 0 {
		heapApply(n.left, n.lazy)
		heapApply(n.right, n.lazy)
		n.lazy = 0
	}
}

func heapMerge(a, b *HeapNode) *HeapNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	heapPush(a)
	heapPush(b)
	if a.w > b.w {
		a, b = b, a
	}
	a.right = heapMerge(a.right, b)
	a.left, a.right = a.right, a.left
	return a
}

func heapPop(n *HeapNode) *HeapNode {
	heapPush(n)
	res := heapMerge(n.left, n.right)
	n.left = nil
	n.right = nil
	return res
}

func solveEmbedded(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024*10)

	out := &bytes.Buffer{}

	scanInt := func() int {
		scanner.Scan()
		x, _ := strconv.Atoi(scanner.Text())
		return x
	}

	if !scanner.Scan() {
		return ""
	}
	n, _ := strconv.Atoi(scanner.Text())
	m := scanInt()

	edges := make([]SolverEdge, m+1)
	adj := make([][]SolverEdge, n+1)

	for i := 1; i <= m; i++ {
		u := scanInt()
		v := scanInt()
		w := scanInt()
		edges[i] = SolverEdge{i, u, v, w}
		adj[u] = append(adj[u], edges[i])
	}

	reached := make([]bool, n+1)
	q := []int{1}
	reached[1] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, e := range adj[u] {
			if !reached[e.v] {
				reached[e.v] = true
				q = append(q, e.v)
			}
		}
	}

	for i := 1; i <= n; i++ {
		if !reached[i] {
			fmt.Fprintln(out, "-1")
			return strings.TrimSpace(out.String())
		}
	}

	maxNodes := 2*n + 1
	heaps := make([]*HeapNode, maxNodes)
	for i := 1; i <= m; i++ {
		e := edges[i]
		if e.v == 1 {
			continue
		}
		nNode := &HeapNode{id: e.id, u: e.u, v: e.v, w: e.w}
		heaps[e.v] = heapMerge(heaps[e.v], nNode)
	}

	parent := make([]int, maxNodes)
	for i := 1; i < maxNodes; i++ {
		parent[i] = i
	}

	find := func(i int) int {
		root := i
		for root != parent[root] {
			root = parent[root]
		}
		curr := i
		for curr != root {
			nxt := parent[curr]
			parent[curr] = root
			curr = nxt
		}
		return root
	}

	visited := make([]int, maxNodes)
	visited[1] = 2

	enter := make([]*HeapNode, maxNodes)
	inEdge := make([]*HeapNode, n+1)

	compId := n

	for i := 2; i <= n; i++ {
		if visited[i] != 0 {
			continue
		}

		curr := i
		path := []int{curr}
		visited[curr] = 1

		for {
			var minEdge *HeapNode
			for heaps[curr] != nil {
				cand := heaps[curr]
				heaps[curr] = heapPop(heaps[curr])
				if find(cand.u) == curr {
					continue
				}
				minEdge = cand
				break
			}

			if minEdge == nil {
				break
			}

			enter[curr] = minEdge
			inEdge[minEdge.v] = minEdge

			nxt := find(minEdge.u)

			if visited[nxt] == 2 {
				break
			} else if visited[nxt] == 0 {
				visited[nxt] = 1
				path = append(path, nxt)
				curr = nxt
			} else if visited[nxt] == 1 {
				idx := -1
				for j := 0; j < len(path); j++ {
					if path[j] == nxt {
						idx = j
						break
					}
				}

				compId++
				newC := compId

				for j := idx; j < len(path); j++ {
					c := path[j]
					parent[c] = newC
					lazyVal := -enter[c].w
					if heaps[c] != nil {
						heapApply(heaps[c], lazyVal)
						heaps[newC] = heapMerge(heaps[newC], heaps[c])
					}
				}

				path = path[:idx]
				path = append(path, newC)
				visited[newC] = 1
				curr = newC
			}
		}

		for _, c := range path {
			visited[c] = 2
		}
	}

	ans := []int{}
	for i := 2; i <= n; i++ {
		if inEdge[i] != nil && edges[inEdge[i].id].w == 1 {
			ans = append(ans, inEdge[i].id)
		}
	}

	if len(ans) == 0 {
		fmt.Fprintln(out, 0)
	} else {
		fmt.Fprintln(out, len(ans))
		for i, id := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, id)
		}
		fmt.Fprintln(out)
	}
	return strings.TrimSpace(out.String())
}

// ---------- verifier infrastructure ----------

type edgeData struct {
	from, to   int
	needRepair bool
}

type verTestCase struct {
	input  string
	expect int
	n      int
	edges  []edgeData
}

type refEdge struct {
	from, to int
	cost     int
}

const inf = int(1e9)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkAnswer(tc, out); err != nil {
			// Also run embedded solver for comparison
			refOut := solveEmbedded(tc.input)
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\nexpected minimal repairs: %d\nreference output:\n%s\n", i+1, err, tc.input, out, tc.expect, refOut)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkAnswer(tc verTestCase, output string) error {
	output = strings.TrimSpace(output)
	if len(output) == 0 {
		return fmt.Errorf("empty output")
	}
	tokens := strings.Fields(output)
	if tokens[0] == "-1" {
		if tc.expect != -1 {
			return fmt.Errorf("expected %d repairs but got -1", tc.expect)
		}
		if len(tokens) != 1 {
			return fmt.Errorf("unexpected extra tokens after -1")
		}
		return nil
	}
	if tc.expect == -1 {
		return fmt.Errorf("expected -1 but got %s", tokens[0])
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil || k < 0 {
		return fmt.Errorf("invalid number of repairs %q", tokens[0])
	}
	if k == 0 {
		if tc.expect != 0 {
			return fmt.Errorf("expected %d repairs but reported 0", tc.expect)
		}
		// "0" can appear alone or with "0" on second line
		if len(tokens) != 1 {
			return fmt.Errorf("extra tokens after 0")
		}
		return nil
	}
	if len(tokens) < k+1 {
		return fmt.Errorf("reported %d roads but provided %d identifiers", k, len(tokens)-1)
	}
	if k != tc.expect {
		return fmt.Errorf("expected %d repairs but reported %d", tc.expect, k)
	}
	selected := make(map[int]struct{}, k)
	for i := 0; i < k; i++ {
		id, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return fmt.Errorf("invalid road index %q", tokens[i+1])
		}
		if id < 1 || id > len(tc.edges) {
			return fmt.Errorf("road index %d out of range", id)
		}
		if _, ok := selected[id]; ok {
			return fmt.Errorf("road index %d listed multiple times", id)
		}
		if !tc.edges[id-1].needRepair {
			return fmt.Errorf("road %d is already good and cannot be repaired", id)
		}
		selected[id] = struct{}{}
	}
	if !allReachable(tc.n, tc.edges, selected) {
		return fmt.Errorf("not all cities reachable from the capital with reported repairs")
	}
	return nil
}

func allReachable(n int, edges []edgeData, repaired map[int]struct{}) bool {
	adj := make([][]int, n+1)
	for idx, e := range edges {
		if !e.needRepair {
			adj[e.from] = append(adj[e.from], e.to)
			continue
		}
		if _, ok := repaired[idx+1]; ok {
			adj[e.from] = append(adj[e.from], e.to)
		}
	}
	vis := make([]bool, n+1)
	queue := []int{1}
	vis[1] = true
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				queue = append(queue, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return false
		}
	}
	return true
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), err
}

func genTests() []verTestCase {
	rand.Seed(42)
	var tests []verTestCase
	tests = append(tests, newTestCase(1, nil))
	tests = append(tests, newTestCase(2, []edgeData{
		{from: 1, to: 2, needRepair: false},
	}))
	tests = append(tests, newTestCase(2, nil))
	tests = append(tests, newTestCase(3, []edgeData{
		{from: 1, to: 2, needRepair: true},
		{from: 2, to: 3, needRepair: false},
	}))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTestCase())
	}
	return tests
}

func randomTestCase() verTestCase {
	n := rand.Intn(10) + 1
	maxPossible := n * (n - 1)
	limit := n*3 + rand.Intn(n+1)
	if limit > maxPossible {
		limit = maxPossible
	}
	m := 0
	if limit > 0 {
		m = rand.Intn(limit + 1)
	}
	edges := make([]edgeData, 0, m)
	used := make(map[int]struct{})
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		key := (u-1)*n + (v - 1)
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		edges = append(edges, edgeData{
			from:       u,
			to:         v,
			needRepair: rand.Intn(2) == 1,
		})
	}
	return newTestCase(n, edges)
}

func newTestCase(n int, edges []edgeData) verTestCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		c := 0
		if e.needRepair {
			c = 1
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.from, e.to, c))
	}
	edgesCopy := make([]edgeData, len(edges))
	copy(edgesCopy, edges)
	expect := calcMinRepairs(n, edgesCopy)
	return verTestCase{
		input:  sb.String(),
		expect: expect,
		n:      n,
		edges:  edgesCopy,
	}
}

func calcMinRepairs(n int, edges []edgeData) int {
	if n == 0 {
		return 0
	}
	refEdges := make([]refEdge, 0, len(edges))
	for _, e := range edges {
		c := 0
		if e.needRepair {
			c = 1
		}
		refEdges = append(refEdges, refEdge{
			from: e.from - 1,
			to:   e.to - 1,
			cost: c,
		})
	}
	cost, ok := directedMST(n, 0, refEdges)
	if !ok {
		return -1
	}
	return cost
}

func directedMST(n, root int, edges []refEdge) (int, bool) {
	total := 0
	for {
		in := make([]int, n)
		pre := make([]int, n)
		id := make([]int, n)
		vis := make([]int, n)
		for i := range in {
			in[i] = inf
			id[i] = -1
			vis[i] = -1
		}
		for _, e := range edges {
			if e.from == e.to {
				continue
			}
			if e.cost < in[e.to] {
				in[e.to] = e.cost
				pre[e.to] = e.from
			}
		}
		in[root] = 0
		for v := 0; v < n; v++ {
			if v == root {
				continue
			}
			if in[v] == inf {
				return 0, false
			}
		}
		cnt := 0
		for v := 0; v < n; v++ {
			total += in[v]
			u := v
			for vis[u] != v && id[u] == -1 && u != root {
				vis[u] = v
				u = pre[u]
			}
			if u != root && id[u] == -1 {
				for x := pre[u]; x != u; x = pre[x] {
					id[x] = cnt
				}
				id[u] = cnt
				cnt++
			}
		}
		if cnt == 0 {
			return total, true
		}
		for v := 0; v < n; v++ {
			if id[v] == -1 {
				id[v] = cnt
				cnt++
			}
		}
		newEdges := make([]refEdge, 0, len(edges))
		for _, e := range edges {
			u := id[e.from]
			v := id[e.to]
			if u != v {
				newEdges = append(newEdges, refEdge{
					from: u,
					to:   v,
					cost: e.cost - in[e.to],
				})
			}
		}
		root = id[root]
		n = cnt
		edges = newEdges
	}
}
