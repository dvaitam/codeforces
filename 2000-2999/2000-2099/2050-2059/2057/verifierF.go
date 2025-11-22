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
	refSource2057F = "2000-2999/2000-2099/2050-2059/2057/2057F.go"
	maxTotalN      = 20000
	maxTotalQ      = 20000
)

type testCase struct {
	n  int
	q  int
	a  []int64
	ks []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2057F)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)
	totalAns := totalQueries(tests)

	expected, err := runAndParse(refBin, input, totalAns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, totalAns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	idx := 0
	for tIdx, tc := range tests {
		for qIdx := 0; qIdx < tc.q; qIdx++ {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "Mismatch in test %d query %d: expected %d, got %d\n", tIdx+1, qIdx+1, expected[idx], got[idx])
				fmt.Fprintf(os.Stderr, "Test case input:\n%s", formatCase(tc))
				os.Exit(1)
			}
			idx++
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2057F-ref-*")
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
	sumN := totalN(tests)
	sumQ := totalQueries(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for sumN < maxTotalN && sumQ < maxTotalQ {
		remainN := maxTotalN - sumN
		remainQ := maxTotalQ - sumQ
		n := rng.Intn(min(2000, remainN)) + 1
		q := rng.Intn(min(2000, remainQ)) + 1
		a := randomComfortableArray(rng, n)
		ks := make([]int64, q)
		for i := 0; i < q; i++ {
			ks[i] = rng.Int63n(1_000_000_000) + 1
		}
		tests = append(tests, testCase{n: n, q: q, a: a, ks: ks})
		sumN += n
		sumQ += q
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, q: 3, a: []int64{1}, ks: []int64{1, 5, 10}},
		{n: 2, q: 4, a: []int64{5, 8}, ks: []int64{1, 2, 4, 10}},
		{n: 3, q: 3, a: []int64{10, 15, 20}, ks: []int64{1, 5, 100}},
		{n: 5, q: 5, a: []int64{1, 1, 1, 1, 1}, ks: []int64{1, 2, 4, 8, 16}},
		{n: 6, q: 6, a: []int64{100, 150, 225, 330, 490, 730}, ks: []int64{1, 2, 3, 10, 100, 1000}},
	}
}

func randomComfortableArray(rng *rand.Rand, n int) []int64 {
	a := make([]int64, n)
	a[0] = rng.Int63n(1_000_000_000) + 1
	for i := 1; i < n; i++ {
		maxVal := minInt64(2*a[i-1], 1_000_000_000)
		if maxVal < 1 {
			maxVal = 1
		}
		a[i] = rng.Int63n(maxVal) + 1
	}
	return a
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, k := range tc.ks {
			sb.WriteString(strconv.FormatInt(k, 10))
			if i+1 < tc.q {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
}

func totalQueries(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.q
	}
	return sum
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, k := range tc.ks {
		sb.WriteString(strconv.FormatInt(k, 10))
		if i+1 < tc.q {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
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
