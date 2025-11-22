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

type edge struct {
	u, v int
	x, y int64
}

type testCase struct {
	n     int
	edges []edge
}

type testRun struct {
	input string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		expPerms, err := parseOutputs(expRaw, tr.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, expRaw)
			os.Exit(1)
		}

		actRaw, err := runBinary(candBin, tr.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i+1, err, tr.input)
			os.Exit(1)
		}
		actPerms, err := parseOutputs(actRaw, tr.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tr.input, actRaw)
			os.Exit(1)
		}

		if len(expPerms) != len(actPerms) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d test cases, got %d\n", i+1, len(expPerms), len(actPerms))
			os.Exit(1)
		}
		for tcIdx, tc := range tr.cases {
			if len(expPerms[tcIdx]) != tc.n || len(actPerms[tcIdx]) != tc.n {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: permutation length mismatch\n", i+1, tcIdx+1)
				os.Exit(1)
			}
			expVal := valueOf(tc, expPerms[tcIdx])
			actVal := valueOf(tc, actPerms[tcIdx])
			if actVal != expVal {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected value %d, got %d\ninput:\n%sreference perm: %v\ncandidate perm: %v\n",
					i+1, tcIdx+1, expVal, actVal, tr.input, expPerms[tcIdx], actPerms[tcIdx])
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
	tmpDir, err := os.MkdirTemp("", "verifier-2143C-ref-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref2143C")
	cmd := exec.Command("go", "build", "-o", outPath, "2143C.go")
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
	tmpDir, err := os.MkdirTemp("", "verifier-2143C-cand-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "candidate2143C")
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

func parseOutputs(out string, cases []testCase) ([][]int, error) {
	res := make([][]int, len(cases))
	r := bufio.NewReader(strings.NewReader(out))
	for idx, tc := range cases {
		perm := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(r, &perm[i]); err != nil {
				return nil, fmt.Errorf("test case %d: failed to read permutation value %d/%d: %v", idx+1, i+1, tc.n, err)
			}
		}
		if err := validatePermutation(perm); err != nil {
			return nil, fmt.Errorf("test case %d: %v", idx+1, err)
		}
		res[idx] = perm
	}
	var extra string
	if _, err := fmt.Fscan(r, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected: %q", extra)
	} else if err != nil && err != io.EOF {
		return nil, err
	}
	return res, nil
}

func validatePermutation(p []int) error {
	n := len(p)
	seen := make([]bool, n+1)
	for i, v := range p {
		if v < 1 || v > n {
			return fmt.Errorf("value %d at position %d is out of range 1..%d", v, i+1, n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
	}
	return nil
}

func valueOf(tc testCase, perm []int) int64 {
	val := int64(0)
	for _, e := range tc.edges {
		if perm[e.u] > perm[e.v] {
			val += e.x
		} else {
			val += e.y
		}
	}
	return val
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			u := e.u + 1
			v := e.v + 1
			if u > v {
				u, v = v, u
			}
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u, v, e.x, e.y))
		}
	}
	return sb.String()
}

func sampleCases() []testCase {
	return []testCase{
		{
			n: 3,
			edges: []edge{
				{u: 0, v: 1, x: 2, y: 1},
				{u: 1, v: 2, x: 3, y: 2},
			},
		},
		{
			n: 5,
			edges: []edge{
				{u: 0, v: 1, x: 1, y: 3},
				{u: 0, v: 4, x: 2, y: 1},
				{u: 1, v: 3, x: 5, y: 7},
				{u: 1, v: 2, x: 1, y: 100},
			},
		},
		{
			n: 5,
			edges: []edge{
				{u: 1, v: 4, x: 5, y: 2},
				{u: 2, v: 4, x: 4, y: 6},
				{u: 3, v: 4, x: 1, y: 5},
				{u: 0, v: 4, x: 3, y: 5},
			},
		},
	}
}

