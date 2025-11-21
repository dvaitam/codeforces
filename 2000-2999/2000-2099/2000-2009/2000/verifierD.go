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

const referenceSolutionRel = "2000-2999/2000-2099/2000-2009/2000/2000D.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2000D.go")
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
	name string
	a    []int64
	s    string
}

func formatTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		n := len(tc.a)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name: "sample1",
			a:    []int64{3, 5, 1, 4, 3, 2},
			s:    "LRLLLR",
		},
		{
			name: "sample2",
			a:    []int64{2, 8},
			s:    "LR",
		},
		{
			name: "sample3",
			a:    []int64{3, 9},
			s:    "RL",
		},
		{
			name: "sample4",
			a:    []int64{1, 2, 3, 4, 5},
			s:    "LRLRR",
		},
		{
			name: "allL_then_R",
			a:    []int64{5, 4, 3, 2, 1},
			s:    "LLLRR",
		},
		{
			name: "allR_then_L",
			a:    []int64{5, 4, 3, 2, 1},
			s:    "RRRLL",
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	for i := 0; i < 80; i++ {
		n := rng.Intn(50) + 2
		if totalN+n > 200000 {
			break
		}
		totalN += n
		a := make([]int64, n)
		builder := strings.Builder{}
		for j := 0; j < n; j++ {
			a[j] = int64(rng.Intn(100000) + 1)
			if rng.Intn(2) == 0 {
				builder.WriteByte('L')
			} else {
				builder.WriteByte('R')
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			a:    a,
			s:    builder.String(),
		})
	}
	return tests
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
	tmpDir, err := os.MkdirTemp("", "2000D-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2000D")
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

func parseOutputs(output string, t int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	results := make([]int64, t)
	for i := 0; i < t; i++ {
		var val int64
		if _, err := fmt.Sscan(tokens[i], &val); err != nil {
			return nil, fmt.Errorf("failed to parse integer %q: %v", tokens[i], err)
		}
		results[i] = val
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra tokens after reading %d answers", t)
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatTests(tests)

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
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if got[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "test %s (%d) failed: expected %d, got %d\ninput:\n%s", tc.name, i+1, expected[i], got[i], formatTests([]testCase{tc}))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
