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

const mod int64 = 998244353

type edge struct {
	u, v int
	p, q int64
}

type testCase struct {
	n, m int
	e    []edge
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	ref, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for i, tc := range tests {
		input := serialize(tc)

		refOut, err := runProgram(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\noutput:\n%s", i+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", i+1, err, candOut)
			os.Exit(1)
		}

		if refAns != candAns {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\ninput:\n%s", i+1, refAns, candAns, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine current path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2029H-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2029H")
	cmd := exec.Command("go", "build", "-o", outPath, "2029H.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func serialize(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.m * 40)
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.e {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", e.u, e.v, e.p, e.q))
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 0 || val >= mod {
		return 0, fmt.Errorf("value out of range: %d", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		// Sample 1
		{
			n: 2, m: 1,
			e: []edge{{u: 1, v: 2, p: 1, q: 10}},
		},
		// Sample 2
		{
			n: 3, m: 3,
			e: []edge{
				{u: 1, v: 2, p: 1, q: 2},
				{u: 1, v: 3, p: 1, q: 2},
				{u: 2, v: 3, p: 1, q: 2},
			},
		},
		// Sample 3
		{
			n: 1, m: 0,
		},
		// Sample 4
		{
			n: 5, m: 8,
			e: []edge{
				{u: 1, v: 2, p: 1, q: 11},
				{u: 1, v: 3, p: 2, q: 11},
				{u: 1, v: 4, p: 3, q: 11},
				{u: 1, v: 5, p: 4, q: 11},
				{u: 2, v: 4, p: 5, q: 11},
				{u: 2, v: 5, p: 6, q: 11},
				{u: 3, v: 4, p: 7, q: 11},
				{u: 4, v: 5, p: 8, q: 11},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		tc, ok := buildRandomTest(rng)
		if !ok {
			continue
		}
		tests = append(tests, tc)
	}
	return tests
}

func buildRandomTest(rng *rand.Rand) (testCase, bool) {
	n := rng.Intn(10) + 2 // ensure at least 2 vertices for spread dynamics
	if n > 12 {
		n = 12
	}
	maxEdges := n*(n-1)/2
	m := rng.Intn(maxEdges-n+2) + n - 1 // at least tree

	edges := make([]edge, 0, m)
	parent := randTree(rng, n)
	for v := 2; v <= n; v++ {
		edges = append(edges, randomEdge(rng, parent[v], v))
	}

	existing := make(map[[2]int]struct{})
	for v := 2; v <= n; v++ {
		a := parent[v]
		if a > v {
			a, v = v, a
		}
		existing[[2]int{a, v}] = struct{}{}
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
		if _, ok := existing[[2]int{u, v}]; ok {
			continue
		}
		edges = append(edges, randomEdge(rng, u, v))
		existing[[2]int{u, v}] = struct{}{}
	}

	tc := testCase{n: n, m: m, e: edges}
	if !satisfiesConstraint(tc) {
		return testCase{}, false
	}
	return tc, true
}

func randTree(rng *rand.Rand, n int) []int {
	parent := make([]int, n+1)
	for v := 2; v <= n; v++ {
		parent[v] = rng.Intn(v-1) + 1
	}
	return parent
}

func randomEdge(rng *rand.Rand, u, v int) edge {
	q := int64(rng.Intn(500000) + 2)
	if q >= mod {
		q = mod - 2
	}
	var p int64
	for {
		p = int64(rng.Intn(int(q-1)) + 1)
		if gcd(p, q) == 1 {
			break
		}
	}
	return edge{u: u, v: v, p: p, q: q}
}

func satisfiesConstraint(tc testCase) bool {
	if tc.n == 1 {
		return true
	}
	g := make([][]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		g[i] = make([]int64, tc.n)
	}
	for _, e := range tc.e {
		prob := e.p * modPow(e.q, mod-2) % mod
		g[e.u-1][e.v-1] = prob
		g[e.v-1][e.u-1] = prob
	}
	// Check non-empty proper subsets.
	for mask := 1; mask < (1<<tc.n)-1; mask++ {
		prod := int64(1)
		for i := 0; i < tc.n; i++ {
			if (mask>>i)&1 == 0 {
				continue
			}
			for j := 0; j < tc.n; j++ {
				if (mask>>j)&1 != 0 {
					continue
				}
				prod = prod * ((1 - g[i][j] + mod) % mod) % mod
			}
		}
		if prod == 1 {
			return false
		}
	}
	return true
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
