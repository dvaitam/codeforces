package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseB struct {
	n int
	s string
}

func generateTestsB() []testCaseB {
	r := rand.New(rand.NewSource(1))
	tests := []testCaseB{{3, "010"}, {5, "11111"}, {4, "0000"}}
	for len(tests) < 120 {
		n := r.Intn(20) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests = append(tests, testCaseB{n, sb.String()})
	}
	return tests
}

func expectedB(tc testCaseB) string {
	zeros := strings.Count(tc.s, "0")
	if zeros == 0 || zeros == tc.n {
		return "0"
	}
	target := strings.Repeat("0", zeros) + strings.Repeat("1", tc.n-zeros)
	if tc.s == target {
		return "0"
	}
	indices := make([]int, 0)
	for i := 0; i < tc.n; i++ {
		if tc.s[i] != target[i] {
			indices = append(indices, i+1)
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(len(indices)))
	for _, idx := range indices {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(idx))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verifyOutputB(tc testCaseB, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	exp := expectedB(tc)
	expLines := strings.Split(exp, "\n")
	if len(lines) != len(expLines) {
		return fmt.Errorf("expected %d lines got %d", len(expLines), len(lines))
	}
	for i := range lines {
		if strings.TrimSpace(lines[i]) != strings.TrimSpace(expLines[i]) {
			return fmt.Errorf("line %d mismatch: expected %q got %q", i+1, expLines[i], lines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsB()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verifyOutputB(tc, got); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
