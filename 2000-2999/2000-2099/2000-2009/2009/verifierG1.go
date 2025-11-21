package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2000-2009/2009/2009G1.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2009G1.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n, k, q int
	a       []int
	queries [][2]int
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 7, k: 5, q: 3,
			a:       []int{1, 2, 3, 2, 1, 2, 3},
			queries: [][2]int{{1, 5}, {2, 6}, {3, 7}},
		},
		{
			n: 8, k: 4, q: 2,
			a:       []int{4, 3, 1, 1, 2, 4, 3, 2},
			queries: [][2]int{{3, 6}, {2, 5}},
		},
		{
			n: 5, k: 4, q: 2,
			a:       []int{4, 5, 1, 2, 3},
			queries: [][2]int{{1, 4}, {2, 5}},
		},
		{
			n: 4, k: 2, q: 3,
			a:       []int{1, 2, 3, 4},
			queries: [][2]int{{1, 2}, {2, 3}, {3, 4}},
		},
	}
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	return arr
}

func randomQueries(rng *rand.Rand, n, k, q int) [][2]int {
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n-k+1) + 1
		queries[i] = [2]int{l, l + k - 1}
	}
	return queries
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	totalQ := 0
	for len(tests) < 60 && totalN < 200000 && totalQ < 200000 {
		n := rng.Intn(200) + 1
		k := rng.Intn(n) + 1
		q := rng.Intn(200) + 1
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		if totalQ+q > 200000 {
			q = 200000 - totalQ
		}
		if n <= 0 || q <= 0 || k > n {
			break
		}
		tests = append(tests, testCase{
			n:       n,
			k:       k,
			q:       q,
			a:       randomArray(rng, n),
			queries: randomQueries(rng, n, k, q),
		})
		totalN += n
		totalQ += q
	}
	return tests
}

func formatTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.q))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2009G1-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2009G1")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseInts(output string, expected int) ([]int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	values := make([]int64, 0, expected)
	for scanner.Scan() {
		var val int64
		token := scanner.Text()
		if _, err := fmt.Sscan(token, &val); err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", token, err)
		}
		values = append(values, val)
	}
	if len(values) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(values))
	}
	return values, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatTests(tests)
	totalAnswers := 0
	for _, tc := range tests {
		totalAnswers += tc.q
	}

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected, err := parseInts(refOut, totalAnswers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseInts(userOut, totalAnswers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	idx := 0
	for ti, tc := range tests {
		for qi := 0; qi < tc.q; qi++ {
			if got[idx] != expected[idx] {
				fmt.Fprintf(os.Stderr, "test %d query %d mismatch: expected %d got %d\ninputs: n=%d k=%d q=%d\narray=%v query=%v\n", ti+1, qi+1, expected[idx], got[idx], tc.n, tc.k, tc.q, tc.a, tc.queries[qi])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
