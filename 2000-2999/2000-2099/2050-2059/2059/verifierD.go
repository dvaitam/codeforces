package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	refSource2059D = "2059D.go"
	refBinary2059D = "ref2059D.bin"
	maxTotalN      = 900
	totalTests     = 70
)

type edge struct {
	u int
	v int
}

type pair struct {
	a int
	b int
}

type testCase struct {
	n  int
	s1 int
	s2 int
	e1 []edge
	e2 []edge
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on case %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2059D, refSource2059D)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2059D), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.s1, tc.s2)
		fmt.Fprintf(&sb, "%d\n", len(tc.e1))
		for _, e := range tc.e1 {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
		fmt.Fprintf(&sb, "%d\n", len(tc.e2))
		for _, e := range tc.e2 {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2059))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += tc.n
	}

	add(testCase{
		n:  2,
		s1: 1, s2: 2,
		e1: []edge{{1, 2}},
		e2: []edge{{1, 2}},
	})

	add(testCase{
		n:  3,
		s1: 1, s2: 1,
		e1: []edge{{1, 2}, {2, 3}},
		e2: []edge{{1, 3}, {3, 2}, {2, 1}},
	})

	for len(tests) < totalTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		maxN := 40
		if remain < maxN {
			maxN = remain
		}
		n := rnd.Intn(maxN-1) + 2
		s1 := rnd.Intn(n) + 1
		s2 := rnd.Intn(n) + 1
		e1 := randomGraph(n, rnd)
		e2 := randomGraph(n, rnd)
		add(testCase{n: n, s1: s1, s2: s2, e1: e1, e2: e2})
	}

	return tests
}

func randomGraph(n int, rnd *rand.Rand) []edge {
	edges := make(map[pair]struct{})
	// build a spanning tree
	for v := 2; v <= n; v++ {
		u := rnd.Intn(v-1) + 1
		addEdge(u, v, edges)
	}
	// add extra edges
	extra := rnd.Intn(n)
	attempts := 0
	for len(edges) < n-1+extra && attempts < extra*5+10 {
		u := rnd.Intn(n) + 1
		v := rnd.Intn(n) + 1
		if u == v {
			attempts++
			continue
		}
		addEdge(u, v, edges)
	}
	res := make([]edge, 0, len(edges))
	for p := range edges {
		res = append(res, edge{u: p.a, v: p.b})
	}
	return res
}

func addEdge(u, v int, edges map[pair]struct{}) {
	if u > v {
		u, v = v, u
	}
	key := pair{a: u, b: v}
	edges[key] = struct{}{}
}
