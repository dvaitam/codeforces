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
	t int // 1 directed, 0 undirected
	u int
	v int
}

type testCase struct {
	n, m     int
	edges    []edge
	input    string
	possible bool
}

func topoPossible(n int, edges []edge) bool {
	deg := make([]int, n+1)
	adj := make([][]int, n+1)
	for _, e := range edges {
		if e.t == 1 {
			adj[e.u] = append(adj[e.u], e.v)
			deg[e.v]++
		}
	}
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if deg[i] == 0 {
			q = append(q, i)
		}
	}
	cnt := 0
	for i := 0; i < len(q); i++ {
		v := q[i]
		cnt++
		for _, to := range adj[v] {
			deg[to]--
			if deg[to] == 0 {
				q = append(q, to)
			}
		}
	}
	return cnt == n
}

func buildCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2
	var edges []edge
	possible := true
	if rng.Float64() < 0.3 {
		// create directed cycle
		if n < 3 {
			n = 3
		}
		a, b, c := 1, 2, 3
		edges = append(edges, edge{1, a, b}, edge{1, b, c}, edge{1, c, a})
		m := rng.Intn(3)
		for i := 0; i < m; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				v = (v % n) + 1
			}
			edges = append(edges, edge{0, u, v})
		}
		possible = false
	} else {
		perm := rng.Perm(n)
		for i := range perm {
			perm[i]++
		}
		m := rng.Intn(n*2) + n
		for i := 0; i < m; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for u == v {
				v = rng.Intn(n) + 1
			}
			if perm[u-1] < perm[v-1] {
				if rng.Intn(2) == 0 {
					edges = append(edges, edge{1, u, v})
				} else {
					edges = append(edges, edge{0, u, v})
				}
			} else {
				if rng.Intn(2) == 0 {
					edges = append(edges, edge{1, v, u})
				} else {
					edges = append(edges, edge{0, v, u})
				}
			}
		}
		possible = true
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.t, e.u, e.v))
	}
	return testCase{n: n, m: len(edges), edges: edges, input: sb.String(), possible: possible}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	first := strings.TrimSpace(lines[0])
	if first == "NO" {
		if tc.possible {
			return fmt.Errorf("expected YES got NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("expected YES or NO got %s", first)
	}
	if !tc.possible {
		return fmt.Errorf("expected NO but got YES")
	}
	if len(lines)-1 != tc.m {
		return fmt.Errorf("expected %d edges, got %d", tc.m, len(lines)-1)
	}
	orient := make([]edge, 0, tc.m)
	for _, line := range lines[1:] {
		var u, v int
		if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d %d", &u, &v); err != nil {
			return fmt.Errorf("bad edge line: %v", err)
		}
		orient = append(orient, edge{1, u, v})
	}
	// check edges correspond
	used := make([]bool, tc.m)
	for _, o := range orient {
		found := false
		for i, e := range tc.edges {
			if used[i] {
				continue
			}
			if e.t == 1 {
				if e.u == o.u && e.v == o.v {
					used[i] = true
					found = true
					break
				}
			} else {
				// undirected can be either direction
				if (e.u == o.u && e.v == o.v) || (e.u == o.v && e.v == o.u) {
					used[i] = true
					found = true
					break
				}
			}
		}
		if !found {
			return fmt.Errorf("output edge %d %d not in input", o.u, o.v)
		}
	}
	for i, u := range used {
		if !u {
			return fmt.Errorf("missing edge %d in output", i)
		}
	}
	// check acyclicity
	deg := make([]int, tc.n+1)
	adj := make([][]int, tc.n+1)
	for _, e := range orient {
		adj[e.u] = append(adj[e.u], e.v)
		deg[e.v]++
	}
	q := make([]int, 0)
	for i := 1; i <= tc.n; i++ {
		if deg[i] == 0 {
			q = append(q, i)
		}
	}
	cnt := 0
	for i := 0; i < len(q); i++ {
		v := q[i]
		cnt++
		for _, to := range adj[v] {
			deg[to]--
			if deg[to] == 0 {
				q = append(q, to)
			}
		}
	}
	if cnt != tc.n {
		return fmt.Errorf("output edges form a cycle")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0)
	for i := 0; i < 100; i++ {
		cases = append(cases, buildCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
