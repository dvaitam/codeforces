package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type testCase struct {
	a, b, c int
}

func buildReferenceBinary() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine current file path")
	}
	dir := filepath.Dir(file)
	refPath := filepath.Join(dir, "2087A_ref.bin")
	cmd := exec.Command("go", "build", "-o", refPath, "2087A.go")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.Remove(refPath)
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return refPath, nil
}

func parseInput(input string) ([]testCase, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &cases[i].a, &cases[i].b, &cases[i].c); err != nil {
			return nil, fmt.Errorf("failed to read test case %d: %v", i+1, err)
		}
	}
	return cases, nil
}

func readOutputLines(out string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimRight(scanner.Text(), "\r"))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}

func validatePassword(tc testCase, pw string) error {
	expectedLen := tc.a + tc.b + tc.c
	if len(pw) != expectedLen {
		return fmt.Errorf("expected length %d, got %d", expectedLen, len(pw))
	}

	var digits, upper, lower int
	for i := 0; i < len(pw); i++ {
		ch := pw[i]
		if i > 0 && ch == pw[i-1] {
			return fmt.Errorf("adjacent characters at positions %d and %d are both %q", i, i+1, ch)
		}
		switch {
		case ch >= '0' && ch <= '9':
			digits++
		case ch >= 'A' && ch <= 'Z':
			upper++
		case ch >= 'a' && ch <= 'z':
			lower++
		default:
			return fmt.Errorf("invalid character %q at position %d", ch, i+1)
		}
	}
	if digits != tc.a || upper != tc.b || lower != tc.c {
		return fmt.Errorf("expected counts (digits=%d, upper=%d, lower=%d) but got (%d, %d, %d)",
			tc.a, tc.b, tc.c, digits, upper, lower)
	}
	return nil
}

func validateOutput(cases []testCase, out string) error {
	lines, err := readOutputLines(out)
	if err != nil {
		return fmt.Errorf("failed to read output: %v", err)
	}
	if len(lines) != len(cases) {
		return fmt.Errorf("expected %d lines, got %d", len(cases), len(lines))
	}
	for i, tc := range cases {
		if err := validatePassword(tc, lines[i]); err != nil {
			return fmt.Errorf("test case %d: %v", i+1, err)
		}
	}
	return nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(stderr.String()))
		}
		return "", err
	}
	return stdout.String(), nil
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.a))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.b))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.c))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	var tests []string

	tests = append(tests, buildInput([]testCase{{1, 1, 1}}))
	tests = append(tests, buildInput([]testCase{{10, 10, 10}}))
	tests = append(tests, buildInput([]testCase{{1, 10, 1}, {10, 1, 1}, {1, 1, 10}, {5, 6, 7}}))
	tests = append(tests, buildInput([]testCase{{2, 1, 6}, {3, 3, 6}}))

	for i := 0; i < 120; i++ {
		t := rng.Intn(20) + 1
		cases := make([]testCase, t)
		for j := 0; j < t; j++ {
			cases[j] = testCase{
				a: rng.Intn(10) + 1,
				b: rng.Intn(10) + 1,
				c: rng.Intn(10) + 1,
			}
		}
		tests = append(tests, buildInput(cases))
	}

	for i := 0; i < 20; i++ {
		t := rng.Intn(200) + 50
		cases := make([]testCase, t)
		for j := 0; j < t; j++ {
			cases[j] = testCase{
				a: 10 - (j % 10),
				b: (j % 10) + 1,
				c: rng.Intn(10) + 1,
			}
		}
		tests = append(tests, buildInput(cases))
	}

	large := make([]testCase, 1000)
	for i := range large {
		large[i] = testCase{
			a: rng.Intn(10) + 1,
			b: rng.Intn(10) + 1,
			c: rng.Intn(10) + 1,
		}
	}
	tests = append(tests, buildInput(large))

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refPath, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, input := range tests {
		cases, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := validateOutput(cases, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := validateOutput(cases, userOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, userOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
