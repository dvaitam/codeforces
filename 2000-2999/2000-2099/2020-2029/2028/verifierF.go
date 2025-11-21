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

const referenceSolutionRel = "2000-2999/2000-2099/2020-2029/2028/2028F.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2028F.go")
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
	n int
	m int
	a []int
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 5, m: 4, a: []int{2, 1, 1, 1, 2}},
		{n: 5, m: 5, a: []int{2, 1, 1, 1, 2}},
		{n: 5, m: 6, a: []int{2, 1, 1, 1, 2}},
		{n: 5, m: 7, a: []int{2, 1, 1, 1, 2}},
		{n: 5, m: 8, a: []int{2, 1, 1, 1, 2}},
		{n: 5, m: 6, a: []int{2, 0, 2, 2, 3}},
		{n: 1, m: 0, a: []int{0}},
		{n: 1, m: 5, a: []int{5}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	for len(tests) < 80 && totalN < 4000 {
		n := rng.Intn(20) + 1
		if totalN+n > 4000 {
			n = 4000 - totalN
		}
		m := rng.Intn(100) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if rng.Float64() < 0.2 {
				a[i] = 0
			} else {
				a[i] = rng.Intn(20)
			}
		}
		tests = append(tests, testCase{n: n, m: m, a: a})
		totalN += n
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
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
	tmpDir, err := os.MkdirTemp("", "2028F-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2028F")
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

func parseAnswers(out string, count int) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	ans := make([]string, 0, count)
	for scanner.Scan() {
		ans = append(ans, strings.ToUpper(scanner.Text()))
	}
	if len(ans) != count {
		return nil, fmt.Errorf("expected %d answers, got %d", count, len(ans))
	}
	return ans, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
	expected, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseAnswers(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\nn=%d m=%d a=%v\n", i+1, expected[i], got[i], tests[i].n, tests[i].m, tests[i].a)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
