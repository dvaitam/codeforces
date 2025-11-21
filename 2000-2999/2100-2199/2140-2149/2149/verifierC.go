package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type testCase struct {
	n    int
	perm []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}

	tests, err := parseInput(inputBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid input:", err)
		os.Exit(1)
	}

	out, err := runProgram(candidate, inputBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := verifyOutput(out, tests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Accepted")
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("failed to read n: %v", err)
		}
		m := n*n + 1
		perm := make([]int, m)
		for j := 0; j < m; j++ {
			if _, err := fmt.Fscan(reader, &perm[j]); err != nil {
				return nil, fmt.Errorf("failed to read permutation value: %v", err)
			}
		}
		tests[i] = testCase{n: n, perm: perm}
	}
	return tests, nil
}

func verifyOutput(output string, tests []testCase) error {
	lines := readNonEmptyLines(output)
	if len(lines) != len(tests) {
		return fmt.Errorf("expected %d non-empty output lines, got %d", len(tests), len(lines))
	}
	for idx, line := range lines {
		tc := tests[idx]
		tokens := strings.Fields(line)
		if len(tokens) != tc.n+1 {
			return fmt.Errorf("test %d: expected %d indices, got %d", idx+1, tc.n+1, len(tokens))
		}
		indices := make([]int, tc.n+1)
		for i, tok := range tokens {
			v, err := strconv.Atoi(tok)
			if err != nil {
				return fmt.Errorf("test %d: invalid integer %q", idx+1, tok)
			}
			if v < 1 || v > len(tc.perm) {
				return fmt.Errorf("test %d: index %d out of range", idx+1, v)
			}
			if i > 0 && v <= indices[i-1] {
				return fmt.Errorf("test %d: indices must be strictly increasing", idx+1)
			}
			indices[i] = v
		}
		values := make([]int, tc.n+1)
		for i, pos := range indices {
			values[i] = tc.perm[pos-1]
		}
		if !isMonotone(values) {
			return fmt.Errorf("test %d: sequence is not monotone", idx+1)
		}
	}
	return nil
}

func readNonEmptyLines(out string) []string {
	scanner := bufio.NewScanner(strings.NewReader(out))
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func isMonotone(vals []int) bool {
	inc := true
	for i := 1; i < len(vals); i++ {
		if vals[i] <= vals[i-1] {
			inc = false
			break
		}
	}
	if inc {
		return true
	}
	dec := true
	for i := 1; i < len(vals); i++ {
		if vals[i] >= vals[i-1] {
			dec = false
			break
		}
	}
	return dec
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref2149C-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2149C.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}
