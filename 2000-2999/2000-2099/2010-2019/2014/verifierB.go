package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int64
	k int64
}

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/solution")
		os.Exit(1)
	}
	target := args[0]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicCases(), randomCases()...)
	input := buildInput(tests)

	expectedOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\noutput:\n%s", err, expectedOut)
		os.Exit(1)
	}
	candOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	expTokens, err := parseAnswers(expectedOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse oracle output: %v\noutput:\n%s", err, expectedOut)
		os.Exit(1)
	}
	candTokens, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range expTokens {
		if candTokens[i] != expTokens[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s, got %s\ninput:\n%s", i+1, expTokens[i], candTokens[i], input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2014B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2014B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
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
		return stdout.String() + stderr.String(), fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	}
	return sb.String()
}

func deterministicCases() []testCase {
	return []testCase{
		{n: 1, k: 1},
		{n: 2, k: 1},
		{n: 2, k: 2},
		{n: 3, k: 2},
		{n: 4, k: 4},
		{n: 5, k: 1},
		{n: 6, k: 3},
		{n: 100, k: 1},
		{n: 100, k: 100},
		{n: 999999937, k: 1},
		{n: 1000000000, k: 1000000000},
	}
}

func randomCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 200)
	for len(tests) < 200 {
		var n int64
		switch {
		case len(tests) < 50:
			n = int64(rng.Intn(100) + 1)
		case len(tests) < 120:
			n = int64(rng.Intn(10000) + 1)
		default:
			n = int64(rng.Intn(1000000000) + 1)
		}
		k := int64(rng.Intn(int(n)) + 1)
		tests = append(tests, testCase{n: n, k: k})
	}
	return tests
}

func parseAnswers(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		val := strings.ToUpper(tok)
		if val != "YES" && val != "NO" {
			return nil, fmt.Errorf("token %d invalid: %q", i+1, tok)
		}
		res[i] = val
	}
	return res, nil
}
