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

type treeSpec struct {
	n, k  int
	edges [][2]int
}

type testCase struct {
	caseCount int
	input     string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.caseCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", i+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		gotVals, err := parseOutputs(gotOut, tc.caseCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", i+1, err, gotOut)
			os.Exit(1)
		}

		for idx := 0; idx < tc.caseCount; idx++ {
			if refVals[idx] != gotVals[idx] {
				fmt.Fprintf(os.Stderr, "test %d case %d failed: expected %d got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					i+1, idx+1, refVals[idx], gotVals[idx], tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2167F_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2167F.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 40, 40)...)
	tests = append(tests, randomTests(rng, 40, 400)...)
	tests = append(tests, randomTests(rng, 25, 4000)...)
	tests = append(tests, extremeTests(rng)...)
	return tests
}

func manualTests() []testCase {
	return []testCase{
		makeTestCase([]treeSpec{
			{
				n: 2,
				k: 2,
				edges: [][2]int{
					{1, 2},
				},
			},
		}),
		makeTestCase([]treeSpec{
			{
				n: 3,
				k: 2,
				edges: [][2]int{
					{1, 2},
					{2, 3},
				},
			},
			{
				n: 6,
				k: 3,
				edges: [][2]int{
					{1, 2},
					{1, 3},
					{2, 4},
					{2, 5},
					{3, 6},
				},
			},
		}),
		makeTestCase([]treeSpec{
			{
				n: 7,
				k: 4,
				edges: [][2]int{
					{1, 2},
					{1, 3},
					{2, 4},
					{3, 5},
					{3, 6},
					{5, 7},
				},
			},
		}),
	}
}

func randomTests(rng *rand.Rand, batches int, maxCaseN int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		targetCases := rng.Intn(4) + 1
		cases := make([]treeSpec, 0, targetCases)
		sumN := 0
		for len(cases) < targetCases {
			remaining := 200000 - sumN
			if remaining < 2 {
				break
			}
			limit := maxCaseN
			if limit > remaining {
				limit = remaining
			}
			if limit < 2 {
				break
			}
			n := rng.Intn(limit-1) + 2
			cases = append(cases, randomTreeSpec(rng, n))
			sumN += n
		}
		if len(cases) == 0 {
			cases = append(cases, randomTreeSpec(rng, 2))
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func extremeTests(rng *rand.Rand) []testCase {
	var tests []testCase
	tests = append(tests, makeTestCase([]treeSpec{
		pathSpec(200000, 2),
	}))
	tests = append(tests, makeTestCase([]treeSpec{
		starSpec(200000, 200000),
	}))
	tests = append(tests, makeTestCase([]treeSpec{
		randomTreeSpecWithEdges(200000, max(2, 200000/3), randomTreeEdges(200000, rng)),
	}))
	return tests
}

func randomTreeSpec(rng *rand.Rand, n int) treeSpec {
	return randomTreeSpecWithEdges(n, chooseK(rng, n), randomTreeEdges(n, rng))
}

func randomTreeSpecWithEdges(n, k int, edges [][2]int) treeSpec {
	if k < 2 {
		k = 2
	}
	if k > n {
		k = n
	}
	return treeSpec{n: n, k: k, edges: edges}
}

func chooseK(rng *rand.Rand, n int) int {
	switch rng.Intn(4) {
	case 0:
		return 2
	case 1:
		return n
	case 2:
		return max(2, n/2)
	default:
		return rng.Intn(n-1) + 2
	}
}

func randomTreeEdges(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		parent := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{parent, v})
	}
	rng.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
	return edges
}

func pathSpec(n, k int) treeSpec {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return treeSpec{n: n, k: k, edges: edges}
}

func starSpec(n, k int) treeSpec {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return treeSpec{n: n, k: k, edges: edges}
}

func makeTestCase(cases []treeSpec) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(strconv.Itoa(c.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(c.k))
		sb.WriteByte('\n')
		if len(c.edges) != c.n-1 {
			panic("edge count mismatch")
		}
		for _, e := range c.edges {
			sb.WriteString(strconv.Itoa(e[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(e[1]))
			sb.WriteByte('\n')
		}
	}
	return testCase{
		caseCount: len(cases),
		input:     sb.String(),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
