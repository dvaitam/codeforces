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
	u int
	v int
	w int
}

type query struct {
	a int
	b int
	k int
}

type testCase struct {
	n       int
	edges   []edge
	queries []query
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2057E1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", outPath, "2057E1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	totalQ := 0
	for _, tc := range tests {
		totalQ += len(tc.queries)
	}
	sb.Grow(totalQ*16 + len(tests)*64)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, len(tc.edges), len(tc.queries)))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
		}
		for _, q := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", q.a, q.b, q.k))
		}
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]int, expected)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", f, i+1)
		}
		res[i] = val
	}
	return res, nil
}

func compareAnswers(tests []testCase, expected, actual []int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch: expected %d, got %d", len(expected), len(actual))
	}
	idx := 0
	for ti, tc := range tests {
		for qi := range tc.queries {
			if expected[idx] != actual[idx] {
				return fmt.Errorf("test %d query %d mismatch: expected %d, got %d", ti+1, qi+1, expected[idx], actual[idx])
			}
			idx++
		}
	}
	return nil
}

func deterministicTests() []testCase {
	case1 := testCase{
		n: 2,
		edges: []edge{
			{u: 1, v: 2, w: 5},
		},
		queries: []query{
			{a: 1, b: 2, k: 1},
			{a: 2, b: 1, k: 1},
		},
	}

	case2 := testCase{
		n: 4,
		edges: []edge{
			{u: 1, v: 2, w: 3},
			{u: 2, v: 3, w: 2},
			{u: 3, v: 4, w: 7},
			{u: 1, v: 4, w: 5},
		},
		queries: []query{
			{a: 1, b: 3, k: 1},
			{a: 1, b: 3, k: 2},
			{a: 2, b: 4, k: 2},
		},
	}

	return []testCase{case1, case2}
}

func generateRandomCase(n, m, q int, rng *rand.Rand) testCase {
	used := make(map[int]bool)
	edges := make([]edge, 0, m)

	addEdge := func(u, v, w int) {
		if u > v {
			u, v = v, u
		}
		key := u*n + v
		if used[key] {
			return
		}
		used[key] = true
		edges = append(edges, edge{u: u, v: v, w: w})
	}

	// Build a random tree first to ensure connectivity.
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		w := rng.Intn(1_000_000_000) + 1
		addEdge(p, v, w)
	}

	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		if u == v {
			continue
		}
		w := rng.Intn(1_000_000_000) + 1
		addEdge(u, v, w)
	}

	adj := make([][]int, n)
	for _, e := range edges {
		u := e.u - 1
		v := e.v - 1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// All-pairs shortest paths by edge count to comply with k constraint.
	const inf = int(1e9)
	dist := make([][]int, n)
	for i := range dist {
		d := make([]int, n)
		for j := range d {
			d[j] = inf
		}
		dist[i] = d
	}

	qb := make([]int, n)
	for s := 0; s < n; s++ {
		head, tail := 0, 0
		qb[tail] = s
		tail++
		dist[s][s] = 0
		for head < tail {
			v := qb[head]
			head++
			for _, to := range adj[v] {
				if dist[s][to] == inf {
					dist[s][to] = dist[s][v] + 1
					qb[tail] = to
					tail++
				}
			}
		}
	}

	buildQuery := func() query {
		for {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			if a == b {
				continue
			}
			d := dist[a-1][b-1]
			if d <= 0 || d >= inf {
				continue
			}
			k := rng.Intn(d) + 1
			return query{a: a, b: b, k: k}
		}
	}

	queries := make([]query, q)
	for i := 0; i < q; i++ {
		queries[i] = buildQuery()
	}

	return testCase{
		n:       n,
		edges:   edges,
		queries: queries,
	}
}

func totalQueries(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += len(tc.queries)
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	currentN, currentM := 0, 0
	for _, tc := range tests {
		currentN += tc.n
		currentM += len(tc.edges)
	}

	const nLimit = 400
	const mLimit = 400
	const qLimit = 300_000

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	remainingN := nLimit - currentN
	remainingM := mLimit - currentM
	remainingQ := qLimit - totalQueries(tests)

	if remainingN > 0 && remainingM > 0 && remainingQ > 0 {
		// First random-heavy case.
		n1 := min(60, remainingN)
		m1 := min(120, remainingM)
		q1 := min(120_000, remainingQ)
		tests = append(tests, generateRandomCase(n1, m1, q1, rng))
		remainingN -= n1
		remainingM -= m1
		remainingQ -= q1
	}

	if remainingN > 0 && remainingM > 0 && remainingQ > 0 {
		n2 := min(70, remainingN)
		m2 := min(200, remainingM)
		q2 := min(150_000, remainingQ)
		tests = append(tests, generateRandomCase(n2, m2, q2, rng))
	}

	input := buildInput(tests)

	expectedOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	actualOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedAns, err := parseOutput(expectedOut, totalQueries(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actualOut, totalQueries(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actualOut)
		os.Exit(1)
	}

	if err := compareAnswers(tests, expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
