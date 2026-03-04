package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct{ u, v int }

type testCase struct {
	n     int
	edges []edge
}

func normalize(u, v int) edge {
	if u > v {
		u, v = v, u
	}
	return edge{u, v}
}

func minP(tc testCase) int {
	deg := make([]int, tc.n+1)
	for _, e := range tc.edges {
		if e.u == e.v {
			deg[e.u] += 2
		} else {
			deg[e.u]++
			deg[e.v]++
		}
	}
	odd := 0
	for i := 1; i <= tc.n; i++ {
		if deg[i]%2 != 0 {
			odd++
		}
	}
	p := len(tc.edges) + odd/2
	if p%2 != 0 {
		p++
	}
	return p
}

func validateOutput(tc testCase, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return errors.New("empty output")
	}
	p64, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil || p64 < 0 {
		return fmt.Errorf("invalid p: %q", fields[0])
	}
	p := int(p64)
	if len(fields) != 1+2*p {
		return fmt.Errorf("expected %d edge endpoints after p, got %d tokens", 2*p, len(fields)-1)
	}

	if p != minP(tc) {
		return fmt.Errorf("non-minimal p: got %d, expected %d", p, minP(tc))
	}

	need := make(map[edge]int)
	for _, e := range tc.edges {
		need[normalize(e.u, e.v)]++
	}
	have := make(map[edge]int)
	inDeg := make([]int, tc.n+1)
	outDeg := make([]int, tc.n+1)

	pos := 1
	for i := 0; i < p; i++ {
		u64, errU := strconv.ParseInt(fields[pos], 10, 64)
		v64, errV := strconv.ParseInt(fields[pos+1], 10, 64)
		pos += 2
		if errU != nil || errV != nil {
			return fmt.Errorf("invalid edge on line %d", i+2)
		}
		u, v := int(u64), int(v64)
		if u < 1 || u > tc.n || v < 1 || v > tc.n {
			return fmt.Errorf("edge endpoint out of range on line %d: %d %d", i+2, u, v)
		}
		outDeg[u]++
		inDeg[v]++
		have[normalize(u, v)]++
	}

	for e, cnt := range need {
		if have[e] < cnt {
			return fmt.Errorf("missing original edge multiplicity for (%d,%d): need %d, got %d", e.u, e.v, cnt, have[e])
		}
	}

	for v := 1; v <= tc.n; v++ {
		if inDeg[v]%2 != 0 || outDeg[v]%2 != 0 {
			return fmt.Errorf("parity failed at vertex %d: in=%d out=%d", v, inDeg[v], outDeg[v])
		}
	}
	return nil
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(6) + 2
		edges := make([]edge, 0, 16)
		// Start from a random tree so the graph is connected.
		for v := 2; v <= n; v++ {
			u := rnd.Intn(v-1) + 1
			edges = append(edges, edge{u, v})
		}
		extra := rnd.Intn(8) // extra edges: loops/multi-edges allowed by statement
		for i := 0; i < extra; i++ {
			u := rnd.Intn(n) + 1
			v := rnd.Intn(n) + 1
			edges = append(edges, edge{u, v})
		}
		tests = append(tests, testCase{n: n, edges: edges})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
		input := sb.String()
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if err := validateOutput(tc, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s got:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
