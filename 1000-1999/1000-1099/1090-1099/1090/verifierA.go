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
	referenceSource = "./1090A.go"
	testCases       = 120
)

type test struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary-or-source")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expectOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expectVal, err := parseSingleInt(expectOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, expectOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		gotVal, err := parseSingleInt(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}
		if gotVal != expectVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, tc.name, expectVal, gotVal, tc.input, expectOut, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	tmp, err := os.CreateTemp("", "ref1090A-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmpPath, referenceSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpPath)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() { os.Remove(tmpPath) }
	return tmpPath, cleanup, nil
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runCandidate(target string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		abs, err := filepath.Abs(target)
		if err != nil {
			return "", fmt.Errorf("resolve candidate path: %w", err)
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseSingleInt(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer %d", val)
	}
	return val, nil
}

func generateTests() []test {
	var tests []test
	tests = append(tests, test{
		name:  "sample",
		input: "3\n2 4 3\n2 2 1\n3 1 1 1\n",
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < testCases {
		tests = append(tests, randomTest(rng, len(tests)))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) test {
	n := rng.Intn(8) + 1
	totalEmployees := 0
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		remaining := 200000 - totalEmployees
		maxEmployees := min(5000, remaining-(n-i-1))
		if maxEmployees <= 0 {
			maxEmployees = 1
		}
		m := rng.Intn(min(2000, maxEmployees)) + 1
		totalEmployees += m
		b.WriteString(fmt.Sprintf("%d", m))
		for j := 0; j < m; j++ {
			salary := rng.Intn(1_000_000_000) + 1
			b.WriteString(fmt.Sprintf(" %d", salary))
		}
		b.WriteByte('\n')
	}
	return test{
		name:  fmt.Sprintf("rand_%d", idx),
		input: b.String(),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
