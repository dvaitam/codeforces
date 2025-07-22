package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// expectedDanger computes 2^(n-components)
func expectedDanger(n int, edges [][2]int) uint64 {
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		ra := find(a)
		rb := find(b)
		if ra != rb {
			parent[rb] = ra
		}
	}
	for _, e := range edges {
		union(e[0], e[1])
	}
	roots := make(map[int]struct{})
	for i := 1; i <= n; i++ {
		roots[find(i)] = struct{}{}
	}
	comps := len(roots)
	exp := n - comps
	return 1 << exp
}

func genCase(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(10) + 1 // up to 10 nodes
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	seen := make(map[[2]int]struct{})
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		edges = append(edges, key)
	}
	return n, edges
}

func caseToString(n int, edges [][2]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(2))

	type test struct {
		n     int
		edges [][2]int
	}
	var cases []test
	// fixed small cases
	cases = append(cases, test{n: 1})
	cases = append(cases, test{n: 2, edges: [][2]int{{1, 2}}})
	cases = append(cases, test{n: 3})
	cases = append(cases, test{n: 4, edges: [][2]int{{1, 2}, {3, 4}}})

	for len(cases) < 100 {
		n, edges := genCase(rng)
		cases = append(cases, test{n: n, edges: edges})
	}

	for idx, c := range cases {
		tcStr := caseToString(c.n, c.edges)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tcStr)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		var got uint64
		if _, err := fmt.Sscan(outStr, &got); err != nil {
			fmt.Printf("case %d: cannot parse output %q\n", idx+1, outStr)
			os.Exit(1)
		}
		exp := expectedDanger(c.n, c.edges)
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\ninput:\n%s", idx+1, exp, got, tcStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
