package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n       int
	m       int
	candies []int
}

func parseTestcases(path string) ([]testCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err1 := strconv.Atoi(parts[0])
		m, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: invalid n or m", idx+1)
		}
		if len(parts) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d candies, got %d", idx+1, n, len(parts)-2)
		}
		candies := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(parts[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid candy value", idx+1)
			}
			candies[i] = val
		}
		tests = append(tests, testCase{n: n, m: m, candies: candies})
	}
	return tests, nil
}

func expectedChild(tc testCase) int {
	maxRounds := -1
	ans := 1
	for i, val := range tc.candies {
		rounds := (val + tc.m - 1) / tc.m
		if rounds >= maxRounds {
			maxRounds = rounds
			ans = i + 1
		}
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatInput(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.candies {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		input := formatInput(tc)
		expected := fmt.Sprintf("%d", expectedChild(tc))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
