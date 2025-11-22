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
	refSource2117G = "2117G.go"
	refBinary2117G = "ref2117G.bin"
	maxTests       = 160
	maxTotalN      = 180000
	maxTotalM      = 180000
)

type edge struct {
	u, v int
	w    int64
}

type testCase struct {
	n, m  int
	edges []edge
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2117G, refSource2117G)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2117G), nil
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
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2117))
	var tests []testCase
	totalN, totalM := 0, 0

	add := func(tc testCase) {
		if len(tests) >= maxTests {
			return
		}
		if totalN+tc.n > maxTotalN || totalM+tc.m > maxTotalM {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalM += tc.m
	}

	// Deterministic small graphs
	add(simplePath(2, []int64{1}))
	add(simplePath(3, []int64{1, 3}))
	add(simplePath(3, []int64{5, 5}))
	add(simpleCycle(4, []int64{2, 7, 3, 9}))

	for len(tests) < maxTests && totalN < maxTotalN && totalM < maxTotalM {
		n := rnd.Intn(400) + 2
		remainN := maxTotalN - totalN
		if n > remainN {
			n = remainN
		}
		if n < 2 {
			break
		}

		// Build a random connected graph: start with a tree, then add extra edges.
		maxPossibleM := n * (n - 1) / 2
		targetM := rnd.Intn(minInt(maxPossibleM, 3*n)) + (n - 1)
		remainM := maxTotalM - totalM
		if targetM > remainM {
			targetM = remainM
		}
		if targetM < n-1 {
			targetM = n - 1
		}

		edges := make([]edge, 0, targetM)
		parents := randTree(n, rnd)
		used := make(map[int64]struct{}, targetM)
		for v := 1; v < n; v++ {
			u := parents[v]
			w := randWeight(rnd)
			edges = append(edges, edge{u: u + 1, v: v + 1, w: w})
			key := pairKey(u, v)
			used[key] = struct{}{}
		}

		for len(edges) < targetM {
			u := rnd.Intn(n)
			v := rnd.Intn(n)
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := pairKey(u, v)
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			edges = append(edges, edge{u: u + 1, v: v + 1, w: randWeight(rnd)})
		}

		add(testCase{n: n, m: len(edges), edges: edges})
	}
	return tests
}

func simplePath(n int, weights []int64) testCase {
	edges := make([]edge, 0, len(weights))
	for i, w := range weights {
		edges = append(edges, edge{u: i + 1, v: i + 2, w: w})
	}
	return testCase{n: n, m: len(edges), edges: edges}
}

func simpleCycle(n int, weights []int64) testCase {
	edges := make([]edge, 0, n)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		w := weights[i%len(weights)]
		edges = append(edges, edge{u: i + 1, v: j + 1, w: w})
	}
	return testCase{n: n, m: len(edges), edges: edges}
}

func randTree(n int, rnd *rand.Rand) []int {
	par := make([]int, n)
	par[0] = -1
	for v := 1; v < n; v++ {
		par[v] = rnd.Intn(v)
	}
	return par
}

func randWeight(rnd *rand.Rand) int64 {
	x := rnd.Int63n(1_000_000_000) + 1
	return x
}

func pairKey(u, v int) int64 {
	return int64(u)*1_000_000 + int64(v)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
