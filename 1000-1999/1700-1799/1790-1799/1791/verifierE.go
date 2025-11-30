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
	input    string
	expected string
}

const testcaseData = `100
4 -1 -7 2 5
3 -8 -8 -10
7 7 -1 -9 -3 6 7 1
5 -5 -7 -2 -4 -10
5 -2 -4 -5 -1 -1
6 -8 9 0 2 6 -3
3 -3 5 -2
2 7 -1
1 -1
10 -1 6 -4 3 3 9 -1 3 4 -5
4 -1 -2 -9 -8
1 4
5 6 7 10 5 0
3 -4 -8 3
4 10 10 4 -2
3 1 3 8
6 10 7 -4 0 -7 -9
4 -2 8 9 -3
2 0 -5
5 4 -10 -9 1 -8
5 0 -10 0 -1 0
3 10 3 9
2 -1 9
4 4 -1 -6 -2
7 9 -5 0 8 -10 1 -9
8 -5 1 1 -1 8 -7 4 -4
7 -4 -7 -9 -9 -9 -5 9
3 9 -9 7
8 8 -3 0 -9 -7 6 -1 3
4 5 -4 -3 4
7 5 -9 -3 3 4 -3 10
7 -4 5 -4 -9 -9 -2 -2
4 6 -4 -3 3
5 -6 0 -9 0 8
2 8 2
1 5
7 -8 3 -4 8 -5 0 -1
8 10 0 3 6 -4 10 -2 0
7 5 -8 -2 10 -4 -9 2
10 -6 -2 -9 -5 10 4 8 5 2 2
4 -10 -4 -5 -10
10 -2 -7 2 2 -3 7 -9 -4 -5 9
6 7 5 6 4 -10 -8
1 9
2 5 7
5 9 -6 -9 1 -8
9 -10 -1 1 -8 -8 7 4 2 -4
5 2 -3 5 2 -7
2 -7 9
6 6 3 3 4 -8 10
4 10 -1 5 3
2 7 -5
6 -5 -5 -6 0 5 0
5 7 -10 -5 -10 10
5 -7 7 -7 5 9
8 6 -8 6 -3 3 -1 1 -3
3 10 -10 -9
10 0 7 4 8 -1 6 4 9 9 4
7 -6 -2 9 1 0 -6 3
2 9 -6
10 -5 -1 1 -4 8 1 9 -8 -8 2
3 0 10 1
6 -5 -1 -10 9 -10 6
2 1 -7
3 -5 8 5
10 -8 -7 -5 10 5 -3 9 -1 2 9
4 5 -3 -1 1
4 0 7 10 6
8 2 6 2 0 -1 4 3 8
1 -2
3 7 4 7
10 1 2 2 9 -10 -6 6 -8 4 10
6 9 -1 -8 -2 5 -3
8 9 -8 -6 -3 -8 -1 -6 -9
3 2 8 9
9 7 -2 -10 6 -3 -5 -7 -4 -9
6 -8 -7 -2 -9 10 -1
10 10 -5 -6 10 3 -6 -8 7 1 -10
9 -6 6 3 -6 -4 -1 5 6 -8
7 -5 -5 -2 6 2 7 -1
7 0 -5 2 -9 3 -10 -2
1 -1
3 -8 -5 -7
10 -10 -3 -3 7 -10 5 7 -5 4 2
6 -5 6 8 -4 -7 5
10 1 -1 10 3 9 -9 10 -4 -2 8
10 -1 5 10 0 -8 -3 0 -7 -9 9
6 6 5 1 -8 -5 -9
8 6 7 9 -3 -9 -4 -8 0
10 -6 -1 -7 6 6 -10 -8 -3 9 -2
10 -10 -9 -10 5 -6 -4 1 -3 1 -1
7 10 9 2 -9 -5 3 5
2 6 6
5 7 -1 -4 3 0
3 -4 -10 5
7 -2 0 3 10 3 0 9
4 -2 -2 6 -8
1 2
5 9 -1 10 -6 9
6 -8 -5 9 -8 -4 -2`

func solve(arr []int64) int64 {
	neg := 0
	hasZero := false
	var sum int64
	minAbs := int64(1<<63 - 1)
	for _, v := range arr {
		if v < 0 {
			neg++
		}
		if v == 0 {
			hasZero = true
		}
		abs := v
		if abs < 0 {
			abs = -abs
		}
		sum += abs
		if abs < minAbs {
			minAbs = abs
		}
	}
	if neg%2 == 1 && !hasZero {
		sum -= 2 * minAbs
	}
	return sum
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing value", caseNum+1)
			}
			v, err := strconv.ParseInt(fields[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", caseNum+1, err)
			}
			arr[i] = v
			pos++
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')

		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(arr), 10),
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
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(bin, tc.input)
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
