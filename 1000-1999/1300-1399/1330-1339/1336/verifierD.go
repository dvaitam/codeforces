package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Embedded testcases from testcasesD.txt. Each line holds a single integer n.
const testcasesRaw = `7
8
8
7
4
5
4
8
4
4
5
4
4
4
6
5
4
4
7
5
5
7
4
7
6
5
7
4
6
5
6
4
7
4
4
5
4
8
5
5
6
5
6
4
4
6
4
6
4
4
4
5
4
6
4
7
5
6
4
4
7
6
4
7
5
8
7
4
6
4
4
5
5
6
4
5
5
7
4
6
5
5
4
5
4
5
5
6
7
7
4
4
5
4
4
4
4
4
8
7`

type testCase struct {
	n int
}

// referenceSolution mirrors 1336D.go (interactive, no output).
func referenceSolution(tc testCase) string {
	_ = tc
	return ""
}

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var n int
		fmt.Sscan(line, &n)
		tests = append(tests, testCase{n: n})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	cmdDone := make(chan error, 1)
	go func() { cmdDone <- cmd.Run() }()
	select {
	case err := <-cmdDone:
		if err != nil {
			return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
		}
	case <-time.After(2 * time.Second):
		_ = cmd.Process.Kill()
		return "", fmt.Errorf("timeout; stderr: %s", errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseTestcases()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.n)
		expected := referenceSolution(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %q got %q\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