func deterministicRuns() []testRun {
	var runs []testRun
	runs = append(runs, makeRun(sampleCases()))

	runs = append(runs, makeRun([]testCase{
		// n=2 edge equal weights
		{
			n: 2,
			edges: []edge{
				{u: 0, v: 1, x: 5, y: 5},
			},
		},
		// chain with mixed weights
		{
			n: 4,
			edges: []edge{
				{u: 0, v: 1, x: 10, y: 1},
				{u: 1, v: 2, x: 2, y: 20},
				{u: 2, v: 3, x: 3, y: 3},
			},
		},
	}))

	// star with different edge preferences
	runs = append(runs, makeRun([]testCase{
		buildStar(6, []int64{1, 10, 5, 7, 3}, []int64{9, 2, 6, 4, 8}),
	}))

	return runs
}

func buildStar(n int, xs, ys []int64) testCase {
	edges := make([]edge, n-1)
	for i := 1; i < n; i++ {
		x := xs[(i-1)%len(xs)]
		y := ys[(i-1)%len(ys)]
		edges[i-1] = edge{u: 0, v: i, x: x, y: y}
	}
	return testCase{n: n, edges: edges}
}

func randomRuns(rng *rand.Rand) []testRun {
	var runs []testRun

	for i := 0; i < 30; i++ {
		cases := make([]testCase, rng.Intn(3)+1)
		for j := range cases {
			n := rng.Intn(8) + 2
			cases[j] = randomCase(rng, n)
		}
		runs = append(runs, makeRun(cases))
	}

	for i := 0; i < 6; i++ {
		cases := make([]testCase, rng.Intn(2)+1)
		for j := range cases {
			n := rng.Intn(200) + 50
			cases[j] = randomCase(rng, n)
		}
		runs = append(runs, makeRun(cases))
	}

	for i := 0; i < 2; i++ {
		// one bigger case per run
		n := rng.Intn(4000) + 1000
		runs = append(runs, makeRun([]testCase{randomCase(rng, n)}))
	}

	return runs
}

func randomCase(rng *rand.Rand, n int) testCase {
	edges := make([]edge, 0, n-1)
	// generate random tree via Prufer
	if n == 2 {
		x := rng.Int63n(1_000_000_000) + 1
		y := rng.Int63n(1_000_000_000) + 1
		edges = append(edges, edge{u: 0, v: 1, x: x, y: y})
		return testCase{n: n, edges: edges}
	}
	prufer := make([]int, n-2)
	for i := range prufer {
		prufer[i] = rng.Intn(n)
	}
	deg := make([]int, n)
	for i := 0; i < n; i++ {
		deg[i] = 1
	}
	for _, v := range prufer {
		deg[v]++
	}
	for _, v := range prufer {
		var u int
		for i := 0; i < n; i++ {
			if deg[i] == 1 {
				u = i
				break
			}
		}
		deg[u]--
		deg[v]--
		addEdge(&edges, u, v, rng)
	}
	// last two leaves
	var u, v int
	u = -1
	for i := 0; i < n; i++ {
		if deg[i] == 1 {
			if u == -1 {
				u = i
			} else {
				v = i
				break
			}
		}
	}
	addEdge(&edges, u, v, rng)

	return testCase{n: n, edges: edges}
}

func addEdge(edges *[]edge, u, v int, rng *rand.Rand) {
	x := rng.Int63n(1_000_000_000) + 1
	y := rng.Int63n(1_000_000_000) + 1
	if u > v {
		u, v = v, u
	}
	*edges = append(*edges, edge{u: u, v: v, x: x, y: y})
}

func makeRun(cases []testCase) testRun {
	return testRun{input: buildInput(cases), cases: cases}
}

func generateTests(seed int64) []testRun {
	rng := rand.New(rand.NewSource(seed))
	tests := deterministicRuns()
	tests = append(tests, randomRuns(rng)...)
	return tests
}
