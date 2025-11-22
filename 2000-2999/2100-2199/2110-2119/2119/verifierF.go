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

const refSource = "2119F.go"

type testCase struct {
	name  string
	n, st int
	w     []int
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candPath := os.Args[len(os.Args)-1]

	refBin, cleanupRef, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput preview:\n%s\n", err, previewInput(input))
		os.Exit(1)
	}
	candAns, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate produced invalid output: %v\noutput:\n%s\ninput preview:\n%s\n", err, candOut, previewInput(input))
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d got %d\n", i+1, tests[i].name, refAns[i], candAns[i])
			fmt.Fprintln(os.Stderr, previewInput(buildInput([]testCase{tests[i]})))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "ref2119F-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	if filepath.Ext(abs) != ".go" {
		return abs, func() {}, nil
	}

	tmpDir, err := os.MkdirTemp("", "cand2119F-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "candidate")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = filepath.Dir(abs)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("candidate build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	ans := make([]int, expected)
	for i, tok := range tokens {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		ans[i] = v
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.st))
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.w[i]))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []testCase{}
	tests = append(tests, sampleLikeTests()...)
	tests = append(tests, lineTest(), starTest(), deepPath(), negativeSea())
	tests = append(tests, randomTests(rng, 8, 5, 25, "rand_small")...)
	tests = append(tests, randomTests(rng, 6, 100, 600, "rand_mid")...)
	tests = append(tests, randomTests(rng, 4, 5000, 12000, "rand_large")...)
	tests = append(tests, edgeHeavyTests(rng)...)
	return tests
}

func sampleLikeTests() []testCase {
	// Small structured instances to catch off-by-one and lava timing.
	return []testCase{
		buildLine(7, 4, []int{-1, -1, -1, 1, 1, 1, -1}, "sample_line"),
		{
			name: "short_cycle_like", // still a tree (star)
			n:    6,
			st:   2,
			w:    []int{1, 1, 1, 1, 1, -1},
			edges: [][2]int{
				{1, 2}, {1, 3}, {1, 4}, {4, 5}, {4, 6},
			},
		},
	}
}

func lineTest() testCase {
	return buildLine(10, 5, []int{1, 1, -1, -1, 1, -1, 1, -1, 1, -1}, "line_mix")
}

func starTest() testCase {
	n := 9
	w := []int{1, 1, -1, 1, -1, -1, 1, 1, -1}
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return testCase{name: "star_center_root", n: n, st: 3, w: w, edges: edges}
}

func deepPath() testCase {
	n := 30
	w := make([]int, n)
	for i := range w {
		if i%2 == 0 {
			w[i] = 1
		} else {
			w[i] = -1
		}
	}
	return buildLineWithWeights(w, 15, "deep_path")
}

func negativeSea() testCase {
	n := 20
	w := make([]int, n)
	for i := range w {
		w[i] = -1
	}
	w[9] = 1
	w[0] = 1
	return testCase{
		name: "mostly_negative",
		n:    n,
		st:   10,
		w:    w,
		edges: func() [][2]int {
			edges := make([][2]int, 0, n-1)
			for i := 2; i <= n; i++ {
				edges = append(edges, [2]int{i - 1, i})
			}
			return edges
		}(),
	}
}

func randomTests(rng *rand.Rand, count, minN, maxN int, tag string) []testCase {
	res := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxN-minN+1) + minN
		tc := randomTreeCase(rng, n, tag, i+1)
		res = append(res, tc)
	}
	return res
}

func edgeHeavyTests(rng *rand.Rand) []testCase {
	return []testCase{
		randomTreeCase(rng, 80000, "heavy1", 1),
		randomTreeCase(rng, 90000, "heavy2", 2),
		randomTreeCase(rng, 110000, "heavy3", 3),
	}
}

func randomTreeCase(rng *rand.Rand, n int, tag string, idx int) testCase {
	st := rng.Intn(n-1) + 2 // ensure st >= 2
	w := make([]int, n)
	for i := 0; i < n; i++ {
		if i+1 == st {
			w[i] = 1
		} else {
			if rng.Intn(2) == 0 {
				w[i] = -1
			} else {
				w[i] = 1
			}
		}
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return testCase{
		name:  fmt.Sprintf("%s_%d_n%d", tag, idx, n),
		n:     n,
		st:    st,
		w:     w,
		edges: edges,
	}
}

func buildLine(n, st int, weights []int, name string) testCase {
	if len(weights) != n {
		panic("weights length mismatch")
	}
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{i - 1, i})
	}
	return testCase{name: name, n: n, st: st, w: weights, edges: edges}
}

func buildLineWithWeights(weights []int, st int, name string) testCase {
	n := len(weights)
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{i - 1, i})
	}
	if weights[st-1] != 1 {
		weights[st-1] = 1
	}
	return testCase{name: name, n: n, st: st, w: weights, edges: edges}
}

func previewInput(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 12 {
		return strings.Join(lines[:12], "\n") + "\n..."
	}
	return input
}
