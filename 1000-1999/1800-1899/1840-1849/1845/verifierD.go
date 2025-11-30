package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `2 4 1
7 -2 -2 -4 4 1 -4 -1
1 1
5 2 1 4 -4 -2
6 -5 5 5 -4 4 2
6 -3 -1 2 -3 2 -1
1 2
5 -1 5 5 -1 3
4 2 4 -2 5
7 -2 2 3 -4 2 -5 3
4 -2 5 -4 5
8 1 1 4 2 -4 1 -1 1
6 3 -4 -1 5 -3 -1
8 4 -2 3 1 4 4 2 -3
1 4
8 5 1 -2 -1 -5 -1 -5 2
8 -5 2 -2 2 -2 1 -4 -4
3 -5 -5 -5
8 5 -1 3 2 -5 -5 -3 -2
4 4 -4 5 1
5 -4 1 5 5 -3
3 -2 2 -1
6 -4 -4 -1 -2 5 3
5 -2 1 3 -3 1
3 -3 -4 5
8 4 3 1 -4 -2 4 -1 4
1 1
6 1 -5 2 -4 -1 3
1 -1
8 -1 -2 -3 -3 5 5 5 2
2 2 -5
4 -2 -4 4 -4
3 -1 4 -4
1 5
3 1 3 3
3 4 1 -1
8 -2 2 4 2 1 -4 -2 4
1 1
8 3 -2 4 4 3 1 1 4
2 -1 -2
1 5
2 2 -4
3 1 -2 2
5 -1 -5 5 2 1
7 5 4 -5 2 -5 1 2
1 1
5 5 4 1 2 -4
1 -5
8 1 -4 4 5 1 -2 1 -3
8 1 -3 -3 -4 1 -4 5 5
8 -5 -5 3 -2 -3 1 1 1
7 -1 -1 -3 3 3 -5 5
2 1 -4
6 2 -5 2 4 4 2
2 2 1
5 1 -1 2 1 -1
2 -1 4
7 -3 -1 -5 -5 1 -5 -4
5 2 5 1 -4 5
8 4 5 -3 -3 -2 2 3 1
5 1 1 5 -4 2
8 -1 5 -2 3 2 2 4 -4
7 -5 -5 2 -4 -1 -2 1
8 1 -2 1 3 1 -1 -3 -5
2 2 1
5 5 -3 -3 -1 2
6 4 -1 -1 5 1 -2
1 3
8 3 1 4 -1 -2 -3 -3 4
4 3 -4 1 -3
4 -4 -5 1 5
7 -3 -4 4 -3 5 -2 3
7 5 2 5 1 -3 4 3
6 2 1 3 1 -4 3
3 2 -3 -1
6 5 1 5 2 -1 2
6 -1 -3 1 -4 5 -1
8 -2 -5 -5 1 1 -2 -3 5
3 1 -1 -5
3 4 5 5
8 2 4 1 -5 1 3 -3 5
6 5 -2 1 3 1 -1
2 3 4
3 -2 -2 -2
6 -1 4 -3 1 -3 -2
5 -2 -5 3 3 4
7 -4 -2 2 2 -3 -2 1
5 -1 3 -2 1 -4
1 -5
6 -2 5 1 5 -3 -5
3 -3 1 5
5 3 1 1 3 3
5 1 -3 1 4 -5
3 5 4 1
2 2 1
2 2 4
1 1
8 4 3 5 1 4 -5 -1 -1
8 -5 -5 2 -3 4 1 3 2
4 -4 -3 2 -3
6 -2 -5 1 -4 3 -3
5 -3 -2 -5 -1 -1
8 5 -5 -1 1 -4 2 -4 -5
7 2 -2 -2 2 1 -1 4
7 -2 2 1 -4 -2 3 -1
7 5 -3 5 3 2 2 -1
8 1 1 1 1 -2 -4 -4 5
7 2 1 2 -4 5 5 -3
6 1 -5 1 4 4 -4
5 3 3 1 -2 -2
6 -2 1 2 -2 5 -2
4 -4 1 5 -3
3 1 1 1
3 4 5 -2
2 1 -3
8 -5 3 1 -1 2 2 -4 1
3 -2 2 5
3 2 -1 -3
3 -2 -2 -4
4 -2 -3 2 -1`

type testCase struct {
	input    string
	expected string
}

func solve(arr []int64) int64 {
	const INF int64 = 1_000_000_000_000_000_000
	n := len(arr)
	sum := make([]int64, n+1)
	for i := 0; i < n; i++ {
		sum[i+1] = sum[i] + arr[i]
	}
	mins := make([]int64, n+1)
	copy(mins, sum)
	for i := n - 1; i >= 0; i-- {
		if mins[i+1] < mins[i] {
			mins[i] = mins[i+1]
		}
	}
	maxv := sum[n]
	ans := INF
	for i := 0; i < n; i++ {
		diff := sum[i] - mins[i]
		if diff < 0 {
			diff = 0
		}
		v := sum[n] + diff
		if v > maxv {
			maxv = v
			ans = sum[i]
		}
	}
	return ans
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
		if len(parts) < 1 {
			continue
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", idx+1, err)
		}
		if len(parts) != n+1 {
			return nil, fmt.Errorf("case %d: expected %d numbers got %d", idx+1, n+1, len(parts))
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(parts[1+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", idx+1, err)
			}
			arr[i] = val
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
		fmt.Println("usage: verifierD /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
