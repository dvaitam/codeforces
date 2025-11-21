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
	c int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-690D2-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleD2")
	cmd := exec.Command("go", "build", "-o", bin, "690D2.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
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

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n", tc.n, tc.c)
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d: %q", len(fields), out)
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 0 || val >= 1000003 {
		// modulo result must lie in [0, mod)
		return 0, fmt.Errorf("result %d outside modulo range", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, c: 1},
		{n: 1, c: 2},
		{n: 2, c: 1},
		{n: 3, c: 2},
		{n: 5, c: 5},
		{n: 10, c: 1},
		{n: 500000, c: 1},
		{n: 500000, c: 200000},
	}
}

func randomTest(rng *rand.Rand) testCase {
	var n int
	var c int
	mode := rng.Intn(4)
	switch mode {
	case 0:
		n = rng.Intn(10) + 1
		c = rng.Intn(10) + 1
	case 1:
		n = rng.Intn(1000) + 1
		c = rng.Intn(1000) + 1
	case 2:
		n = rng.Intn(500000) + 1
		c = rng.Intn(200000) + 1
	default:
		n = 500000 - rng.Intn(1000)
		if n < 1 {
			n = 1
		}
		c = 200000 - rng.Intn(1000)
		if c < 1 {
			c = 1
		}
	}
	return testCase{n: n, c: c}
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
	for i := 0; i < 300; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expStr)
			os.Exit(1)
		}

		gotStr, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotStr, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
