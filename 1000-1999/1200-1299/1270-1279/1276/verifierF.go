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
}

func solveRef(input string) (int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return 0, err
	}
	n := len(s)
	uniq := make(map[string]struct{})
	uniq[""] = struct{}{}
	addSubs := func(str string) {
		for i := 0; i < len(str); i++ {
			for j := i + 1; j <= len(str); j++ {
				uniq[str[i:j]] = struct{}{}
			}
		}
	}
	addSubs(s)
	for i := 0; i < n; i++ {
		bytes := []byte(s)
		bytes[i] = '*'
		addSubs(string(bytes))
	}
	return len(uniq), nil
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

func parseOutput(out string) (int, error) {
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	var val int
	if _, err := fmt.Sscan(out, &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val, nil
}

func makeCase(name, s string) testCase {
	return testCase{name: name, input: fmt.Sprintf("%s\n", s)}
}

func randomString(rng *rand.Rand, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(3))) // limit alphabet for variation
	}
	return sb.String()
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_a", "a"),
		makeCase("single_b", "b"),
		makeCase("double", "ab"),
		makeCase("repeat", "aaaa"),
		makeCase("mixed", "abcab"),
		makeCase("pal", "abba"),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 80; i++ {
		n := rng.Intn(12) + 1
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", i+1), randomString(rng, n)))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%d\nactual:%d\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
