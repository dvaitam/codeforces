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
14
2 12 11 6 4 3 10 14 8 9 5 1 7 13
18
4 8 14 11 7 6 1 12 3 15 9 16 17 2 13 18 10 5
13
2 13 12 1 3 5 11 8 9 4 10 6 7
4
2 1 3 4
17
14 3 1 16 5 13 7 9 17 10 4 2 15 6 12 8 11
12
10 7 4 3 11 1 6 12 5 2 8 9
19
11 9 12 8 1 14 2 16 5 4 17 6 7 13 18 3 15 10 19
4
4 3 1 2
4
3 4 1 2
20
13 17 5 7 18 11 1 20 4 19 14 8 2 10 6 3 12 16 15 9
13
5 4 12 9 11 8 10 2 1 7 6 3 13
2
2 1
8
8 5 7 6 3 1 4 2
8
5 7 8 2 4 1 6 3
2
2 1
10
8 5 1 9 7 3 6 10 4 2
18
11 9 10 18 17 1 14 13 3 8 6 5 16 7 12 4 2 15
5
2 1 3 4 5
2
2 1
19
14 10 7 2 13 16 18 6 8 4 17 5 9 19 12 11 3 15 1
14
10 3 5 9 8 13 14 6 4 1 2 11 7 12
3
1 3 2
3
3 2 1
4
2 1 4 3
16
2 3 13 14 11 4 7 5 6 8 12 1 10 9 15 16
6
3 6 2 5 4 1
8
2 3 4 7 5 8 6 1
8
4 7 8 1 2 6 3 5
4
4 3 2 1
5
1 4 5 2 3
5
5 4 1 3 2
3
3 2 1
19
15 10 1 13 5 4 17 3 9 6 7 14 8 19 18 16 2 12 11
6
4 1 2 5 6 3
12
4 12 7 6 8 9 2 3 5 11 1 10
11
9 7 6 10 5 11 3 1 4 8 2
12
12 11 9 1 6 3 4 8 10 2 7 5
17
10 4 17 15 8 9 1 16 7 13 3 11 12 6 5 2 14
2
2 1
14
12 10 13 6 2 3 1 7 4 11 8 5 14 9
15
12 5 7 3 14 10 11 15 8 9 1 13 2 6 4
8
8 7 6 5 4 2 1 3
5
3 2 1 4 5
15
14 3 12 1 2 7 5 4 10 13 6 11 8 9 15
12
12 2 6 8 11 4 7 10 5 3 9 1
4
2 4 3 1
4
4 3 1 2
11
11 9 5 10 4 2 7 3 6 8 1
18
1 12 15 8 6 5 16 4 7 18 10 9 14 13 11 3 2 17
20
11 3 18 14 12 4 7 15 13 20 2 17 5 16 9 10 6 19 1 8
6
6 4 1 5 3 2
17
13 15 11 9 4 16 1 7 6 3 2 12 10 17 8 5 14
5
5 3 1 2 4
6
4 2 5 1 3 6
8
6 2 3 5 1 4 8 7
12
6 7 9 12 1 8 4 3 5 10 2 11
16
7 16 5 10 15 3 14 9 1 2 11 6 8 4 13 12
18
5 6 18 9 4 8 1 11 12 15 13 16 7 3 2 10 14 17
16
4 10 9 12 3 16 13 15 2 5 6 8 1 7 11 14
7
6 4 5 7 3 2 1
8
4 3 7 5 8 6 2 1
2
2 1
6
5 3 6 1 4 2
10
9 5 2 6 8 7 10 4 1 3
7
4 2 7 3 6 1 5
6
5 1 4 3 6 2
9
6 8 7 2 1 4 9 5 3
14
1 11 8 4 9 13 6 10 5 7 3 14 12 2
4
1 3 4 2
18
3 16 11 9 18 15 10 5 17 4 8 6 12 7 1 14 13 2
19
18 2 19 9 1 11 10 12 16 8 4 3 14 7 17 6 13 15 5
2
2 1
8
7 8 4 2 3 6 1 5
4
3 2 1 4
10
8 4 9 3 7 10 6 5 2 1
16
9 10 7 15 13 3 8 1 2 16 11 14 6 5 12 4
7
2 7 6 5 3 4 1
20
20 3 1 13 2 16 10 6 8 15 5 12 19 9 11 14 4 18 7 17
14
1 11 13 2 14 4 12 6 5 9 3 8 10 7
18
17 14 7 18 6 9 13 10 4 1 15 12 3 5 8 2 16 11
4
3 2 1 4
3
2 3 1
10
6 9 1 2 8 10 7 3 4 5
16
10 4 1 9 5 2 13 8 14 12 16 15 6 3 7 11
7
1 6 2 5 7 3 4
9
7 3 6 9 1 5 2 8 4
4
3 1 4 2
7
4 7 6 1 3 5 2
13
11 2 12 3 6 4 5 13 10 9 1 8 7
8
7 2 4 5 8 1 3 6
2
2 1
6
6 5 3 2 4 1
19
13 12 4 5 1 17 3 19 9 16 6 14 18 2 10 11 7 15 8
4
1 4 3 2
8
1 2 8 5 7 6 3 4
10
9 2 10 7 5 3 1 4 6 8
17
3 5 4 6 14 10 1 9 16 11 12 8 7 13 17 2 15
16
11 13 6 5 10 12 2 8 4 16 1 15 9 7 14 3
5
2 1 4 3 5
11
6 4 7 10 1 2 3 5 11 9 8`

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve(arr []int) int {
	g := 0
	for i, v := range arr {
		diff := v - (i + 1)
		if diff < 0 {
			diff = -diff
		}
		g = gcd(g, diff)
	}
	return g
}

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: incomplete permutation", caseIdx+1)
			}
			val, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", caseIdx+1, err)
			}
			arr[i] = val
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
