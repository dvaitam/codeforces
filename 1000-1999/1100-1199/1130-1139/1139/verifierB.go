package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
3 18 2 8
2 15 14
8 20 12 6 3 15 0 12 13
10 0 14 8 7 18 3 10 0 0 0
9 0 12 6 13 0 16 7 14 15
9 7 11 7 7 14 9 0 13 17
2 5 20
5 3 10 16 13 16
4 9 9 18 15
9 12 18 1 15 7 12 13 5 11
9 11 2 14 16 3 5 16 12 11
8 0 15 1 9 19 18 18 12
3 5 16 7
1 6
9 17 7 12 16 11 18 11 14 8
9 19 0 12 16 4 16 17 6 13
1 15
6 18 17 6 16 13 15
6 13 11 0 17 17 19
10 10 14 19 0 7 20 5 17 18 5
2 17 8
1 2
2 0 14
1 8
4 8 3 19 5
6 9 2 5 5 8 16
3 8 20 9
8 10 15 15 3 0 9 12 10
7 6 8 3 8 16 6 19
7 0 7 0 12 4 1 5
8 16 13 17 7 20 16 14 7
9 20 0 12 18 10 20 13 1 9
3 6 1 9
2 2 9
5 5 13 18 8 4
1 17
1 18
4 18 14 5 19
9 1 12 6 11 3 6 18 13 18
4 15 3 12 9
9 15 0 10 19 12 9 0 5 6
6 18 4 10 13 6 8
2 12 17
6 17 15 17 7 2 1
2 4 5
3 17 6 8
6 19 16 8 11 10 10
2 9 7
10 15 4 18 17 3 10 1 13 2 12
3 4 10 3
10 18 12 2 18 17 7 18 2 8 11
5 18 17 3 14 8
2 1 9
1 19
1 2
7 3 1 6 7 18 13 5
2 14 5
4 5 3 13 12
9 9 17 8 15 10 3 6 20 10
1 0
1 9
10 10 14 12 10 12 2 2 10 19 14
2 8 6
10 17 15 11 8 5 17 6 9 6 7
6 2 8 2 14 2 20
10 20 10 7 12 9 1 10 5 10 18
5 7 10 3 17 19
10 19 2 7 7 0 7 12 2 8 17
2 2 0
1 9
6 15 15 4 3 16 10
2 16 5
3 4 4 10
5 3 16 19 9 4
4 4 17 1 10
10 17 6 5 9 13 17 5 1 7 8
2 14 13
9 8 17 14 17 14 0 12 10 5
5 15 0 20 13 18
1 1
6 18 4 18 4 4 8
5 12 18 12 5 19
2 7 15
1 5
9 10 16 20 14 20 7 7 10 15
8 7 13 10 17 19 20 8 20
4 1 2 16 20
6 5 16 6 9 9 9
9 11 5 14 19 2 3 19 16 18
7 5 4 8 13 6 18 1
8 12 20 11 12 16 5 17 1
9 2 8 20 3 8 2 4 19 2
8 7 12 13 12 5 10 14 4
10 15 6 3 13 19 17 13 3 9 8
4 12 17 0 6
9 14 18 0 0 20 19 7 8 6
3 9 4 17
4 8 9 18 8
8 5 17 11 15 13 3 6 18
7 6 9 3 0 3 18 0
`

type testCase struct {
	n int
	a []int64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesData, "\n")
	var cases []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("length mismatch for n=%d in line %q", n, line)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				return nil, err
			}
			arr[i] = val
		}
		cases = append(cases, testCase{n: n, a: arr})
	}
	return cases, nil
}

// solve mirrors 1139B.go so we don't need an external oracle.
func solve(tc testCase) int64 {
	n := tc.n
	a := tc.a
	if n == 0 {
		return 0
	}
	h := a[n-1]
	ans := h
	for i := n - 2; i >= 0 && h > 0; i-- {
		if h > a[i] {
			h = a[i]
		} else {
			h--
		}
		ans += h
	}
	return ans
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := buildInput(tc)
		expected := strconv.FormatInt(solve(tc), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
