package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2171E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "2171E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "example1", input: "88 94 95\n"},
		{name: "example2", input: "100 80 81\n"},
		{name: "example3", input: "98 99 98\n"},
		{name: "example4", input: "95 86 85\n"},
		{name: "all_equal_low", input: "80 80 80\n"},
		{name: "all_equal_high", input: "100 100 100\n"},
		{name: "diff_exact_9", input: "80 89 85\n"},
		{name: "diff_exact_10", input: "90 100 95\n"},
		{name: "ascending", input: "80 85 90\n"},
		{name: "descending", input: "100 95 90\n"},
		{name: "mixed_mid", input: "81 90 100\n"},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	for i := 0; i < 100; i++ {
		g := rng.Intn(21) + 80
		c := rng.Intn(21) + 80
		l := rng.Intn(21) + 80
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: fmt.Sprintf("%d %d %d\n", g, c, l),
		})
	}
	return tests
}

func exhaustiveTests() []testCase {
	tests := make([]testCase, 0, 9261)
	for g := 80; g <= 100; g++ {
		for c := 80; c <= 100; c++ {
			for l := 80; l <= 100; l++ {
				tests = append(tests, testCase{
					name:  fmt.Sprintf("exhaustive_%d_%d_%d", g, c, l),
					input: fmt.Sprintf("%d %d %d\n", g, c, l),
				})
			}
		}
	}
	return tests
}

func medianResult(input string) string {
	fields := strings.Fields(input)
	if len(fields) != 3 {
		return ""
	}
	vals := make([]int, 3)
	for i := 0; i < 3; i++ {
		fmt.Sscanf(fields[i], "%d", &vals[i])
	}
	sort.Ints(vals)
	if vals[2]-vals[0] >= 10 {
		return "check again"
	}
	return fmt.Sprintf("final %d", vals[1])
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	tests = append(tests, exhaustiveTests()...)

	for idx, tc := range tests {
		expected, err := runBinary(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		actual, err := runBinary(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if expected != actual {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %q got %q\ninput:\n%s", idx+1, tc.name, expected, actual, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
