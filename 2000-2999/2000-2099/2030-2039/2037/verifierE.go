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
	refSource2037E = "2000-2999/2000-2099/2030-2039/2037/2037E.go"
	maxTotalN      = 10000
)

type testCase struct {
	n int
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2037E)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch in test %d: expected %s, got %s\n", i+1, expected[i], got[i])
			fmt.Fprintf(os.Stderr, "Test case input:\n%s", formatCase(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2037E-ref-*")
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

func runAndParse(target, input string, expected int) ([]string, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(out)
	res := make([]string, expected)
	for i := 0; i < expected; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("expected %d outputs, got %d (%v)", expected, i, err)
		}
		res[i] = strings.ToUpper(res[i])
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
	total := totalN(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for total < maxTotalN {
		remain := maxTotalN - total
		n := rng.Intn(min(5000, remain)) + 2
		s := randomString(rng, n)
		tests = append(tests, testCase{n: n, s: s})
		total += n
	}

	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, s: "01"},
		{n: 2, s: "00"},
		{n: 2, s: "11"},
		{n: 3, s: "010"},
		{n: 3, s: "101"},
		{n: 4, s: "1111"},
		{n: 5, s: "11000"},
		{n: 6, s: "011111"},
		{n: 7, s: "1000000"},
	}
}

func randomString(rng *rand.Rand, n int) string {
	if rng.Intn(4) == 0 {
		// Monotone string likely yields IMPOSSIBLE.
		ch := byte('0' + byte(rng.Intn(2)))
		return strings.Repeat(string(ch), n)
	}
	bs := make([]byte, n)
	for i := 0; i < n; i++ {
		bs[i] = byte('0' + rng.Intn(2))
	}
	return string(bs)
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
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

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	sb.WriteString(tc.s)
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
