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

const referenceSolutionRel = "2000-2999/2000-2099/2030-2039/2032/2032E.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2032E.go")
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
	a []int64
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, a: []int64{5}},
		{n: 3, a: []int64{2, 1, 2}},
		{n: 5, a: []int64{2, 1, 1, 1, 2}},
		{n: 7, a: []int64{1, 2, 1, 2, 1, 2, 1}},
		{n: 9, a: []int64{10, 1, 2, 3, 4, 5, 6, 7, 8}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	for len(tests) < 80 && totalN < 2000 {
		n := rng.Intn(19)*2 + 1
		if totalN+n > 2000 {
			n = (totalN + n) - 2000
			if n%2 == 0 {
				n++
			}
		}
		if n <= 0 {
			break
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(rng.Intn(50) + 1)
		}
		tests = append(tests, testCase{n: n, a: a})
		totalN += n
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
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
	tmpDir, err := os.MkdirTemp("", "2032E-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2032E")
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

func parseOutput(out string, expectedLines int) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != expectedLines {
		return nil, fmt.Errorf("expected %d lines, got %d", expectedLines, len(lines))
	}
	return lines, nil
}

func simulate(a []int64, ops []int64) ([]int64, error) {
	n := len(a)
	res := append([]int64(nil), a...)
	for i, times := range ops {
		if times < 0 {
			return nil, fmt.Errorf("negative operations at index %d", i+1)
		}
		left := (i - 1 + n) % n
		right := (i + 1) % n
		res[left] += times
		res[i] += 2 * times
		res[right] += times
	}
	return res, nil
}

func checkLine(line string, n int) ([]int64, bool) {
	fields := strings.Fields(line)
	if len(fields) == 1 && fields[0] == "-1" {
		return nil, true
	}
	if len(fields) != n {
		return nil, false
	}
	ops := make([]int64, n)
	for i, f := range fields {
		var val int64
		if _, err := fmt.Sscan(f, &val); err != nil {
			return nil, false
		}
		ops[i] = val
	}
	return ops, false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	refLines, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	userLines, err := parseOutput(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "participant output parse error: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		_, refImpossible := checkLine(refLines[idx], tc.n)
		userOps, userImpossible := checkLine(userLines[idx], tc.n)
		if refImpossible {
			if !userImpossible {
				if userOps != nil {
					fmt.Fprintf(os.Stderr, "test %d: expected -1 but participant gave operations\n", idx+1)
					os.Exit(1)
				}
				fmt.Fprintf(os.Stderr, "test %d: expected -1 but participant gave other output\n", idx+1)
				os.Exit(1)
			}
			continue
		}
		if userImpossible {
			fmt.Fprintf(os.Stderr, "test %d: solution exists but participant output -1\n", idx+1)
			os.Exit(1)
		}
		final, err := simulate(tc.a, userOps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid operations: %v\n", idx+1, err)
			os.Exit(1)
		}
		for i := 1; i < len(final); i++ {
			if final[i] != final[0] {
				fmt.Fprintf(os.Stderr, "test %d: not balanced after applying operations\n", idx+1)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
