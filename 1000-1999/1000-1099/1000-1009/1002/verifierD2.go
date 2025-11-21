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
	n int
	b string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1002D2-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleD2")
	cmd := exec.Command("go", "build", "-o", path, "1002D2.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
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
	return stdout.String(), nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d\n%s\n", tc.n, tc.b)
}

func parseOutput(out string) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	cnt, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid operation count: %v", err)
	}
	if cnt != len(lines)-1 {
		return nil, fmt.Errorf("operation count mismatch: declared %d got %d", cnt, len(lines)-1)
	}
	ops := make([]string, cnt)
	for i := 0; i < cnt; i++ {
		ops[i] = strings.TrimSpace(lines[i+1])
	}
	return ops, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, b: "0"},
		{n: 1, b: "1"},
		{n: 2, b: "01"},
		{n: 3, b: "101"},
		{n: 4, b: "0000"},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return testCase{n: n, b: sb.String()}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
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

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}

		expOps, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}
		gotOps, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}
		if len(expOps) != len(gotOps) {
			fmt.Fprintf(os.Stderr, "test %d: mismatch in operation count: expected %d got %d\ninput:\n%s\n", idx+1, len(expOps), len(gotOps), input)
			os.Exit(1)
		}
		for i := range expOps {
			if gotOps[i] != expOps[i] {
				fmt.Fprintf(os.Stderr, "test %d: operation %d mismatch: expected %q got %q\ninput:\n%s\n", idx+1, i+1, expOps[i], gotOps[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
