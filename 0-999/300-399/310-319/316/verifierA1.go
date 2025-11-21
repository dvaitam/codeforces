package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	refSource        = "316A1.go"
	tempOraclePrefix = "oracle-316A1-"
	randomTests      = 500
	maxLen           = 5
)

type testCase struct {
	hint string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestCases(randomTests, rng)...)

	for idx, tc := range tests {
		input := tc.hint + "\n"
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (s=%q): %v\n", idx+1, tc.hint, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (s=%q): %v\n", idx+1, tc.hint, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expVal, err := parseBigInt(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotVal, err := parseBigInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		if expVal.Cmp(gotVal) != 0 {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s (s=%q)\n", idx+1, expVal.String(), gotVal.String(), tc.hint)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA1")
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

func parseBigInt(out string) (*big.Int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	val, ok := new(big.Int).SetString(fields[0], 10)
	if !ok {
		return nil, fmt.Errorf("value %q is not an integer", fields[0])
	}
	if len(fields) > 1 {
		return nil, fmt.Errorf("unexpected extra tokens in output")
	}
	return val, nil
}

func deterministicTests() []testCase {
	patterns := []string{
		"?",
		"A",
		"1",
		"9",
		"A?",
		"?A",
		"AB",
		"AA",
		"A?A",
		"A1",
		"A1?",
		"A1B2C",
		"ABCDE",
		"?????",
		"JJJJJ",
		"AABBC",
		"??A??",
		"A2345",
		"A?9?A",
		"B0?A?", // first char not zero (B)
	}
	tests := make([]testCase, 0, len(patterns))
	for _, p := range patterns {
		tests = append(tests, testCase{hint: p})
	}
	return tests
}

func randomTestCases(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	firstChars := []rune("?123456789ABCDEFGHIJ")
	allChars := []rune("?0123456789ABCDEFGHIJ")
	for len(tests) < count {
		length := rng.Intn(maxLen) + 1
		var builder strings.Builder
		for i := 0; i < length; i++ {
			if i == 0 {
				builder.WriteRune(firstChars[rng.Intn(len(firstChars))])
			} else {
				builder.WriteRune(allChars[rng.Intn(len(allChars))])
			}
		}
		tests = append(tests, testCase{hint: builder.String()})
	}
	return tests
}
