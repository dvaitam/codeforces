package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	s string
}

func parseTestcases(path string) ([]testCase, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cases []testCase
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		s := fields[1]
		if len(s) != n {
			return nil, fmt.Errorf("len mismatch")
		}
		cases = append(cases, testCase{n: n, s: s})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func cmpStrings(a, b string) bool {
	i := 0
	for i < len(a) && a[i] == '0' {
		i++
	}
	a = a[i:]
	j := 0
	for j < len(b) && b[j] == '0' {
		j++
	}
	b = b[j:]
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}

func solveCase(tc testCase) string {
	n := tc.n
	s := tc.s
	best := strings.Repeat("9", n)
	digits := make([]int, n)
	for i := 0; i < n; i++ {
		digits[i] = int(s[i] - '0')
	}
	for rot := 0; rot < n; rot++ {
		rotDigits := make([]int, n)
		for i := 0; i < n; i++ {
			rotDigits[i] = digits[(i-rot+n)%n]
		}
		for add := 0; add < 10; add++ {
			b := make([]byte, n)
			for i := 0; i < n; i++ {
				b[i] = byte((rotDigits[i]+add)%10) + '0'
			}
			cand := string(b)
			if cmpStrings(cand, best) {
				best = cand
			}
		}
	}
	return best
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		expected := solveCase(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
