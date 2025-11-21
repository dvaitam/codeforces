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
)

const mod = 1000000007

type edge struct {
	u, v int
	seq  []int
}

type testCase struct {
	n     int
	edges []edge
	desc  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		input := buildInput(tc)
		expStdout, _, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.desc, err, input)
			os.Exit(1)
		}
		expVals, err := parseOutput(expStdout, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.desc, err, expStdout)
			os.Exit(1)
		}

		gotStdout, gotStderr, err := runBinary(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d (%s): %v\nstderr:\n%s\n", idx+1, tc.desc, err, gotStderr)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotStdout, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.desc, err, input, gotStdout)
			os.Exit(1)
		}
		if err := compareOutputs(expVals, gotVals); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.desc, err, input, formatLines(expVals), formatLines(gotVals))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-331E2-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE2")
	cmd := exec.Command("go", "build", "-o", outPath, "331E2.go")
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(output))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(path, input string) (string, string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutput(out string, n int) ([]int, error) {
	expected := 2 * n
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	res := make([]int, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: not an integer (%v)", i+1, err)
		}
		norm := int((val%mod + mod) % mod)
		res[i] = norm
	}
	return res, nil
}

func compareOutputs(expected, actual []int) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d numbers, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("line %d: expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d", e.u, e.v, len(e.seq))
		for _, val := range e.seq {
			fmt.Fprintf(&sb, " %d", val)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func formatLines(vals []int) string {
	var sb strings.Builder
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		n: 6,
		edges: []edge{
			{1, 2, []int{1, 2}},
			{2, 3, []int{3}},
			{3, 4, []int{4, 5}},
			{4, 5, []int{}},
			{5, 3, []int{3}},
			{6, 1, []int{6}},
		},
		desc: "sample",
	})
	tests = append(tests, testCase{
		n:    1,
		desc: "single-node",
	})
	tests = append(tests, testCase{
		n: 2,
		edges: []edge{
			{1, 2, []int{1}},
		},
		desc: "single-edge",
	})
	tests = append(tests, testCase{
		n: 3,
		edges: []edge{
			{1, 2, []int{1, 2}},
			{2, 3, []int{2, 3}},
		},
		desc: "chain-with-prefix",
	})
	tests = append(tests, testCase{
		n: 4,
		edges: []edge{
			{1, 2, []int{}},
			{2, 3, []int{1, 2, 3}},
			{3, 4, []int{2, 3, 4}},
		},
		desc: "empty-vision-edge",
	})
	tests = append(tests, denseCase(6))
	tests = append(tests, cycleCase())

	rng := rand.New(rand.NewSource(1331331))
	for i := 0; i < 60; i++ {
		n := rng.Intn(8) + 2
		tc := randomCase(rng, n)
		tc.desc = fmt.Sprintf("random-%d", i+1)
		tests = append(tests, tc)
	}
	return tests
}

func denseCase(n int) testCase {
	var edges []edge
	flag := false
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			u, v := i, j
			if flag {
				u, v = v, u
			}
			flag = !flag
			edges = append(edges, edge{u: u, v: v, seq: []int{u, v}})
		}
	}
	return testCase{
		n:     n,
		edges: edges,
		desc:  "dense",
	}
}

func cycleCase() testCase {
	return testCase{
		n: 5,
		edges: []edge{
			{1, 2, []int{1}},
			{2, 3, []int{}},
			{3, 4, []int{3, 4}},
			{4, 5, []int{5}},
			{5, 1, []int{1}},
		},
		desc: "cycle",
	}
}

func randomCase(rng *rand.Rand, n int) testCase {
	type pair struct{ a, b int }
	var pairs []pair
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			pairs = append(pairs, pair{i, j})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})
	m := 0
	if len(pairs) > 0 {
		m = rng.Intn(len(pairs) + 1)
	}
	edges := make([]edge, 0, m)
	for i := 0; i < m; i++ {
		x, y := pairs[i].a, pairs[i].b
		if rng.Intn(2) == 0 {
			x, y = y, x
		}
		k := rng.Intn(n + 1)
		seq := make([]int, k)
		for j := 0; j < k; j++ {
			seq[j] = rng.Intn(n) + 1
		}
		edges = append(edges, edge{u: x, v: y, seq: seq})
	}
	return testCase{n: n, edges: edges}
}
