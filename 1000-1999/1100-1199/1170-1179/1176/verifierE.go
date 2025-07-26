package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type graphCase struct {
	n, m  int
	edges [][2]int
}

func genCase(r *rand.Rand) graphCase {
	n := r.Intn(8) + 2
	maxEdges := n * (n - 1) / 2
	m := r.Intn(maxEdges-(n-1)+1) + (n - 1)
	edges := make([][2]int, 0, m)
	// create tree first
	for i := 2; i <= n; i++ {
		p := r.Intn(i-1) + 1
		edges = append(edges, [2]int{i, p})
	}
	edgeSet := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		edgeSet[e] = true
	}
	for len(edges) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if edgeSet[key] {
			continue
		}
		edgeSet[key] = true
		edges = append(edges, [2]int{u, v})
	}
	return graphCase{n: n, m: m, edges: edges}
}

func (gc graphCase) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", gc.n, gc.m))
	for _, e := range gc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func check(tc graphCase, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	if k < 1 || k > tc.n/2 {
		return fmt.Errorf("k out of range")
	}
	if len(fields)-1 != k {
		return fmt.Errorf("expected %d vertices, got %d", k, len(fields)-1)
	}
	chosen := make([]bool, tc.n+1)
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil || v < 1 || v > tc.n {
			return fmt.Errorf("invalid vertex")
		}
		if chosen[v] {
			return fmt.Errorf("duplicate vertex")
		}
		chosen[v] = true
	}
	adj := make([][]int, tc.n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	for v := 1; v <= tc.n; v++ {
		if chosen[v] {
			continue
		}
		ok := false
		for _, u := range adj[v] {
			if chosen[u] {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("vertex %d not covered", v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		tc := genCase(rng)
		input := tc.Input()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if err := check(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
