package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
4
6 11 -4 -15
1 2
2 3
1 4
2
10 -20
1 2
4
-13 -1 -18 19
1 2
1 3
2 4
3
-7 7 -16
1 2
1 3
2
-18 -9
1 2
4
2 0 -14 6
1 2
1 3
2 4
4
4 13 14 20
1 2
1 3
2 4
4
-7 2 -11 20
1 2
1 3
2 4
6
7 -2 -7 6 4 4
1 2
2 3
3 4
1 5
5 6
5
-17 12 19 -19 19
1 2
2 3
2 4
2 5
3
-20 19 -3
1 2
2 3
2
-20 19
1 2
6
-3 9 8 -8 1 -3
1 2
1 3
2 4
1 5
2 6
2
18 -4
1 2
2
-9 19
1 2
2
-15 -16
1 2
5
18 -10 -15 8 4
1 2
2 3
2 4
4 5
3
13 -19 -16
1 2
2 3
3
0 3 8
1 2
2 3
6
14 -20 6 17 -18 20
1 2
1 3
1 4
2 5
2 6
2
-17 -15
1 2
6
-2 0 8 13 4 17
1 2
2 3
3 4
4 5
4 6
6
-18 4 -5 4 -17 -15
1 2
2 3
1 4
2 5
4 6
4
-6 -20 6 0
1 2
1 3
3 4
6
6 18 17 -14 14 -12
1 2
1 3
3 4
1 5
2 6
4
-3 15 20 -11
1 2
2 3
3 4
4
4 -6 -19 -7
1 2
1 3
3 4
2
11 -17
1 2
6
1 1 8 20 12 -4
1 2
2 3
1 4
4 5
3 6
3
10 14 -2
1 2
1 3
3
18 9 13
1 2
1 3
2
16 12
1 2
4
-8 13 1 -11
1 2
1 3
3 4
5
20 -11 6 -8 5
1 2
2 3
3 4
1 5
4
-18 -14 -5 17
1 2
2 3
2 4
3
3 -18 1
1 2
1 3
3
14 -1 -18
1 2
2 3
6
14 18 -9 1 17 19
1 2
2 3
1 4
2 5
2 6
2
5 4
1 2
4
-14 -9 6 3
1 2
2 3
3 4
5
4 19 12 4 -11
1 2
1 3
2 4
1 5
4
-15 -1 -3 -16
1 2
1 3
3 4
3
20 -2 -19
1 2
1 3
2
7 -5
1 2
6
6 -9 12 12 19 -3
1 2
2 3
3 4
1 5
3 6
3
6 9 11
1 2
1 3
5
1 2 3 -8 1
1 2
1 3
2 4
2 5
3
-12 -20 -12
1 2
2 3
5
14 19 2 -3 -14
1 2
1 3
2 4
3 5
5
-4 -5 -9 7 16
1 2
2 3
3 4
2 5
3
-1 -4 -1
1 2
1 3
5
13 -19 -9 9 12
1 2
2 3
1 4
3 5
6
-4 -15 -6 11 -7 16
1 2
1 3
3 4
3 5
5 6
3
-6 17 -10
1 2
2 3
3
-6 -2 18
1 2
2 3
5
4 -1 12 16 6
1 2
1 3
2 4
3 5
2
8 14
1 2
2
11 -5
1 2
6
13 11 6 -9 2 7
1 2
1 3
3 4
1 5
4 6
3
5 -14 13
1 2
1 3
4
20 19 4 2
1 2
1 3
2 4
6
12 10 -3 -7 -11 2
1 2
1 3
2 4
3 5
3 6
2
20 14
1 2
5
-2 -12 -13 17 9
1 2
1 3
2 4
2 5
2
1 -5
1 2
3
-18 4 2
1 2
1 3
5
15 -16 16 -9 18
1 2
1 3
1 4
4 5
4
15 3 10 4
1 2
2 3
1 4
6
-3 20 -20 -10 15 17
1 2
2 3
3 4
2 5
5 6
6
14 -1 -6 14 0 -8
1 2
2 3
1 4
4 5
3 6
2
-13 -2
1 2
4
1 2 18 14
1 2
2 3
1 4
2
10 -19
1 2
4
-18 18 -14 -6
1 2
1 3
1 4
6
-6 -12 20 0 6 -15
1 2
1 3
3 4
3 5
3 6
6
-14 -9 14 13 10 -8
1 2
2 3
1 4
2 5
2 6
6
-17 1 20 13 4 11
1 2
2 3
3 4
4 5
5 6
3
10 11 8
1 2
1 3
6
8 -5 -4 -4 -4 -12
1 2
2 3
1 4
4 5
1 6
5
18 -18 -16 -17 -5
1 2
1 3
1 4
4 5
6
17 -17 17 16 -8 -6
1 2
2 3
2 4
1 5
2 6
4
5 -10 -10 -8
1 2
2 3
1 4
6
-11 -6 -4 9 12 2
1 2
2 3
3 4
1 5
1 6
4
-11 -8 -18 4
1 2
1 3
2 4
6
-7 -18 -7 8 1 -5
1 2
2 3
3 4
3 5
5 6
5
13 9 -10 6 20
1 2
1 3
3 4
2 5
5
-15 -16 2 19 19
1 2
2 3
1 4
2 5
4
-3 19 -8 -16
1 2
1 3
2 4
4
-11 -14 -10 -7
1 2
1 3
3 4
5
-17 -6 -7 -5 -15
1 2
1 3
2 4
3 5
3
4 -13 -3
1 2
2 3
2
-10 -1
1 2
5
12 18 13 -14 13
1 2
2 3
2 4
2 5
3
12 -8 13
1 2
2 3
4
-19 5 -7 3
1 2
2 3
1 4
2
-12 8
1 2
5
14 2 12 -19 7
1 2
2 3
1 4
1 5
5
-13 -11 -7 -7 14
1 2
2 3
1 4
4 5
5
-8 -15 -4 9 -6
1 2
1 3
1 4
4 5
5
-14 -13 12 14 19
1 2
2 3
1 4
4 5`

type testCase struct {
	input    string
	expected string
}

func expected(n int, edges [][2]int) string {
	degree := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		degree[u]++
		degree[v]++
	}
	leaves := 0
	for i := 2; i <= n; i++ {
		if degree[i] == 1 {
			leaves++
		}
	}
	return fmt.Sprintf("%d", leaves)
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
		values := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing value", caseIdx+1)
			}
			v, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", caseIdx+1, err)
			}
			values[i] = v
			pos++
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d: missing edge", caseIdx+1)
			}
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad edge u: %w", caseIdx+1, err)
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad edge v: %w", caseIdx+1, err)
			}
			edges[i] = [2]int{u, v}
			pos += 2
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(values[i]))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: expected(n, edges),
		})
	}
	return cases, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
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
		got, err := runExe(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
