package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "0-999/800-899/880-889/883/883B.go"

type testInput struct {
	n, m, k int
	ranks   []int
	edges   [][2]int
}

type testCase struct {
	input string
	data  testInput
}

type candidateResult struct {
	impossible bool
	ranks      []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		refRes, err := parseOutput(refOut, tc.data.n)
		if err != nil {
			fail("failed to parse reference output on test %d: %v\noutput:\n%s", i+1, err, refOut)
		}
		candRes, err := parseOutput(candOut, tc.data.n)
		if err != nil {
			fail("failed to parse candidate output on test %d: %v\noutput:\n%s", i+1, err, candOut)
		}
		if refRes.impossible {
			if !candRes.impossible {
				fail("test %d: expected -1 but candidate produced assignment\ninput:\n%s\ncandidate output:\n%s", i+1, tc.input, candOut)
			}
			continue
		}
		if candRes.impossible {
			fail("test %d: candidate reported -1 but reference found a solution\ninput:\n%s", i+1, tc.input)
		}
		if err := validateAssignment(&tc.data, candRes.ranks); err != nil {
			fail("test %d: invalid assignment: %v\ninput:\n%s\ncandidate output:\n%s", i+1, err, tc.input, candOut)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "883B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, n int) (candidateResult, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	sc.Buffer(make([]byte, 1024), 4<<20)
	var tokens []string
	for sc.Scan() {
		tokens = append(tokens, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return candidateResult{}, err
	}
	if len(tokens) == 0 {
		return candidateResult{}, errors.New("empty output")
	}
	if len(tokens) == 1 && tokens[0] == "-1" {
		return candidateResult{impossible: true}, nil
	}
	if tokens[0] == "-1" && len(tokens) > 1 {
		return candidateResult{}, errors.New("extra data after -1")
	}
	if len(tokens) != n {
		return candidateResult{}, fmt.Errorf("expected %d ranks, got %d", n, len(tokens))
	}
	res := candidateResult{ranks: make([]int, n)}
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return candidateResult{}, fmt.Errorf("invalid integer %q", tok)
		}
		res.ranks[i] = val
	}
	return res, nil
}

