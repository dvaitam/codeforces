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

const (
	refSource        = "1357A3.go"
	tempOraclePrefix = "oracle-1357A3-"
	randomTests      = 50
)

type testCase struct {
	label string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA3.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestsCases(rng, randomTests)...)

	for idx, tc := range tests {
		exp, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.label, err)
			fmt.Println("Input:")
			fmt.Print(tc.input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.label, err)
			fmt.Println("Input:")
			fmt.Print(tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %s got %s\n", idx+1, tc.label, strings.TrimSpace(exp), strings.TrimSpace(got))
			fmt.Println("Input:")
			fmt.Print(tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA3")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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
		{label: "single_zero", input: "1\n0\n"},
		{label: "single_one", input: "1\n1\n"},
		{label: "empty", input: "0\n"},
	}
}

func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		k := rng.Intn(10)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", k)
		for j := 0; j < k; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteString("0\n")
			} else {
				sb.WriteString("1\n")
			}
		}
		tests = append(tests, testCase{
			label: fmt.Sprintf("random_%d", i+1),
			input: sb.String(),
		})
	}
	return tests
}
