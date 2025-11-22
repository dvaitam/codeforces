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
	refSource2127G2 = "2000-2999/2100-2199/2120-2129/2127/2127G2.go"
	maxTotalSq      = 9500 // stay below the 1e4 bound comfortably
)

type testCase struct {
	n int
	p []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2127G2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeTests(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	expect, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output invalid: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if err := compareCase(expect[i], got[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, singleTestInput(tests[i]), candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2127G2-ref-*")
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

func parseOutputs(out string, tests []testCase) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([][]int, len(tests))
	for idx, tc := range tests {
		ans[idx] = make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(reader, &ans[idx][i]); err != nil {
				return nil, fmt.Errorf("test %d: expected %d numbers, got %d (%v)", idx+1, tc.n, i, err)
			}
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected starting with %q", extra)
	}
	return ans, nil
}

func compareCase(expect, got []int) error {
	if len(expect) != len(got) {
		return fmt.Errorf("length mismatch: expected %d values, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			return fmt.Errorf("position %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}
	return nil
}

func serializeTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func singleTestInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0)
	totalSq := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		totalSq += tc.n * tc.n
	}

	// Deterministic small cases.
	add(testCase{n: 4, p: []int{2, 3, 4, 1}})
	add(testCase{n: 5, p: []int{2, 3, 5, 1, 4}})
	add(testCase{n: 6, p: []int{2, 1, 4, 3, 6, 5}})
	add(testCase{n: 10, p: derangement(10, rng)})
	add(testCase{n: 15, p: derangement(15, rng)})

	// Random cases up to the squared budget.
	for totalSq < maxTotalSq-100 { // leave room for one large case
		n := rng.Intn(25) + 4
		if totalSq+(n*n) > maxTotalSq-100 {
			break
		}
		add(testCase{n: n, p: derangement(n, rng)})
	}

	// One larger stress case if budget allows.
	for size := 50; size >= 30; size-- {
		if totalSq+size*size <= maxTotalSq {
			add(testCase{n: size, p: derangement(size, rng)})
			break
		}
	}

	return tests
}

func derangement(n int, rng *rand.Rand) []int {
	for {
		perm := rng.Perm(n)
		ok := true
		for i, v := range perm {
			if v == i {
				ok = false
				break
			}
		}
		if ok {
			for i := 0; i < n; i++ {
				perm[i]++
			}
			return perm
		}
	}
}
