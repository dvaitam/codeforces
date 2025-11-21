package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "1000-1999/1700-1799/1770-1779/1773/1773J.go"

type edge struct {
	u int
	v int
	x int64
}

type graphData struct {
	n     int
	m     int
	mod   int64
	edges []edge
}

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		data, err := parseGraph(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%sreference output:\n%s\n", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		possible, err := isPossible(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to interpret reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if err := verifyCandidate(candOut, data, possible); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\ncandidate output:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "verifier1773J-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), cleanPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

func isPossible(output string) (bool, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return false, fmt.Errorf("empty output")
	}
	if fields[0] == "-1" {
		if len(fields) > 1 {
			return false, fmt.Errorf("unexpected extra tokens after -1")
		}
		return false, nil
	}
	return true, nil
}

func parseGraph(input string) (graphData, error) {
	reader := strings.NewReader(input)
	var n, m int
	var mod int64
	if _, err := fmt.Fscan(reader, &n, &m, &mod); err != nil {
		return graphData{}, err
	}
	edges := make([]edge, m+1)
	for i := 1; i <= m; i++ {
		var u, v int
		var x int64
		if _, err := fmt.Fscan(reader, &u, &v, &x); err != nil {
			return graphData{}, err
		}
		edges[i] = edge{u: u, v: v, x: ((x % mod) + mod) % mod}
	}
	return graphData{n: n, m: m, mod: mod, edges: edges}, nil
}

func verifyCandidate(output string, data graphData, expectPossible bool) error {
	output = strings.TrimSpace(output)
	if !expectPossible {
		if output != "-1" {
			return fmt.Errorf("expected -1 but got %q", output)
		}
		return nil
	}

	reader := strings.NewReader(output)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return fmt.Errorf("failed to read number of operations: %v", err)
	}
	if t < 0 || t > 2*data.m {
		return fmt.Errorf("operation count %d out of allowed range [0, %d]", t, 2*data.m)
	}

	counters := make([]int64, data.m+1)
	treeSize := data.n - 1
	for op := 0; op < t; op++ {
		var val int64
		if _, err := fmt.Fscan(reader, &val); err != nil {
			return fmt.Errorf("failed to read value for operation %d: %v", op+1, err)
		}
		if val < 0 || val >= data.mod {
			return fmt.Errorf("operation %d uses invalid increment %d", op+1, val)
		}

		edgeIDs := make([]int, treeSize)
		for i := 0; i < treeSize; i++ {
			if _, err := fmt.Fscan(reader, &edgeIDs[i]); err != nil {
				return fmt.Errorf("failed to read edge #%d for operation %d: %v", i+1, op+1, err)
			}
			if edgeIDs[i] < 1 || edgeIDs[i] > data.m {
				return fmt.Errorf("operation %d references invalid edge id %d", op+1, edgeIDs[i])
			}
		}

		if err := ensureSpanningTree(edgeIDs, data); err != nil {
			return fmt.Errorf("operation %d: %v", op+1, err)
		}

		if data.n == 1 {
			continue
		}
		for _, eid := range edgeIDs {
			counters[eid] = (counters[eid] + val) % data.mod
		}
	}

	var leftover string
	if _, err := fmt.Fscan(reader, &leftover); err == nil {
		return fmt.Errorf("unexpected extra output after listed operations")
	}

	for i := 1; i <= data.m; i++ {
		if counters[i]%data.mod != data.edges[i].x {
			return fmt.Errorf("edge %d reached %d, expected %d", i, counters[i]%data.mod, data.edges[i].x)
		}
	}
	return nil
}

func ensureSpanningTree(edgeIDs []int, data graphData) error {
	if data.n == 1 {
		if len(edgeIDs) != 0 {
			return fmt.Errorf("expected 0 edges for n=1 but got %d", len(edgeIDs))
		}
		return nil
	}
	if len(edgeIDs) != data.n-1 {
		return fmt.Errorf("expected %d edges but got %d", data.n-1, len(edgeIDs))
	}
	used := make([]bool, data.m+1)
	dsu := newDSU(data.n)
	for _, eid := range edgeIDs {
		if used[eid] {
			return fmt.Errorf("edge %d used multiple times in the same tree", eid)
		}
		used[eid] = true
		u := data.edges[eid].u
		v := data.edges[eid].v
		if !dsu.union(u, v) {
			return fmt.Errorf("edges form a cycle, not a tree")
		}
	}
	root := dsu.find(1)
	for i := 2; i <= data.n; i++ {
		if dsu.find(i) != root {
			return fmt.Errorf("edges do not connect all vertices")
		}
	}
	return nil
}

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &dsu{parent: parent, size: size}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "sample1", input: "3 3 101\n1 2 30\n2 3 40\n3 1 50\n"},
		{name: "sample2", input: "2 2 37\n1 2 8\n1 2 15\n"},
		{name: "sample3_impossible", input: "5 4 5\n1 3 1\n2 3 2\n2 5 3\n4 1 4\n"},
	}

	tests = append(tests, deterministicCases()...)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random_small_%d", i+1), 3, 8, 10))
	}
	for i := 0; i < 6; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random_mid_%d", i+1), 30, 80, 150))
	}
	tests = append(tests,
		randomCase(rng, "random_dense", 80, 200, 400),
		randomCase(rng, "random_max", 500, 1000, 1000),
		randomCase(rng, "random_sparse_large", 500, 600, 700),
	)
	return tests
}

func deterministicCases() []testCase {
	return []testCase{
		buildCase("multi_edge_small", 2, 2, 11, []edge{
			{u: 1, v: 2, x: 3},
			{u: 1, v: 2, x: 7},
		}),
		buildCase("triangle_prime2", 3, 3, 2, []edge{
			{1, 2, 1},
			{2, 3, 0},
			{1, 3, 1},
		}),
		buildCase("line_three", 3, 2, 17, []edge{
			{1, 2, 5},
			{2, 3, 9},
		}),
	}
}

func buildCase(name string, n, m int, mod int64, edges []edge) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, m, mod)
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d %d\n", e.u, e.v, e.x%mod)
	}
	return testCase{name: name, input: b.String()}
}

func randomCase(rng *rand.Rand, name string, maxN, maxM, extraEdges int) testCase {
	if maxN < 2 {
		maxN = 2
	}
	n := rng.Intn(maxN-1) + 2
	mLimit := maxM
	if mLimit < n-1 {
		mLimit = n - 1
	}
	if mLimit > 1000 {
		mLimit = 1000
	}
	m := n - 1
	if extraEdges > 0 {
		additional := rng.Intn(extraEdges)
		if m+additional > mLimit {
			additional = mLimit - m
		}
		m += additional
	}

	edges := make([][2]int, 0, m)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		edges = append(edges, [2]int{u, v})
	}

	primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 101, 103, 1_000_003, 1_000_033, 1_000_237, 1_000_000_007}
	mod := primes[rng.Intn(len(primes))]

	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, len(edges), mod)
	for _, e := range edges {
		x := rng.Int63n(mod)
		fmt.Fprintf(&b, "%d %d %d\n", e[0], e[1], x)
	}
	return testCase{name: name, input: b.String()}
}
