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
	refSource  = "1000-1999/1800-1899/1860-1869/1866/1866D.go"
	totalTests = 60
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, tc.name, refVal, candVal, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref1866D-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1866D.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "single_element", input: "1 1 1\n5\n"},
		{name: "small_equal", input: "2 2 1\n1 2\n3 4\n"},
		{name: "k_equals_m", input: "3 3 3\n1 2 3\n4 5 6\n7 8 9\n"},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-2 {
		tests = append(tests, randomCase(rng, len(tests)+1))
	}
	tests = append(tests,
		randomLarge(rng, "large1", 10, 200),
		randomLarge(rng, "large2", 10, 500),
	)
	return tests
}

func randomCase(rng *rand.Rand, idx int) testCase {
	N := rng.Intn(5) + 1
	M := rng.Intn(8) + 1
	K := rng.Intn(min(10, M)) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", N, M, K))
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rng.Intn(100) + 1))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: fmt.Sprintf("rand_%d", idx), input: sb.String()}
}

func randomLarge(rng *rand.Rand, name string, N, M int) testCase {
	K := rng.Intn(min(10, M)) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", N, M, K))
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rng.Intn(1_000_000) + 1))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