func validateAssignment(ti *testInput, ranks []int) error {
	if len(ranks) != ti.n {
		return fmt.Errorf("expected %d ranks, got %d", ti.n, len(ranks))
	}
	used := make([]bool, ti.k+1)
	for i := 0; i < ti.n; i++ {
		v := ranks[i]
		if v < 1 || v > ti.k {
			return fmt.Errorf("rank %d out of range for soldier %d", v, i+1)
		}
		if ti.ranks[i] > 0 && ti.ranks[i] != v {
			return fmt.Errorf("soldier %d rank mismatch: expected %d got %d", i+1, ti.ranks[i], v)
		}
		used[v] = true
	}
	for _, e := range ti.edges {
		if ranks[e[0]] <= ranks[e[1]] {
			return fmt.Errorf("order constraint violated for (%d,%d)", e[0]+1, e[1]+1)
		}
	}
	for r := 1; r <= ti.k; r++ {
		if !used[r] {
			return fmt.Errorf("rank %d unused", r)
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	// manual simple cases
	tests = append(tests, buildTestCase(testInput{
		n: 1, k: 1,
		ranks: []int{0},
	}))
	tests = append(tests, buildTestCase(testInput{
		n: 2, k: 2,
		ranks: []int{1, 2},
		edges: [][2]int{{0, 1}},
	}))
	tests = append(tests, buildTestCase(testInput{
		n: 3, k: 2,
		ranks: []int{0, 0, 0},
		edges: [][2]int{{0, 1}, {0, 2}},
	}))
	tests = append(tests, buildTestCase(testInput{
		n: 3, k: 3,
		ranks: []int{1, 2, 3},
		edges: [][2]int{{0, 1}, {1, 2}, {2, 0}},
	}))

	// random solvable tests of varying sizes
	for i := 0; i < 25; i++ {
		tests = append(tests, randomSolvableTest(rng, 6, 4))
	}
	for i := 0; i < 25; i++ {
		tests = append(tests, randomSolvableTest(rng, 40, 10))
	}
	for i := 0; i < 15; i++ {
		tests = append(tests, randomSolvableTest(rng, 150, 20))
	}
	tests = append(tests, randomSolvableTest(rng, 500, 40))
	tests = append(tests, randomSolvableTest(rng, 2000, 80))

	// random impossible cases
	for i := 0; i < 10; i++ {
		tests = append(tests, randomImpossibleCycleTest(rng, 30, 10))
	}
	tests = append(tests, randomImpossibleCycleTest(rng, 200, 40))
	tests = append(tests, impossibleDueToLargeK(50, 80))

	return tests
}

func randomSolvableTest(rng *rand.Rand, maxN, maxK int) testCase {
	if maxN < 1 {
		maxN = 1
	}
	n := rng.Intn(maxN) + 1
	k := rng.Intn(maxInt(1, minInt(maxK, n))) + 1
	if k > n {
		k = n
	}
	base := make([]int, n)
	for i := 0; i < k && i < n; i++ {
		base[i] = i + 1
	}
	for i := k; i < n; i++ {
		base[i] = rng.Intn(k) + 1
	}
	rng.Shuffle(n, func(i, j int) { base[i], base[j] = base[j], base[i] })
	ranks := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Float64() < 0.4 {
			ranks[i] = base[i]
		}
	}
	maxEdges := minInt(n*(n-1), 3*n+10)
	target := rng.Intn(maxEdges + 1)
	edges := make([][2]int, 0, target)
	seen := make(map[int]struct{})
	limit := target*5 + 100
	attempts := 0
	for len(edges) < target && attempts < limit {
		x := rng.Intn(n)
		y := rng.Intn(n)
		if x == y || base[x] <= base[y] {
			attempts++
			continue
		}
		key := x*n + y
		if _, ok := seen[key]; ok {
			attempts++
			continue
		}
		seen[key] = struct{}{}
		edges = append(edges, [2]int{x, y})
		attempts++
	}
	ti := testInput{
		n:     n,
		m:     len(edges),
		k:     k,
		ranks: ranks,
		edges: edges,
	}
	return buildTestCase(ti)
}

func randomImpossibleCycleTest(rng *rand.Rand, maxN, maxK int) testCase {
	tc := randomSolvableTest(rng, maxN, maxK)
	ti := cloneTestInput(tc.data)
	if ti.n < 2 {
		return buildTestCase(testInput{
			n:     2,
			k:     2,
			ranks: []int{1, 2},
			edges: [][2]int{{0, 1}, {1, 0}},
		})
	}
	addEdgeUnique(&ti, 0, 1)
	addEdgeUnique(&ti, 1, 0)
	ti.m = len(ti.edges)
	return buildTestCase(ti)
}

func impossibleDueToLargeK(n, k int) testCase {
	if k <= n {
		k = n + 1
	}
	if k > 200000 {
		k = 200000
	}
	ranks := make([]int, n)
	ti := testInput{
		n:     n,
		k:     k,
		ranks: ranks,
	}
	return buildTestCase(ti)
}

func addEdgeUnique(ti *testInput, x, y int) {
	for _, e := range ti.edges {
		if e[0] == x && e[1] == y {
			return
		}
	}
	ti.edges = append(ti.edges, [2]int{x, y})
}

func cloneTestInput(src testInput) testInput {
	dst := testInput{
		n: src.n, m: src.m, k: src.k,
		ranks: append([]int(nil), src.ranks...),
		edges: append([][2]int(nil), src.edges...),
	}
	return dst
}

func buildTestCase(data testInput) testCase {
	data.m = len(data.edges)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", data.n, data.m, data.k)
	for i, v := range data.ranks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range data.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	return testCase{
		input: sb.String(),
		data: testInput{
			n:     data.n,
			m:     data.m,
			k:     data.k,
			ranks: append([]int(nil), data.ranks...),
			edges: append([][2]int(nil), data.edges...),
		},
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ensure dependency on math/rand/time via dummy usage when limits small
var _ = math.MaxInt
