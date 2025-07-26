package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type edge struct{ u, v int }

type testCase struct {
	n     int
	edges []edge
}

func solveB(tc testCase) int64 {
	n := tc.n
	m := len(tc.edges)
	d := newDSU(n)
	deg := make([]int64, n+1)
	has := make([]bool, n+1)
	var loops int64
	for _, e := range tc.edges {
		u, v := e.u, e.v
		if u == v {
			loops++
			has[u] = true
		} else {
			deg[u]++
			deg[v]++
			d.union(u, v)
			has[u] = true
			has[v] = true
		}
	}
	root := -1
	for i := 1; i <= n; i++ {
		if has[i] {
			if root == -1 {
				root = d.find(i)
			} else if d.find(i) != root {
				return 0
			}
		}
	}
	if m < 2 {
		return 0
	}
	ans := loops*(loops-1)/2 + loops*(int64(m)-loops)
	for i := 1; i <= n; i++ {
		if deg[i] >= 2 {
			ans += deg[i] * (deg[i] - 1) / 2
		}
	}
	return ans
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	maxEdges := n * n
	m := rng.Intn(maxEdges + 1)
	edges := make([]edge, 0, m)
	used := make(map[string]bool)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		key := fmt.Sprintf("%d-%d", u, v)
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, edge{u, v})
	}
	return testCase{n: n, edges: edges}
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func expected(tc testCase) string {
	return fmt.Sprintf("%d", solveB(tc))
}

func runCase(bin string, tc testCase) error {
	input := formatInput(tc)
	exp := expected(tc)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if out != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s\ninput:\n%s", exp, out, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// some fixed edge cases
	cases = append(cases, testCase{n: 2, edges: []edge{{1, 1}, {2, 2}}})
	cases = append(cases, testCase{n: 3, edges: []edge{{1, 2}, {2, 3}}})
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
