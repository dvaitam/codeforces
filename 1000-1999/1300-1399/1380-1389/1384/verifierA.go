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

type caseA struct {
	n   int
	arr []int
}

func generateTests() []caseA {
	r := rand.New(rand.NewSource(1))
	tests := []caseA{
		{1, []int{0}},
		{1, []int{1}},
		{5, []int{0, 1, 2, 3, 4}},
	}
	for len(tests) < 120 {
		n := r.Intn(10) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = r.Intn(51)
		}
		tests = append(tests, caseA{n, arr})
	}
	return tests
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func lcp(a, b string) int {
	m := len(a)
	if len(b) < m {
		m = len(b)
	}
	for i := 0; i < m; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return m
}

func verify(tc caseA, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != tc.n+1 {
		return fmt.Errorf("expected %d lines, got %d", tc.n+1, len(lines))
	}
	for i, line := range lines {
		s := strings.TrimSpace(line)
		if len(s) == 0 {
			return fmt.Errorf("line %d empty", i+1)
		}
		if len(s) > 200 {
			return fmt.Errorf("line %d too long", i+1)
		}
		for _, ch := range s {
			if ch < 'a' || ch > 'z' {
				return fmt.Errorf("invalid char in line %d", i+1)
			}
		}
		lines[i] = s
	}
	for i := 0; i < tc.n; i++ {
		got := lcp(lines[i], lines[i+1])
		if got != tc.arr[i] {
			return fmt.Errorf("expected lcp %d at pair %d, got %d", tc.arr[i], i+1, got)
		}
	}
	return nil
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
