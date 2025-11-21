package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2160-2169/2162/2162B.go"

type testCase struct {
	name  string
	input string
	t     int
}

func runProgram(bin, input string) (string, error) {
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

func parseOutput(out string, expected int) ([][]int, error) {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	result := make([][]int, 0, expected)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, nil)
			fields := strings.Fields(line)
			for _, f := range fields {
				val, err := strconv.Atoi(f)
				if err != nil {
					return nil, fmt.Errorf("invalid integer %q", f)
				}
				result[len(result)-1] = append(result[len(result)-1], val)
			}
		}
	}
	if len(result)%2 != 0 {
		return nil, fmt.Errorf("expected pairs of lines per test case, got %d lines", len(result))
	}
	if len(result)/2 != expected {
		return nil, fmt.Errorf("expected %d cases, got %d", expected, len(result)/2)
	}
	values := make([][]int, expected*2)
	copy(values, result)
	return values, nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_char", input: "1\n1 a\n", t: 1},
		{name: "small_palindrome", input: "1\n4 abba\n", t: 1},
		{name: "multiple_cases", input: "2\n3 abc\n4 abca\n", t: 2},
	}
}

func randomTests(count int) []testCase {
	tests := make([]testCase, 0, count)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for j := 0; j < t; j++ {
			n := rng.Intn(7) + 1
			var s strings.Builder
			for k := 0; k < n; k++ {
				s.WriteByte(byte('a' + rng.Intn(3)))
			}
			sb.WriteString(fmt.Sprintf("%d %s\n", n, s.String()))
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: sb.String(),
			t:     t,
		})
	}
	return tests
}

func compareResults(refLines, candLines [][]int) error {
	if len(refLines) != len(candLines) {
		return fmt.Errorf("line count mismatch")
	}
	for i := range refLines {
		if len(refLines[i]) != len(candLines[i]) {
			return fmt.Errorf("line %d length mismatch: expected %d got %d", i+1, len(refLines[i]), len(candLines[i]))
		}
		for j := range refLines[i] {
			if refLines[i][j] != candLines[i][j] {
				return fmt.Errorf("line %d value mismatch at position %d: expected %d got %d",
					i+1, j+1, refLines[i][j], candLines[i][j])
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refBin, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	tests := append(manualTests(), randomTests(50)...)
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refLines, err := parseOutput(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candLines, err := parseOutput(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if err := compareResults(refLines, candLines); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) mismatch: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
