package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2040-2049/2044/2044H.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2044H.go")
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

type query struct {
	x1, y1, x2, y2 int
}

type testCase struct {
	n, q int
	mat  [][]int64
	qs   []query
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 4, q: 3,
			mat: [][]int64{
				{1, 5, 2, 4},
				{4, 9, 5, 3},
				{4, 5, 2, 3},
				{1, 5, 5, 2},
			},
			qs: []query{
				{1, 1, 4, 4},
				{2, 2, 3, 3},
				{1, 2, 4, 3},
			},
		},
		{
			n: 1, q: 1,
			mat: [][]int64{{7}},
			qs:  []query{{1, 1, 1, 1}},
		},
	}
}

func randomMatrix(rng *rand.Rand, n int) [][]int64 {
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			mat[i][j] = int64(rng.Intn(20) + 1)
		}
	}
	return mat
}

func randomQueries(rng *rand.Rand, n, q int) []query {
	res := make([]query, q)
	for i := 0; i < q; i++ {
		x1 := rng.Intn(n) + 1
		x2 := rng.Intn(n-x1+1) + x1
		y1 := rng.Intn(n) + 1
		y2 := rng.Intn(n-y1+1) + y1
		res[i] = query{x1, y1, x2, y2}
	}
	return res
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	totalQ := 0
	for len(tests) < 50 && totalN < 2000 && totalQ < 5000 {
		n := rng.Intn(20) + 1
		q := rng.Intn(30) + 1
		if totalN+n > 2000 {
			n = 2000 - totalN
		}
		if totalQ+q > 5000 {
			q = 5000 - totalQ
		}
		if n <= 0 || q <= 0 {
			break
		}
		tests = append(tests, testCase{
			n:   n,
			q:   q,
			mat: randomMatrix(rng, n),
			qs:  randomQueries(rng, n, q),
		})
		totalN += n
		totalQ += q
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", tc.mat[i][j]))
			}
			sb.WriteByte('\n')
		}
		for _, qu := range tc.qs {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", qu.x1, qu.y1, qu.x2, qu.y2))
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
	tmpDir, err := os.MkdirTemp("", "2044H-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2044H")
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

func parseOutputs(out string, totalQueries int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != totalQueries {
		return nil, fmt.Errorf("expected %d outputs, got %d", totalQueries, len(tokens))
	}
	res := make([]int64, totalQueries)
	for i, tok := range tokens {
		if _, err := fmt.Sscan(tok, &res[i]); err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", tok, err)
		}
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

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

	totalQueries := 0
	for _, tc := range tests {
		totalQueries += tc.q
	}

	expected, err := parseOutputs(refOut, totalQueries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, totalQueries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "query %d mismatch: expected %d got %d\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
