package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// refSource2109D points to the local reference solution to avoid GOPATH resolution.
const refSource2109D = "2109D.go"

type testCase struct {
	n     int
	m     int
	l     int
	a     []int
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	input := formatInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %q, got %q\nInput used:\n%s", i+1, expected[i], got[i], stringifyCase(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2109D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2109D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2109D)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.Write(errBuf.Bytes())
	}
	return out.String(), err
}

func parseOutputs(out string, tests []testCase) ([]string, error) {
	res := make([]string, len(tests))
	r := bufio.NewScanner(strings.NewReader(out))
	for i, tc := range tests {
		if !r.Scan() {
			return nil, fmt.Errorf("expected %d lines, got %d", len(tests), i)
		}
		line := strings.TrimSpace(r.Text())
		if len(line) != tc.n {
			return nil, fmt.Errorf("test %d: expected string of length %d, got %d (%q)", i+1, tc.n, len(line), line)
		}
		for idx, ch := range line {
			if ch != '0' && ch != '1' {
				return nil, fmt.Errorf("test %d: invalid character at position %d: %q", i+1, idx+1, ch)
			}
		}
		res[i] = line
	}
	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.l)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
	}
	return []byte(sb.String())
}

func stringifyCase(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.l)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 40)

	// Deterministic small cases to cover parity and bipartite/non-bipartite behaviors.
	tests = append(tests,
		// Simple edge, length-1 walk reaches node 2.
		makeLineCase(2, []int{1}),
		// Path length 3, only even distance reachable.
		makePathCase(3, []int{2}),
		// Triangle (non-bipartite), limited total sum.
		makeTriangleCase([]int{1}),
		// Square (bipartite) with only odd length walk.
		makeCycleCase(4, []int{3}),
		// Bigger odd/even mix.
		makeCycleCase(6, []int{5, 2}),
	)

	// Random connected graphs with varying multiset compositions.
	rng := rand.New(rand.NewSource(2109_2025))
	for len(tests) < 35 {
		n := rng.Intn(12) + 3   // 3..14 to keep outputs manageable
		maxL := 8 + rng.Intn(5) // 8..12 elements
		l := 1 + rng.Intn(maxL) // 1..maxL
		a := make([]int, l)
		for i := 0; i < l; i++ {
			a[i] = 1 + rng.Intn(20) // lengths up to 20 keep total reasonable
		}
		edges := generateConnectedGraph(n, rng)
		tc := testCase{
			n:     n,
			m:     len(edges),
			l:     l,
			a:     a,
			edges: edges,
		}
		tests = append(tests, tc)
	}

	return tests
}

func makeLineCase(n int, a []int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return testCase{n: n, m: len(edges), l: len(a), a: a, edges: edges}
}

func makePathCase(n int, a []int) testCase {
	return makeLineCase(n, a)
}

func makeTriangleCase(a []int) testCase {
	edges := [][2]int{{1, 2}, {2, 3}, {1, 3}}
	return testCase{n: 3, m: 3, l: len(a), a: a, edges: edges}
}

func makeCycleCase(n int, a []int) testCase {
	edges := make([][2]int, 0, n)
	for i := 1; i <= n; i++ {
		j := i + 1
		if j > n {
			j = 1
		}
		edges = append(edges, [2]int{i, j})
	}
	return testCase{n: n, m: len(edges), l: len(a), a: a, edges: edges}
}

func generateConnectedGraph(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n+5)
	// Start with a random spanning tree.
	for v := 2; v <= n; v++ {
		u := 1 + rng.Intn(v-1)
		edges = append(edges, normEdge(u, v))
	}
	edgeSet := make(map[[2]int]struct{}, len(edges))
	for _, e := range edges {
		edgeSet[e] = struct{}{}
	}
	// Add extra edges.
	extra := rng.Intn(n) // up to n-1 extra edges
	for extra > 0 {
		u := 1 + rng.Intn(n)
		v := 1 + rng.Intn(n)
		if u == v {
			continue
		}
		e := normEdge(u, v)
		if _, ok := edgeSet[e]; ok {
			continue
		}
		edgeSet[e] = struct{}{}
		edges = append(edges, e)
		extra--
	}
	return edges
}

func normEdge(u, v int) [2]int {
	if u > v {
		u, v = v, u
	}
	return [2]int{u, v}
}
