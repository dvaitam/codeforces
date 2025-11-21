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
	strings []string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1551B1-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleB1")
	cmd := exec.Command("go", "build", "-o", path, "1551B1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string, input string) (string, error) {
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

func parseOutput(out string, t int) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d lines, got %d", t, len(lines))
	}
	res := make([]int, t)
	for i := 0; i < t; i++ {
		val, err := strconv.Atoi(strings.TrimSpace(lines[i]))
		if err != nil {
			return nil, fmt.Errorf("invalid integer on line %d: %v", i+1, err)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative answer on line %d", i+1)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{strings: []string{"kzaaa"}},
		{strings: []string{"codeforces", "abacaba"}},
		{strings: []string{"a", "aaaaa", "abcde"}},
		{strings: []string{"zzzzzzzzzz", "qwerty"}},
	}
}

func randomString(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte(rng.Intn(26)) + 'a')
	}
	return sb.String()
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(10) + 1
	strs := make([]string, t)
	for i := 0; i < t; i++ {
		switch rng.Intn(4) {
		case 0:
			strs[i] = randomString(rng)
		case 1:
			strs[i] = strings.Repeat("a", rng.Intn(50)+1)
		case 2:
			strs[i] = "abcdefghijklmnopqrstuvwxyz"[:rng.Intn(26)+1]
		default:
			strs[i] = "codeforces"
		}
	}
	return testCase{strings: strs}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.strings)))
	for _, s := range tc.strings {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expVals, err := parseOutput(expOut, len(tc.strings))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, len(tc.strings))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		for i := range expVals {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d case %d: expected %d got %d\ninput:\n%s\n", idx+1, i+1, expVals[i], gotVals[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
