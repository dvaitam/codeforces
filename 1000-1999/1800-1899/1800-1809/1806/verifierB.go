package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `2 0 3
4 3 3 1 3
2 2 1
4 1 3 1 2
2 0 0
5 1 0 3 3 2
3 0 1 1
2 1 3
5 2 1 0 3 1
5 3 3 2 3 3
3 1 0 2
4 2 0 1 2
3 3 2 3
2 0 0
2 1 1
2 2 0
5 2 1 1 3 2
5 0 0 0 3 1
4 3 3 3 1
2 0 1
4 2 2 0 3
4 2 3 3 3
2 1 1
3 2 2 0
2 3 3
3 3 0 1
4 3 0 3 3
5 0 0 3 1 0
2 1 2
2 2 1
5 3 0 0 3 2
2 2 1
5 2 0 1 0 3
5 2 1 3 2 0
2 0 3
4 0 1 1 0
5 2 0 1 3 3
2 0 3
2 0 3
3 0 3 3
5 0 3 2 2 0
4 0 0 2 0
4 2 1 0 1
4 0 1 1 0
3 0 0 2
4 0 1 1 1
5 0 3 2 2 1
2 1 2
4 3 2 2 2
3 0 0 2
3 3 1 1
3 2 1 1
3 2 2 3
2 1 0
5 0 0 1 3 2
5 1 3 2 2 0
3 3 2 0
5 3 0 3 2 3
3 2 2 3
2 1 0
4 0 1 3 3
3 3 3 1
3 3 2 1
4 3 0 3 1
4 0 3 3 0
3 2 0 2
3 3 3 0
5 1 3 2 0 0
4 0 0 2 3
5 3 0 2 2 2
3 0 0 0
4 2 2 0 1
3 0 3 2
4 1 1 0 3
5 2 2 3 2 1
3 0 0 3
5 3 1 2 2 3
5 1 3 3 2 3
2 3 2
3 3 0 1
2 2 3
5 0 0 0 3 0
4 0 0 0 2
4 1 1 2 1
2 3 3
4 3 1 2 3
5 1 3 1 2 1
3 1 3 2
5 3 3 3 1 1
5 1 0 3 0 1
2 1 2
2 1 1
4 0 2 2 3
5 0 3 3 3 2
5 1 2 2 0 0
2 1 2
2 2 0
3 0 3 1
5 1 3 1 2 0
2 2 2`

type testCase struct {
	input    string
	expected string
}

func solve(arr []int) int {
	n := len(arr)
	c0, c1 := 0, 0
	for _, v := range arr {
		if v == 0 {
			c0++
		} else if v == 1 {
			c1++
		}
	}
	if c0 <= (n+1)/2 {
		return 0
	}
	others := n - c0 - c1
	if c1 == 0 || others > 0 {
		return 1
	}
	return 2
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %w", idx+1, err)
		}
		if len(parts) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d values got %d", idx+1, n+1, len(parts))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value: %w", idx+1, err)
			}
			arr[i] = v
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.Itoa(solve(arr)),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
