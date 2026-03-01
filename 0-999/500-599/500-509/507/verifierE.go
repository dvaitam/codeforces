package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type edge struct{ u, v, z int }

// Test case structure
type TestE struct {
	n     int
	m     int
	edges []edge
}

func generateTests() []TestE {
	rand.Seed(5)
	tests := make([]TestE, 100)
	for i := range tests {
		n := rand.Intn(6) + 2 // 2..7
		maxExtra := n*(n-1)/2 - (n - 1)
		extra := rand.Intn(maxExtra + 1)
		m := (n - 1) + extra
		edges := make([]edge, 0, m)
		// create tree to ensure connectivity
		for v := 2; v <= n; v++ {
			u := rand.Intn(v-1) + 1
			z := rand.Intn(2)
			edges = append(edges, edge{u, v, z})
		}
		// add extra edges
		pairs := make([][2]int, 0)
		for a := 1; a <= n; a++ {
			for b := a + 1; b <= n; b++ {
				skip := false
				for _, e := range edges {
					if (e.u == a && e.v == b) || (e.u == b && e.v == a) {
						skip = true
						break
					}
				}
				if !skip {
					pairs = append(pairs, [2]int{a, b})
				}
			}
		}
		rand.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
		for i2 := 0; i2 < extra && i2 < len(pairs); i2++ {
			p := pairs[i2]
			z := rand.Intn(2)
			edges = append(edges, edge{p[0], p[1], z})
		}
		tests[i] = TestE{n, len(edges), edges}
	}
	return tests
}

// Implementation of reference algorithm from 507E.go
const inf = int64(1e18)
const cFijo = int64(1000000)

type Edge struct {
	to int
	z  int
}

type Item struct {
	node int
	dist int64
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func expected(t TestE) string {
	N := t.n
	M := t.m
	g := make([][]Edge, N)
	edges := make([]edge, 0, M)
	rotos := 0
	for _, e := range t.edges {
		u, v, z := e.u-1, e.v-1, e.z
		g[u] = append(g[u], Edge{v, z})
		g[v] = append(g[v], Edge{u, z})
		edges = append(edges, e)
		if z == 0 {
			rotos++
		}
	}
	dist := make([]int64, N)
	parent := make([]int, N)
	for i := 0; i < N; i++ {
		dist[i] = inf
		parent[i] = -1
	}
	dist[0] = 0
	pq := &PQ{}
	heap.Init(pq)
	heap.Push(pq, Item{0, 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		u := it.node
		d := it.dist
		if d != dist[u] {
			continue
		}
		for _, e := range g[u] {
			cost := cFijo
			if e.z == 0 {
				cost++
			}
			nd := d + cost
			v := e.to
			if nd < dist[v] {
				dist[v] = nd
				parent[v] = u
				heap.Push(pq, Item{v, nd})
			}
		}
	}
	pathMap := make(map[[2]int]bool)
	zMap := make(map[[2]int]int)
	for _, e := range edges {
		u, v := e.u-1, e.v-1
		key := [2]int{min(u, v), max(u, v)}
		zMap[key] = e.z
	}
	brokenInPath := 0
	pathLen := 0
	cur := N - 1
	for parent[cur] != -1 {
		p := parent[cur]
		key := [2]int{min(p, cur), max(p, cur)}
		pathMap[key] = true
		if zMap[key] == 0 {
			brokenInPath++
		}
		pathLen++
		cur = p
	}
	ops := M - rotos - pathLen + 2*brokenInPath
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d\n", ops))
	for _, e := range edges {
		u, v, z := e.u-1, e.v-1, e.z
		key := [2]int{min(u, v), max(u, v)}
		if pathMap[key] {
			if z == 0 {
				buf.WriteString(fmt.Sprintf("%d %d 1\n", u+1, v+1))
			}
		} else {
			if z == 1 {
				buf.WriteString(fmt.Sprintf("%d %d 0\n", u+1, v+1))
			}
		}
	}
	return strings.TrimSpace(buf.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildInput(t TestE) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.z))
	}
	return sb.String()
}

func normalizeOutput(out string) (string, error) {
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return "", fmt.Errorf("empty output")
	}

	lines := strings.Split(trimmed, "\n")
	var k int
	if _, err := fmt.Sscanf(strings.TrimSpace(lines[0]), "%d", &k); err != nil {
		return "", fmt.Errorf("invalid operation count line: %w", err)
	}

	ops := make([]string, 0, len(lines)-1)
	for idx, raw := range lines[1:] {
		var x, y, op int
		if _, err := fmt.Sscanf(strings.TrimSpace(raw), "%d %d %d", &x, &y, &op); err != nil {
			return "", fmt.Errorf("invalid operation line %d: %w", idx+2, err)
		}
		if x > y {
			x, y = y, x
		}
		ops = append(ops, fmt.Sprintf("%d %d %d", x, y, op))
	}

	if len(ops) != k {
		return "", fmt.Errorf("declared %d operations but found %d", k, len(ops))
	}

	sort.Strings(ops)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", k))
	for i, op := range ops {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(op)
	}
	return strings.TrimSpace(b.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := buildInput(t)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		expRaw := expected(t)
		exp, err := normalizeOutput(expRaw)
		if err != nil {
			fmt.Printf("verifier bug on expected output at test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotNorm, err := normalizeOutput(got)
		if err != nil {
			fmt.Printf("test %d failed: invalid contestant output: %v\nraw output:\n%s\n", i+1, err, strings.TrimSpace(got))
			os.Exit(1)
		}

		if gotNorm != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, exp, gotNorm)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
