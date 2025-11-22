package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n     int
	edges [][2]int
}

type testRun struct {
	input string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, refCleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := prepareCandidate(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	seed := time.Now().UnixNano()
	tests := generateTests(seed)

	for i, tr := range tests {
		expRaw, err := runBinary(refBin, tr.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tr.input)
			os.Exit(1)
		}
		expAns, err := parseOutputs(expRaw, tr.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, expRaw)
			os.Exit(1)
		}

		actRaw, err := runBinary(candBin, tr.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i+1, err, tr.input)
			os.Exit(1)
		}
		actAns, err := parseOutputs(actRaw, tr.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tr.input, actRaw)
			os.Exit(1)
		}

		if len(expAns) != len(actAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d answers, got %d\n", i+1, len(expAns), len(actAns))
			os.Exit(1)
		}
		for idx := range expAns {
			if expAns[idx] != actAns[idx] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected %d, got %d\ninput:\n%sreference: %v\ncandidate: %v\n",
					i+1, idx+1, expAns[idx], actAns[idx], tr.input, expAns, actAns)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed (seed %d).\n", len(tests), seed)
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("unable to determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "verifier-2127H-ref-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2127H")
	cmd := exec.Command("go", "build", "-o", outPath, "2127H.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out)
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	dir := filepath.Dir(abs)
	tmpDir, err := os.MkdirTemp("", "verifier-2127H-cand-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "candidate2127H")
	cmd := exec.Command("go", "build", "-o", outPath, abs)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out)
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, cases []testCase) ([]int, error) {
	res := make([]int, len(cases))
	r := bufio.NewReader(strings.NewReader(out))
	for i := range cases {
		if _, err := fmt.Fscan(r, &res[i]); err != nil {
			return nil, fmt.Errorf("failed to read answer %d/%d: %v", i+1, len(cases), err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(r, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected: %q", extra)
	} else if err != nil && err != io.EOF {
		return nil, err
	}
	return res, nil
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func sampleCases() []testCase {
	return []testCase{
		{
			n: 4,
			edges: [][2]int{
				{1, 2}, {1, 3}, {2, 3}, {3, 4},
			},
		},
		{
			n: 7,
			edges: [][2]int{
				{1, 2}, {1, 3}, {1, 4}, {2, 4}, {3, 4},
				{4, 5}, {4, 6}, {5, 6}, {5, 7}, {6, 7},
			},
		},
		{
			n: 9,
			edges: [][2]int{
				{1, 2}, {1, 3}, {3, 4}, {3, 7}, {4, 5},
				{4, 6}, {5, 6}, {7, 8}, {7, 9}, {8, 9},
			},
		},
	}
}

func deterministicRuns() []testRun {
	var runs []testRun
	runs = append(runs, makeRun(sampleCases()))

	runs = append(runs, makeRun([]testCase{
		lineGraph(3),
		lineGraph(10),
		cycleGraph(6),
	}))

	runs = append(runs, makeRun([]testCase{
		starWithChord(8),
		starWithChord(12),
	}))

	return runs
}

func lineGraph(n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return testCase{n: n, edges: edges}
}

func cycleGraph(n int) testCase {
	edges := make([][2]int, 0, n)
	for i := 1; i <= n; i++ {
		edges = append(edges, [2]int{i, (i % n) + 1})
	}
	return testCase{n: n, edges: edges}
}

func starWithChord(n int) testCase {
	edges := make([][2]int, 0, n-1+1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	if n >= 4 {
		edges = append(edges, [2]int{2, 3})
	}
	return testCase{n: n, edges: edges}
}

func makeRun(cases []testCase) testRun {
	return testRun{input: buildInput(cases), cases: cases}
}

func randomRuns(rng *rand.Rand) []testRun {
	var runs []testRun

	for i := 0; i < 35; i++ {
		runs = append(runs, makeRun(randomBatch(rng, 5, 10, 6)))
	}

	for i := 0; i < 8; i++ {
		runs = append(runs, makeRun(randomBatch(rng, 10, 30, 10)))
	}

	for i := 0; i < 3; i++ {
		runs = append(runs, makeRun(randomBatch(rng, 1, 30, 12)))
	}

	return runs
}

func randomBatch(rng *rand.Rand, maxCases, maxN, maxExtra int) []testCase {
	var cases []testCase
	budget := 900
	numCases := rng.Intn(maxCases) + 1
	for i := 0; i < numCases; i++ {
		n := rng.Intn(maxN-2) + 3
		if n*n > budget && len(cases) > 0 {
			break
		}
		budget -= n * n
		extra := rng.Intn(maxExtra + 1)
		tc := randomCase(rng, n, extra)
		cases = append(cases, tc)
	}
	if len(cases) == 0 {
		cases = append(cases, randomCase(rng, 6, 3))
	}
	return cases
}

func randomCase(rng *rand.Rand, n, extra int) testCase {
	edges := buildGraph(rng, n, extra)
	return testCase{n: n, edges: edges}
}

func buildGraph(rng *rand.Rand, n, extra int) [][2]int {
	edges := make([][2]int, 0, n-1+extra)
	// Random tree for connectivity.
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}

	additional := 0
	targetExtra := extra
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	for attempts := 0; attempts < 200 && additional < targetExtra; attempts++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if edgeExists(adj, u, v) {
			continue
		}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		if ok := cyclesOK(adj); ok {
			edges = append(edges, [2]int{u + 1, v + 1})
			additional++
		} else {
			adj[u] = adj[u][:len(adj[u])-1]
			adj[v] = adj[v][:len(adj[v])-1]
		}
	}

	return edges
}

func edgeExists(adj [][]int, u, v int) bool {
	for _, w := range adj[u] {
		if w == v {
			return true
		}
	}
	return false
}

func cyclesOK(adj [][]int) bool {
	n := len(adj)
	counts := make([]int, n)
	visited := make([]bool, n)
	path := make([]int, 0, n)

	var dfs func(int, int, int)
	dfs = func(u, parent, start int) {
		visited[u] = true
		path = append(path, u)
		for _, v := range adj[u] {
			if v < start {
				continue
			}
			if v == start && len(path) >= 3 {
				for _, x := range path {
					counts[x]++
					if counts[x] > 5 {
						return
					}
				}
			} else if !visited[v] && v > start {
				dfs(v, u, start)
				if countsExceed(counts) {
					return
				}
			}
		}
		path = path[:len(path)-1]
		visited[u] = false
	}

	for s := 0; s < n; s++ {
		dfs(s, -1, s)
		if countsExceed(counts) {
			return false
		}
	}
	return true
}

func countsExceed(cnt []int) bool {
	for _, v := range cnt {
		if v > 5 {
			return true
		}
	}
	return false
}

func generateTests(seed int64) []testRun {
	rng := rand.New(rand.NewSource(seed))
	tests := deterministicRuns()
	tests = append(tests, randomRuns(rng)...)
	return tests
}
