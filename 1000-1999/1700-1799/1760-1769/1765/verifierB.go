package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
	t     int
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", err
	}
	var outputs []string
	for i := 0; i < t; i++ {
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		pos := 0
		double := false
		ok := true
		for pos < n {
			if double {
				if pos+1 >= n || s[pos] != s[pos+1] {
					ok = false
					break
				}
				pos += 2
			} else {
				pos++
			}
			double = !double
		}
		if pos != n {
			ok = false
		}
		if ok {
			outputs = append(outputs, "YES")
		} else {
			outputs = append(outputs, "NO")
		}
	}
	return strings.Join(outputs, "\n"), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	return lines, nil
}

func makeCase(name string, tests []string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, s := range tests {
		fmt.Fprintf(&sb, "%d\n%s\n", len(s), s)
	}
	return testCase{name: name, input: sb.String(), t: len(tests)}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for idx := 0; idx < 40; idx++ {
		count := rng.Intn(10) + 1
		strs := make([]string, count)
		for i := 0; i < count; i++ {
			n := rng.Intn(8) + 1
			var sb strings.Builder
			for j := 0; j < n; j++ {
				sb.WriteByte(byte('a' + rng.Intn(3)))
			}
			strs[i] = sb.String()
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", idx+1), strs))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_yes", []string{"a"}),
		makeCase("double", []string{"aa"}),
		makeCase("sample", []string{"ossu", "aa", "abc", "qwrr"}),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expLines := strings.Fields(expect)
		gotLines, err := parseOutput(out, len(expLines))
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		for i := range expLines {
			if expLines[i] != gotLines[i] {
				fmt.Printf("test %d (%s) mismatch in case %d\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, i+1, tc.input, expect, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
