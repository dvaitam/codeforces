package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	name  string
	input string
}

func solveRef(input string) (int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, err
	}
	votes := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &votes[i]); err != nil {
			return 0, err
		}
	}
	count := 0
	for {
		maxIdx := 0
		for i := 1; i < n; i++ {
			if votes[i] > votes[maxIdx] {
				maxIdx = i
			}
		}
		if maxIdx == 0 {
			break
		}
		votes[0]++
		votes[maxIdx]--
		count++
	}
	return count, nil
}

func makeCase(name string, votes []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(votes)))
	for i, v := range votes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("simple_win", []int{5, 3, 4}),
		makeCase("already_winner", []int{10, 3, 3, 3}),
		makeCase("tie_many", []int{5, 5, 5}),
		makeCase("large_gap", []int{1, 100}),
		makeCase("all_equal", []int{7, 7, 7, 7, 7}),
		makeCase("two_candidates", []int{1, 1}),
		makeCase("descending", []int{2, 5, 4, 3, 1}),
		makeCase("mixed", []int{3, 8, 2, 7, 4}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(574))
	var tests []testCase
	for i := 0; i < 150; i++ {
		n := rng.Intn(10) + 2
		if i > 100 {
			n = rng.Intn(100-2+1) + 2
		}
		votes := make([]int, n)
		for j := 0; j < n; j++ {
			votes[j] = rng.Intn(1000) + 1
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", i+1), votes))
	}
	return tests
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseSingleInt(output string) (int, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %v", fields)
	}
	var val int
	if _, err := fmt.Sscan(fields[0], &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		out, runErr := runCandidate(bin, tc.input)
		if runErr != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, runErr, tc.input)
			os.Exit(1)
		}
		val, parseErr := parseSingleInt(out)
		if parseErr != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, parseErr, tc.input, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Printf("test %d (%s) failed: expect %d got %d\ninput:\n%s\n", idx+1, tc.name, expect, val, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
