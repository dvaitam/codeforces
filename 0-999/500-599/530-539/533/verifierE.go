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
	n    int
	s, t string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-533E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "533E.go")
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
	return stdout.String(), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(len(tc.s) + len(tc.t) + 16)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	sb.WriteString(tc.s)
	sb.WriteByte('\n')
	sb.WriteString(tc.t)
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens (%q)", len(fields), out)
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 7, s: "reading", t: "trading"},
		{n: 5, s: "sweet", t: "sheep"},
		{n: 3, s: "toy", t: "try"},
		{n: 1, s: "a", t: "b"},
		{n: 2, s: "aa", t: "bb"},
		{n: 6, s: "aaaaaa", t: "aaaaba"},
		{n: 6, s: "aaaaaa", t: "baaaaa"},
		{n: 6, s: "abcdef", t: "abceef"},
		{n: 10, s: "aaaaaaaaaa", t: "aaaaabaaaa"},
	}
}

func randomString(rng *rand.Rand, n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte('a' + rng.Intn(26))
	}
	return string(bytes)
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 200)
	for len(tests) < 200 {
		n := rng.Intn(20) + 1
		if rng.Intn(6) == 0 {
			n = rng.Intn(200) + 1
		}
		var s, t string
		for {
			s = randomString(rng, n)
			t = randomString(rng, n)
			if s != t {
				break
			}
		}
		tests = append(tests, testCase{n: n, s: s, t: t})
	}
	// add some worst-case pattern tests
	tests = append(tests,
		testCase{n: 1000, s: strings.Repeat("a", 1000), t: strings.Repeat("a", 999) + "b"},
		testCase{n: 1000, s: strings.Repeat("a", 999) + "b", t: "b" + strings.Repeat("a", 999)},
	)
	return tests
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
	for idx, tc := range tests {
		input := buildInput(tc)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expectedVal, err := parseOutput(expectedOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on test %d: %v\noutput:\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualVal, err := parseOutput(actualOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target produced invalid output on test %d: %v\noutput:\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		if expectedVal != actualVal {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\ninput:\n%s", idx+1, expectedVal, actualVal, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
