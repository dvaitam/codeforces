package main

import (
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
	refSource2002E = "2000-2999/2000-2099/2000-2009/2002/2002E.go"
	maxTotalPairs  = 300000
)

type pair struct {
	a int64
	b int
}

type testCase struct {
	pairs []pair
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2002E)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)
	answerCount := totalPairs(tests)

	expected, err := runAndParse(refBin, input, answerCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, answerCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	idx := 0
	for tIdx, tc := range tests {
		for pIdx := range tc.pairs {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "Mismatch in test %d prefix %d: expected %d, got %d\n", tIdx+1, pIdx+1, expected[idx], got[idx])
				fmt.Fprintf(os.Stderr, "Test case input:\n%s", formatCase(tc))
				os.Exit(1)
			}
			idx++
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2002E-ref-*")
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
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runAndParse(target, input string, expectedCount int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(out)
	res := make([]int64, expectedCount)
	for i := 0; i < expectedCount; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("expected %d numbers, got %d (%v)", expectedCount, i, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected starting with %s", extra)
	}
	return res, nil
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

func buildTests() []testCase {
	tests := deterministicTests()
	sum := totalPairs(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for sum < maxTotalPairs {
		left := maxTotalPairs - sum
		if left <= 0 {
			break
		}
		n := rng.Intn(min(5000, left)) + 1
		tests = append(tests, randomCase(rng, n))
		sum += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{pairs: []pair{{a: 1, b: 0}}},
		{pairs: []pair{{a: 3, b: 1}, {a: 2, b: 0}, {a: 4, b: 2}}},
		{pairs: []pair{{a: 5, b: 1}, {a: 1, b: 2}, {a: 5, b: 1}}},
		{pairs: []pair{{a: 1e9, b: 0}, {a: 2, b: 1}, {a: 1, b: 2}, {a: 1e9, b: 0}}},
	}
}

func randomCase(rng *rand.Rand, n int) testCase {
	pairs := make([]pair, n)
	prev := -1
	for i := 0; i < n; i++ {
		a := rng.Int63n(1_000_000_000) + 1
		b := rng.Intn(n + 1)
		if i > 0 && b == prev {
			b = (b + 1) % (n + 1)
		}
		pairs[i] = pair{a: a, b: b}
		prev = b
	}
	return testCase{pairs: pairs}
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.pairs)))
		sb.WriteByte('\n')
		for _, p := range tc.pairs {
			sb.WriteString(strconv.FormatInt(p.a, 10))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(p.b))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func totalPairs(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.pairs)
	}
	return total
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(len(tc.pairs)))
	sb.WriteByte('\n')
	for _, p := range tc.pairs {
		sb.WriteString(strconv.FormatInt(p.a, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(p.b))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func resolveSourcePath(path string) (string, error) {
	cleaned := filepath.Clean(path)
	if abs, err := filepath.Abs(cleaned); err == nil {
		if _, err := os.Stat(abs); err == nil {
			return abs, nil
		}
	}
	base := filepath.Base(path)
	if abs, err := filepath.Abs(base); err == nil {
		if _, err := os.Stat(abs); err == nil {
			return abs, nil
		}
	}
	return "", fmt.Errorf("reference source not found at %s", path)
}
