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

const referenceSolutionRel = "2000-2999/2000-2099/2040-2049/2043/2043C.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2043C.go")
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
	a []int64
}

func deterministicTests() []testCase {
	return []testCase{
		{a: []int64{-1, 1, 0, 1, -1}},
		{a: []int64{1, 1, 1}},
		{a: []int64{-1, -1, -1}},
		{a: []int64{5}},
		{a: []int64{1, 2, 3, 4}},
		{a: []int64{-1, -1, 10, -1, -1}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	for len(tests) < 80 && totalN < 2000 {
		n := rng.Intn(50) + 1
		if totalN+n > 2000 {
			n = 2000 - totalN
		}
		a := make([]int64, n)
		pos := rng.Intn(n + 1)
		for i := 0; i < n; i++ {
			if i == pos && rng.Intn(2) == 0 {
				a[i] = int64(rng.Intn(101) - 50)
			} else {
				if rng.Intn(2) == 0 {
					a[i] = -1
				} else {
					a[i] = 1
				}
			}
		}
		tests = append(tests, testCase{a: a})
		totalN += n
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.a)))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
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
	tmpDir, err := os.MkdirTemp("", "2043C-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2043C")
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

func parseOutputs(out string, cases int) ([][]int64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	results := make([][]int64, 0, cases)
	i := 0
	for caseIdx := 0; caseIdx < cases; caseIdx++ {
		if i >= len(lines) {
			return nil, fmt.Errorf("not enough lines for case %d", caseIdx+1)
		}
		var count int
		if _, err := fmt.Sscan(lines[i], &count); err != nil {
			return nil, fmt.Errorf("failed to parse count on line %d: %v", i+1, err)
		}
		i++
		if i >= len(lines) {
			return nil, fmt.Errorf("missing values line for case %d", caseIdx+1)
		}
		fields := strings.Fields(lines[i])
		if len(fields) != count {
			return nil, fmt.Errorf("case %d: expected %d sums, got %d", caseIdx+1, count, len(fields))
		}
		sums := make([]int64, count)
		for idx, f := range fields {
			if _, err := fmt.Sscan(f, &sums[idx]); err != nil {
				return nil, fmt.Errorf("case %d: failed to parse sum %d (%q)", caseIdx+1, idx+1, f)
			}
		}
		results = append(results, sums)
		i++
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
	expected, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference parse error: %v\n", err)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "participant parse error: %v\n", err)
		os.Exit(1)
	}

	for idx := range tests {
		if !equalSlices(expected[idx], got[idx]) {
			fmt.Fprintf(os.Stderr, "test %d mismatch:\nexpected %v\ngot %v\n", idx+1, expected[idx], got[idx])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func equalSlices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
