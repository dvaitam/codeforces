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

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-538B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "538B.go")
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

func generateTests() []string {
	tests := []string{
		"1",
		"9",
		"32",
		"101",
		"999999",
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		n := rng.Intn(1_000_000) + 1
		tests = append(tests, strconv.Itoa(n))
	}
	return tests
}

func isQuasiBinary(s string) bool {
	for _, ch := range s {
		if ch != '0' && ch != '1' {
			return false
		}
	}
	return len(s) > 0
}

func parseOutput(output string) (int, []string, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return 0, nil, fmt.Errorf("expected at least 2 lines, got %v", output)
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse k: %v", err)
	}
	parts := strings.Fields(lines[1])
	if len(parts) != k {
		return 0, nil, fmt.Errorf("expected %d numbers, got %v", k, lines[1])
	}
	return k, parts, nil
}

func validate(nStr string, parts []string) error {
	if len(parts) == 0 {
		return fmt.Errorf("no parts provided")
	}
	var sum int
	for _, p := range parts {
		if !isQuasiBinary(p) {
			return fmt.Errorf("not quasi-binary: %s", p)
		}
		val, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("invalid number %s: %v", p, err)
		}
		sum += val
	}
	n, _ := strconv.Atoi(nStr)
	if sum != n {
		return fmt.Errorf("sum %d != %d", sum, n)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, nStr := range tests {
		input := nStr + "\n"

		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\n", idx+1, nStr, err)
			os.Exit(1)
		}
		expK, expParts, err := parseOutput(expected)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expected)
			os.Exit(1)
		}
		expLen := expK

		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d (%s): %v\n", idx+1, nStr, err)
			os.Exit(1)
		}
		actK, actParts, err := parseOutput(actual)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d (%s): %v\n%s", idx+1, nStr, err, actual)
			os.Exit(1)
		}
		if actK < expLen {
			fmt.Fprintf(os.Stderr, "test %d (%s): k too small (expected >= %d, got %d)\n", idx+1, nStr, expLen, actK)
			os.Exit(1)
		}
		if err := validate(nStr, actParts); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed validation: %v\noutput: %s\n", idx+1, nStr, err, actual)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
