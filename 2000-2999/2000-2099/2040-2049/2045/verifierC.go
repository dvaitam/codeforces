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
	s string
	t string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2045C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "2045C.go")
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

func deterministicTests() []testCase {
	return []testCase{
		{s: "sarana", t: "olahraga"},
		{s: "berhiber", t: "wortelhijau"},
		{s: "icpc", t: "icpc"},
		{s: "icpc", t: "jakarta"},
		{s: "aaaa", t: "aaaa"},
	}
}

func randomString(rng *rand.Rand, minLen, maxLen int) string {
	n := rng.Intn(maxLen-minLen+1) + minLen
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = byte('a' + rng.Intn(26))
	}
	return string(bytes)
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 80)
	for len(tests) < cap(tests) {
		s := randomString(rng, 1, 30)
		t := randomString(rng, 1, 30)
		tests = append(tests, testCase{s: s, t: t})
	}
	return tests
}

func isInteresting(s, t, abbr string) bool {
	if len(abbr) < 2 {
		return false
	}
	count := 0
	for i := 1; i < len(abbr); i++ {
		prefix := abbr[:i]
		suffix := abbr[i:]
		if strings.HasPrefix(s, prefix) && strings.HasSuffix(t, suffix) {
			count++
			if count >= 2 {
				return true
			}
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		input := fmt.Sprintf("%s\n%s\n", tc.s, tc.t)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\nS=%s\nT=%s", idx+1, err, tc.s, tc.t)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\nS=%s\nT=%s", idx+1, err, tc.s, tc.t)
			os.Exit(1)
		}
		expectedAns := strings.TrimSpace(expectedOut)
		actualAns := strings.TrimSpace(actualOut)

		if expectedAns == "-1" {
			if actualAns != "-1" {
				fmt.Fprintf(os.Stderr, "test %d: expected -1 but target output %q\nS=%s\nT=%s", idx+1, actualAns, tc.s, tc.t)
				os.Exit(1)
			}
			continue
		}

		if actualAns == "-1" {
			fmt.Fprintf(os.Stderr, "test %d: solution exists but target output -1\nS=%s\nT=%s", idx+1, tc.s, tc.t)
			os.Exit(1)
		}
		if len(actualAns) != len(expectedAns) {
			fmt.Fprintf(os.Stderr, "test %d: expected length %d but got %d\nS=%s\nT=%s\noracle=%s\ntarget=%s", idx+1, len(expectedAns), len(actualAns), tc.s, tc.t, expectedAns, actualAns)
			os.Exit(1)
		}
		if !isInteresting(tc.s, tc.t, actualAns) {
			fmt.Fprintf(os.Stderr, "test %d: target output is not an interesting abbreviation\nS=%s\nT=%s\nabbr=%s", idx+1, tc.s, tc.t, actualAns)
			os.Exit(1)
		}
		if actualAns != expectedAns {
			// ensure minimal by checking there is no shorter abbreviation (oracle gives minimal length).
			// Since lengths are equal, multiple answers allowed; nothing further needed.
		}
	}

	fmt.Println("All tests passed.")
}
