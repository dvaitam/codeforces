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

const refSource = "2000-2999/2000-2099/2010-2019/2011/2011B.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierB.go /path/to/candidate")
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInputData(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refPossible, err := parseReferenceOutput(refOut, tests)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	if err := validateCandidateOutput(candOut, tests, refPossible); err != nil {
		fail("%v", err)
	}

	fmt.Println("OK")
}

func parseInputData(data []byte) ([]int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	if t <= 0 {
		return nil, fmt.Errorf("t must be positive")
	}
	tests := make([]int, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i]); err != nil {
			return nil, err
		}
	}
	return tests, nil
}

func parseReferenceOutput(out string, tests []int) ([]bool, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	possible := make([]bool, len(tests))
	for i, n := range tests {
		token, err := readToken(reader)
		if err != nil {
			return nil, fmt.Errorf("reference output ended early at test %d: %v", i+1, err)
		}
		if token == "-1" {
			possible[i] = false
			continue
		}
		perm, err := parsePermutation(reader, n, token)
		if err != nil {
			return nil, fmt.Errorf("reference test %d: %v", i+1, err)
		}
		if err := validatePermutation(n, perm); err != nil {
			return nil, fmt.Errorf("reference test %d: %v", i+1, err)
		}
		possible[i] = true
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("reference output has extra token %q", extra)
		}
		return nil, err
	}
	return possible, nil
}

func validateCandidateOutput(out string, tests []int, refPossible []bool) error {
	reader := bufio.NewReader(strings.NewReader(out))
	for i, n := range tests {
		token, err := readToken(reader)
		if err != nil {
			return fmt.Errorf("candidate output ended early at test %d: %v", i+1, err)
		}
		if token == "-1" {
			if refPossible[i] {
				return fmt.Errorf("test %d admits a solution but candidate printed -1", i+1)
			}
			continue
		}
		perm, err := parsePermutation(reader, n, token)
		if err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
		if err := validatePermutation(n, perm); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return fmt.Errorf("candidate output has extra token %q", extra)
		}
		return err
	}
	return nil
}

func parsePermutation(reader *bufio.Reader, n int, firstToken string) ([]int, error) {
	perm := make([]int, n)
	val, err := strconv.Atoi(firstToken)
	if err != nil {
		return nil, fmt.Errorf("invalid integer %q", firstToken)
	}
	perm[0] = val
	for idx := 1; idx < n; idx++ {
		tok, err := readToken(reader)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d numbers, got %d", n, idx)
			}
			return nil, err
		}
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		perm[idx] = val
	}
	return perm, nil
}

func validatePermutation(n int, perm []int) error {
	if len(perm) != n {
		return fmt.Errorf("permutation length mismatch, expected %d got %d", n, len(perm))
	}
	seen := make([]bool, n+1)
	for i, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d at position %d out of range [1,%d]", v, i+1, n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
	}
	for i := 0; i+1 < n; i++ {
		a := perm[i]
		b := perm[i+1]
		if a%b == 0 || b%a == 0 {
			if !(a < b) {
				return fmt.Errorf("positions %d and %d must be increasing because one divides the other", i+1, i+2)
			}
		} else {
			if !(a > b) {
				return fmt.Errorf("positions %d and %d must be decreasing because neither divides the other", i+1, i+2)
			}
		}
	}
	return nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2011B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
}

func readToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}
		if !isSpace(b) {
			sb.WriteByte(b)
			break
		}
	}
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if isSpace(b) {
			break
		}
		sb.WriteByte(b)
	}
	return sb.String(), nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t' || b == '\v' || b == '\f'
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
