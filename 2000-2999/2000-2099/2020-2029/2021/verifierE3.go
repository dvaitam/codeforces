package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u, v int
	w    int64
}

type testCase struct {
	n, m, p int
	req     []int
	edges   []edge
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE3.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n%s", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}

	if err := compareOutputs(refOut, candOut, tests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2021E3-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2021E3")
	cmd := exec.Command("go", "build", "-o", outPath, "2021E3.go")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxTotalN = 200000
	const maxTotalM = 200000
	var tests []testCase
	totalN, totalM := 0, 0

	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN || totalM+tc.m > maxTotalM {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalM += tc.m
	}

	// Small deterministic cases
	add(buildPathTest(2, []edge{{1, 2, 5}}, []int{1}))
	add(buildPathTest(3, []edge{{1, 2, 1}, {2, 3, 2}}, []int{1, 3}))
	add(buildPathTest(4, []edge{{1, 2, 7}, {2, 3, 3}, {3, 4, 4}}, []int{2, 4}))

	// Random connected graphs
	for len(tests) < 40 && totalN < 180000 && totalM < 180000 {
		n := rng.Intn(3000) + 2
		maxExtra := n
		m := n - 1 + rng.Intn(maxExtra+1)
		p := rng.Intn(n) + 1
		req := randSample(rng, n, p)
		edges := generateConnectedGraph(rng, n, m)
		add(testCase{n: n, m: m, p: p, req: req, edges: edges})
	}

	// One larger stress case if budget allows
	if totalN+15000 <= maxTotalN && totalM+20000 <= maxTotalM {
		n := 15000
		m := 20000
		p := n/3 + 1
		req := randSample(rng, n, p)
		edges := generateConnectedGraph(rng, n, m)
		add(testCase{n: n, m: m, p: p, req: req, edges: edges})
	}

	return tests
}

func buildPathTest(n int, edges []edge, req []int) testCase {
	return testCase{
		n:     n,
		m:     len(edges),
		p:     len(req),
		req:   append([]int(nil), req...),
		edges: append([]edge(nil), edges...),
	}
}

func generateConnectedGraph(rng *rand.Rand, n, m int) []edge {
	edges := make([]edge, 0, m)
	// Start with a tree (chain) to ensure connectivity
	for i := 1; i < n; i++ {
		edges = append(edges, edge{u: i, v: i + 1, w: rngWeight(rng)})
	}
	seen := make(map[int]struct{}, m)
	for _, e := range edges {
		key := e.u*n + e.v
		seen[key] = struct{}{}
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := u*n + v
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		edges = append(edges, edge{u: u, v: v, w: rngWeight(rng)})
	}
	return edges
}

func rngWeight(rng *rand.Rand) int64 {
	return int64(rng.Intn(1_000_000_000) + 1)
}

func randSample(rng *rand.Rand, n, k int) []int {
	perm := rng.Perm(n)
	res := make([]int, k)
	for i := 0; i < k; i++ {
		res[i] = perm[i] + 1
	}
	return res
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.m, tc.p)
		for i, v := range tc.req {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&b, "%d %d %d\n", e.u, e.v, e.w)
		}
	}
	return b.String()
}

func compareOutputs(refOut, candOut string, tests []testCase) error {
	refTokens := strings.Fields(refOut)
	candTokens := strings.Fields(candOut)

	total := 0
	for _, tc := range tests {
		total += tc.n
	}
	if len(refTokens) != total {
		return fmt.Errorf("reference produced %d tokens, expected %d", len(refTokens), total)
	}
	if len(candTokens) != total {
		return fmt.Errorf("candidate produced %d tokens, expected %d", len(candTokens), total)
	}

	for i := 0; i < total; i++ {
		refVal, err := strconv.ParseInt(refTokens[i], 10, 64)
		if err != nil {
			return fmt.Errorf("reference token %q is not integer", refTokens[i])
		}
		candVal, err := strconv.ParseInt(candTokens[i], 10, 64)
		if err != nil {
			return fmt.Errorf("candidate token %q is not integer", candTokens[i])
		}
		if refVal != candVal {
			testIdx, pos := locateToken(tests, i)
			return fmt.Errorf("mismatch at test %d, position %d: expected %d got %d", testIdx+1, pos+1, refVal, candVal)
		}
	}
	return nil
}

func locateToken(tests []testCase, idx int) (testIndex int, pos int) {
	rem := idx
	for i, tc := range tests {
		if rem < tc.n {
			return i, rem
		}
		rem -= tc.n
	}
	return len(tests) - 1, tests[len(tests)-1].n - 1
}
