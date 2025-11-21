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
	"time"
)

type testCase struct {
	input string
	id    string
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]
	if candidate == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
		os.Exit(1)
	}

	baseDir := currentDir()
	refBin, err := buildReference(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if normalize(got) != normalize(exp) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %s\nInput:\n%sExpected: %s\nGot: %s\n", tc.id, tc.input, exp, got)
			os.Exit(1)
		}
		if (i+1)%10 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d/%d tests...\n", i+1, len(tests))
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get current file location")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref324C1.bin")
	cmd := exec.Command("go", "build", "-o", out, "324C1.go")
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func generateTests() []testCase {
	var tests []testCase
	deterministic := []int{
		0, 1, 2, 5, 9, 10, 11, 19, 24, 37, 100, 101, 111, 191, 999, 1010,
		2020, 9090, 9999, 12345, 54321, 99999, 100000, 123456, 654321, 999999, 1000000,
	}
	for idx, v := range deterministic {
		tests = append(tests, makeTestCase(fmt.Sprintf("fixed-%02d", idx+1), v))
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		val := rng.Intn(1_000_001)
		tests = append(tests, makeTestCase(fmt.Sprintf("rand-%02d", i+1), val))
	}

	// sequential range to stress DP reuse (small numbers)
	for v := 0; v <= 200; v++ {
		tests = append(tests, makeTestCase(fmt.Sprintf("seq-%03d", v), v))
	}
	return tests
}

func makeTestCase(id string, value int) testCase {
	return testCase{
		id:    id,
		input: fmt.Sprintf("%d\n", value),
	}
}
