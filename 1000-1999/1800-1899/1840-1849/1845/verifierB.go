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
3 19 13 13 7 7
3 19 11 3 10 2
11 20 17 10 16 13
19 3 8 11 1 3
20 16 12 6 10 16
6 15 9 2 16 9
10 10 17 7 16 19
7 14 8 16 17 3
15 2 17 20 8 7
4 15 13 13 20 16
3 14 9 11 11 18
4 9 5 9 15 20
8 18 14 19 19 15
6 1 19 15 13 7
10 2 9 2 15 15
2 15 8 16 7 12
3 3 6 2 1 1
16 10 18 16 1 2
6 7 8 20 3 12
19 10 4 11 5 5
7 16 9 19 20 12
4 3 10 8 17 10
7 13 17 6 11 6
6 4 16 19 18 12
3 7 20 9 19 1
12 11 13 2 15 4
9 17 1 10 15 9
8 6 6 15 3 15
1 8 8 3 20 3
6 9 19 4 20 1
5 14 18 18 6 19
12 10 15 8 16 20
15 12 3 7 19 2
14 16 18 8 19 20
17 11 12 20 4 9
7 2 3 15 4 6
13 7 16 9 10 1
15 11 13 19 2 16
1 14 16 2 11 10
19 14 16 3 1 2
19 17 16 13 4 20
14 7 12 5 16 12
5 5 4 11 3 16
2 1 18 8 6 13
13 12 18 13 10 10
6 18 18 1 3 13
3 12 15 2 16 19
20 15 4 15 14 17
10 13 10 15 11 9
3 9 20 13 6 9
1 1 14 2 3 18
9 15 14 3 19 15
20 5 6 7 16 18
14 18 9 14 14 4
16 15 10 7 17 16
16 19 4 17 13 2
1 16 4 10 7 13
15 12 8 14 17 13
9 5 1 17 4 15
11 18 9 6 6 10
15 11 19 9 10 11
8 2 17 17 16 17
14 19 10 8 6 6
19 7 17 4 13 6
8 4 2 13 20 14
6 3 19 6 8 17
17 13 15 12 5 19
18 18 11 16 13 18
13 3 18 17 6 15
5 10 20 12 11 16
10 15 11 9 5 11
3 9 15 7 1 1
12 13 8 6 5 12
9 1 20 6 19 15
15 19 14 1 11 17
5 12 7 12 18 13
10 18 4 18 19 6
7 7 8 20 12 9
19 6 13 6 7 10
7 1 17 18 19 19
19 13 3 7 15 15
5 8 13 18 9 9
18 8 11 3 2 2
11 8 12 5 1 18
6 5 13 9 18 13
13 17 18 10 12 5
14 19 1 5 20 14
19 4 15 13 4 16
19 1 11 15 20 18
14 20 1 10 9 20
15 1 1 16 5 19
11 18 15 7 4 6
15 6 12 7 2 14
4 18 5 10 5 8
1 9 10 15 2 9
12 4 16 4 1 13
15 8 7 15 14 9
19 14 8 16 12 4
8 17 10 13 5 17
16 15 10 16 13 13
14 13 8 4 4 13
15 11 15 4 5 18
12 14 1 14 19 20
4 10 17 17 14 8
7 11 8 12 16 7
8 7 3 13 6 20
5 11 12 12 6 19
8 3 11 5 16 1
17 12 9 16 15 3
13 6 8 16 18 5
16 10 6 5 7 8
4 8 8 5 15 10
9 2 13 7 4 10
14 13 11 16 13 16
7 15 12 14 11 15
17 11 16 9 19 19
2 6 10 12 6 10
17 2 16 11 12 6
16 3 11 1 15 8
1 20 11 18 18 9`

type testCase struct {
	input    string
	expected string
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solve(xA, yA, xB, yB, xC, yC int64) int64 {
	ans := int64(1)
	dx1 := xB - xA
	dx2 := xC - xA
	if dx1*dx2 > 0 {
		ans += minInt64(absInt64(dx1), absInt64(dx2))
	}
	dy1 := yB - yA
	dy2 := yC - yA
	if dy1*dy2 > 0 {
		ans += minInt64(absInt64(dy1), absInt64(dy2))
	}
	return ans
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
		if pos+5 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing coordinates", caseIdx+1)
		}
		vals := make([]int64, 6)
		for i := 0; i < 6; i++ {
			v, err := strconv.ParseInt(fields[pos+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad number: %w", caseIdx+1, err)
			}
			vals[i] = v
		}
		pos += 6
		var sb strings.Builder
		sb.WriteString("1\n")
		for i, v := range vals {
			sb.WriteString(strconv.FormatInt(v, 10))
			if i%2 == 1 {
				sb.WriteByte('\n')
			} else {
				sb.WriteByte(' ')
			}
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]), 10),
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
