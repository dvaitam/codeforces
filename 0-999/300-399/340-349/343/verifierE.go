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

type edge struct {
	u, v, c int
}

type testCase struct {
	n     int
	edges []edge
}

func addEdge(g [][]edgeImpl, u, v, c int) {
	g[u] = append(g[u], edgeImpl{v, len(g[v]), c})
	g[v] = append(g[v], edgeImpl{u, len(g[u]) - 1, c})
}

type edgeImpl struct {
	to, rev, cap int
}

func maxFlow(n int, edges []edge, s, t int) int {
	g := make([][]edgeImpl, n+1)
	for _, e := range edges {
		addEdge(g, e.u, e.v, e.c)
	}
	flow := 0
	for {
		q := []int{s}
		parent := make([]int, n+1)
		pe := make([]int, n+1)
		for i := range parent {
			parent[i] = -1
		}
		parent[s] = s
		for len(q) > 0 && parent[t] == -1 {
			v := q[0]
			q = q[1:]
			for i, e := range g[v] {
				if e.cap > 0 && parent[e.to] == -1 {
					parent[e.to] = v
					pe[e.to] = i
					q = append(q, e.to)
					if e.to == t {
						break
					}
				}
			}
		}
		if parent[t] == -1 {
			break
		}
		aug := 1<<31 - 1
		for v := t; v != s; v = parent[v] {
			e := &g[parent[v]][pe[v]]
			if e.cap < aug {
				aug = e.cap
			}
		}
		for v := t; v != s; v = parent[v] {
			e := &g[parent[v]][pe[v]]
			e.cap -= aug
			rev := e.rev
			g[v][rev].cap += aug
		}
		flow += aug
	}
	return flow
}

func bestSalary(n int, edges []edge) (int, []int) {
	flows := make([][]int, n+1)
	for i := range flows {
		flows[i] = make([]int, n+1)
	}
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			flows[i][j] = maxFlow(n, edges, i, j)
			flows[j][i] = flows[i][j]
		}
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	best := -1
	bestPerm := make([]int, n)
	var dfs func(int)
	dfs = func(idx int) {
		if idx == n {
			sum := 0
			for i := 0; i < n-1; i++ {
				sum += flows[perm[i]][perm[i+1]]
			}
			if sum > best {
				best = sum
				copy(bestPerm, perm)
			}
			return
		}
		for i := idx; i < n; i++ {
			perm[idx], perm[i] = perm[i], perm[idx]
			dfs(idx + 1)
			perm[idx], perm[i] = perm[i], perm[idx]
		}
	}
	dfs(0)
	return best, bestPerm
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.c)
	}
	input := sb.String()
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) < 1+tc.n {
		return fmt.Errorf("not enough output lines")
	}
	var salary int
	if _, err := fmt.Sscan(fields[0], &salary); err != nil {
		return fmt.Errorf("bad salary: %v", err)
	}
	perm := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		if _, err := fmt.Sscan(fields[i+1], &perm[i]); err != nil {
			return fmt.Errorf("bad permutation value: %v", err)
		}
	}
	expSalary, _ := bestSalary(tc.n, tc.edges)
	if salary != expSalary {
		return fmt.Errorf("expected salary %d got %d", expSalary, salary)
	}
	used := make([]bool, tc.n+1)
	for _, v := range perm {
		if v < 1 || v > tc.n || used[v] {
			return fmt.Errorf("invalid permutation")
		}
		used[v] = true
	}
	return nil
}

func generateCases(rng *rand.Rand) []testCase {
	cases := []testCase{}
	for len(cases) < 100 {
		n := rng.Intn(5) + 2
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges-n+1) + n - 1
		// start with a tree to ensure connectivity
		edges := make([]edge, 0, m)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			c := rng.Intn(10) + 1
			edges = append(edges, edge{p, i, c})
		}
		added := n - 1
		for added < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			exists := false
			for _, e := range edges {
				if e.u == u && e.v == v || e.u == v && e.v == u {
					exists = true
					break
				}
			}
			if exists {
				continue
			}
			c := rng.Intn(10) + 1
			edges = append(edges, edge{u, v, c})
			added++
		}
		cases = append(cases, testCase{n, edges})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateCases(rng)
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
