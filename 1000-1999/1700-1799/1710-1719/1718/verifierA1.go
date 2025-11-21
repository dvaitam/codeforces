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
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return "", err
		}
		for j := 0; j < n; j++ {
			var x int
			fmt.Fscan(reader, &x)
		}
	}
	// reference prints 0 for each test case
	lines := make([]string, t)
	for i := 0; i < t; i++ {
		lines[i] = "0"
	}
	return strings.Join(lines, "\n"), nil
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

func parseOutput(out string, expectedLines int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expectedLines {
		return nil, fmt.Errorf("expected %d outputs, got %d", expectedLines, len(lines))
	}
	return lines, nil
}

func makeCase(name string, tests []struct {
	n int
	a []int
}) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, val := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String(), t: len(tests)}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for idx := 0; idx < 40; idx++ {
		tcCount := rng.Intn(5) + 1
		cases := make([]struct {
			n int
			a []int
		}, tcCount)
		for i := 0; i < tcCount; i++ {
			n := rng.Intn(5) + 1
			cases[i].n = n
			cases[i].a = make([]int, n)
			for j := 0; j < n; j++ {
				cases[i].a[j] = rng.Intn(10)
			}
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", idx+1), cases))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_zero", []struct {
			n int
			a []int
		}{
			{n: 1, a: []int{0}},
		}),
		makeCase("two_tests", []struct {
			n int
			a []int
		}{
			{n: 2, a: []int{1, 2}},
			{n: 3, a: []int{3, 3, 3}},
		}),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
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
		expectLines := strings.Fields(expect)
		gotLines, err := parseOutput(out, len(expectLines))
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		for i := range expectLines {
			if expectLines[i] != gotLines[i] {
				fmt.Printf("test %d (%s) mismatch at line %d\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, i+1, tc.input, expect, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
