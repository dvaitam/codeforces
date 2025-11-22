package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource2120E = "2000-2999/2100-2199/2120-2129/2120/2120E.go"
	maxTotalN      = 200000
)

type testCase struct {
	k int64
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2120E)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeTests(tests)

	refAns, err := runAndParse(refBin, tests, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	candAns, err := runAndParse(candidate, tests, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d, got %d\ninput:\n%s", i+1, refAns[i], candAns[i], singleTestInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2120E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	srcPath, err := resolveSourcePath(source)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func resolveSourcePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, path), nil
}

func runAndParse(target string, tests []testCase, input string) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	ans, err := parseOutput(out, len(tests))
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutput(out string, t int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	res := make([]int64, 0, t)
	for i := 0; i < t; i++ {
		var v int64
		if _, err := fmt.Fscan(reader, &v); err != nil {
			return nil, fmt.Errorf("expected %d answers, got %d (%v)", t, i, err)
		}
		res = append(res, v)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected starting with %q", extra)
	}
	return res, nil
}

func serializeTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", len(tc.a), tc.k))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func singleTestInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", len(tc.a), tc.k))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	totalN := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalN += len(tc.a)
	}

	// Deterministic small cases (including sample-like).
	add(testCase{k: 4, a: []int64{13, 7, 4}})
	add(testCase{k: 9, a: []int64{6, 12, 14}})
	add(testCase{k: 5, a: []int64{5, 3}})
	add(testCase{k: 7, a: []int64{6}})
	add(testCase{k: 1, a: []int64{1, 1, 1, 1, 1}})
	add(testCase{k: 1000000, a: []int64{1, 1000000}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Random mid-sized cases.
	for totalN < maxTotalN && len(tests) < 25 {
		n := rng.Intn(5000) + 1
		if totalN+n > maxTotalN {
			n = maxTotalN - totalN
		}
		k := int64(rng.Intn(1_000_000) + 1)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(rng.Intn(1_000_000) + 1)
		}
		add(testCase{k: k, a: a})
	}

	// Large stress case to hit limits.
	if totalN < maxTotalN {
		n := maxTotalN - totalN
		k := int64(rng.Intn(1_000_000) + 1)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			// Keep values moderate but varied.
			a[i] = int64(rng.Intn(2000) + 1)
		}
		add(testCase{k: k, a: a})
	}

	return tests
}
