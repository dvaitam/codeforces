package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcaseData = `
100
5
2 -2 3 2 0 -2 -2 0 -2 0
4
-1 -5 1 -1 -5 -5
2
-3
2
-5
5
-4 -4 -4 -4 2 -4 0 -4 0 -4
3
-1 -1 -1
4
-4 -4 -2 -1 -4 -2
3
0 2 0
4
-1 -2 -1 2 -2 -2
2
-2
3
-2 -4 -4
4
-1 -5 -3 -5 -3 -5
2
-4
2
-5
3
-2 -2 0
4
-4 -5 -5 -4 -5 -4
2
-5
2
-3
5
1 1 -5 -5 -5 1 -5 1 1 2
2
4
3
-1 -1 2
2
-4
5
-2 -3 -2 1 -3 1 3 -3 -2 -3
4
0 0 1 3 0 1
2
-5
5
-2 -2 -3 -3 0 0 -2 -3 -3 3
4
-3 -3 -3 -3 -2 -3
3
-2 2 -2
3
-4 -5 -5
3
2 -5 -5
3
-3 -3 -3
2
-3
3
0 -4 -4
5
-4 0 -4 0 0 -4 -1 -4 -1 -1
3
4 -3 -3
4
-3 -5 -5 -3 -3 -5
3
-4 4 -4
4
-2 -2 0 -4 -4 -4
4
-4 -4 -4 -2 -2 -2
5
1 1 1 1 -3 1 -3 -3 1 -3
2
-5
2
-1
2
-4
3
-4 -3 -4
4
-5 -5 4 -5 2 2
2
-2
3
1 3 1
4
-5 -2 -4 -5 -5 -4
4
-5 -5 -5 -4 -1 -4
2
-3
4
-4 -4 -4 -1 -4 -4
2
3
5
-3 -3 0 -3 -3 -3 -3 0 -3 0
3
-4 -4 -4
5
-5 -5 -2 -5 1 3 -2 -5 -2 1
2
3
3
-4 -3 -4
3
-5 -3 -5
3
-3 -3 4
2
-4
5
0 -1 3 -1 -1 2 2 0 -1 0
2
-4
4
1 -1 0 -1 -1 0
4
-3 1 1 -3 -3 4
3
-1 -5 -5
4
-1 -1 -1 -1 -1 -1
5
-1 -1 -1 0 -1 1 4 1 0 0
2
-1
5
-5 3 -3 -5 -3 -5 -3 -3 -5 -3
4
0 0 2 1 0 1
5
4 2 -3 -3 2 2 -3 2 2 -3
2
-3
3
-2 -2 -2
4
3 -5 -5 1 1 -5
2
-2
4
-1 1 -1 1 -1 1
4
3 4 1 1 3 1
4
-4 -4 -4 2 -1 -1
4
-1 5 -1 -1 -1 -1
4
-3 -5 -5 -3 5 -5
5
-5 -2 -5 0 4 -5 -2 -2 0 -5
4
-4 -3 -4 -4 -1 -3
2
-2
3
-1 0 -1
4
-5 -4 2 -4 -5 -5
2
1
4
-5 -3 -5 -5 -3 -3
3
-3 -4 -4
3
-5 -5 0
4
-3 -3 -3 -2 -3 -3
5
1 1 1 -5 -5 1 1 -5 -5 1
2
-5
3
-2 -1 -2
5
-4 -4 -1 -3 -3 -3 -4 -4 -1 3
4
-3 -3 -3 -2 0 -2
3
-5 3 -5
4
1 -5 1 -5 -5 5
2
-5
2
-4
5
1 -4 -4 -4 -1 2 -1 -4 1 -1
`

type testCase struct {
	input    string
	expected string
}

func solve(n int, arr []int) string {
	sort.Ints(arr)
	m := len(arr)
	if m == 0 {
		return "0"
	}
	k := 0
	res := make([]string, n)
	for i := 1; i < n; i++ {
		res[i-1] = strconv.Itoa(arr[k])
		k += n - i
	}
	res[n-1] = strconv.Itoa(arr[m-1])
	return strings.Join(res, " ")
}

func loadCases() ([]testCase, error) {
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
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		pos++
		m := n * (n - 1) / 2
		if pos+m > len(fields) {
			return nil, fmt.Errorf("case %d: incomplete list", caseIdx+1)
		}
		arr := make([]int, m)
		for i := 0; i < m; i++ {
			val, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", caseIdx+1, err)
			}
			arr[i] = val
		}
		pos += m
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
			expected: solve(n, arr),
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
		fmt.Println("usage: verifierC /path/to/binary")
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
