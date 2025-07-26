package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input  string
	output string
}

func expected(n int, s string) int {
	maxDelay := 0
	cur := 0
	seenA := false
	for i := 0; i < len(s); i++ {
		if s[i] == 'A' {
			seenA = true
			cur = 0
		} else if seenA {
			cur++
			if cur > maxDelay {
				maxDelay = cur
			}
		}
	}
	return maxDelay
}

func buildTests() []testCase {
	r := rand.New(rand.NewSource(1))
	var tests []testCase
	// Add some fixed edge cases
	edge := []string{"A", "P", "AA", "PP", "AP", "PA", "APA", "PAP", strings.Repeat("A", 100), strings.Repeat("P", 100)}
	for _, s := range edge {
		input := fmt.Sprintf("1\n%d %s\n", len(s), s)
		output := fmt.Sprintf("%d\n", expected(len(s), s))
		tests = append(tests, testCase{input, output})
	}
	for len(tests) < 100 {
		n := r.Intn(100) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			if r.Intn(2) == 0 {
				b[i] = 'A'
			} else {
				b[i] = 'P'
			}
		}
		s := string(b)
		input := fmt.Sprintf("1\n%d %s\n", n, s)
		output := fmt.Sprintf("%d\n", expected(n, s))
		tests = append(tests, testCase{input, output})
	}
	return tests
}

func run(binary string, in string) (string, string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(in)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		return "", errBuf.String(), err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return outBuf.String(), errBuf.String(), err
	case <-time.After(2 * time.Second):
		cmd.Process.Kill()
		return outBuf.String(), errBuf.String(), fmt.Errorf("timeout")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := buildTests()
	for i, tc := range tests {
		out, errStr, err := run(binary, tc.input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", i+1, err, errStr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(tc.output) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, tc.output, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
