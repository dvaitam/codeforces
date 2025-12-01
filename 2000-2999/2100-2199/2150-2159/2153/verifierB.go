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
	refSource        = "2153B.go"
	tempOraclePrefix = "oracle-2153B-"
	randomTests      = 80
	maxValue         = 1 << 16
)

type testCase struct {
	x int
	y int
	z int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		input := formatInput(tc)
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		exp := normalize(expOut)
		got := normalize(candOut)
		if exp == "" || got == "" {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", idx+1)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(candOut)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", idx+1, exp, got)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Expected:")
			fmt.Print(expOut)
			fmt.Println("Got:")
			fmt.Print(candOut)
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
	outPath := filepath.Join(tmpDir, "oracleB")
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

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func formatInput(tc testCase) string {
	return fmt.Sprintf("1\n%d %d %d\n", tc.x, tc.y, tc.z)
}

func normalize(out string) string {
	out = strings.TrimSpace(strings.ToLower(out))
	if out == "yes" {
		return "yes"
	}
	if out == "no" {
		return "no"
	}
	return ""
}

func deterministicTests() []testCase {
	return []testCase{
		{x: 0, y: 0, z: 0},
		{x: 1, y: 1, z: 1},
		{x: 3, y: 2, z: 6},
		{x: 4, y: 8, z: 12},
		{x: 5, y: 10, z: 15},
	}
}

func randomTestsCases(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		x := rng.Intn(maxValue)
		y := rng.Intn(maxValue)
		z := rng.Intn(maxValue)
		tests = append(tests, testCase{x: x, y: y, z: z})
	}
	return tests
}
