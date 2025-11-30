package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesB = `100
5 3 11 5
0 -10 -8 -4 -3
4 -2 2 0

1 2 4 3
9 -4 6
-3 3 17 3
8 -7 -5
-4 -3 1 2
-7 10
-5 -3 10 1
9
5 -1 14 5
3 7 -10 -6 -1
-1 2 14 2
-6 -2
2 2 10 2
-6 7
0 -3 11 2
-4 -10
4 0 20 4
-8 -7 4 -3
0 -3 13 2
0 -2
5 3 5 3
10 -10 -9
-1 -1 2 3
-5 -2 6
-4 -3 7 2
1 -4
2 2 18 5
4 8 -7 -4 -3
3 2 13 1
-8
-2 3 1 1
5
-2 3 6 4
8 2 3 -4
0 0 3 2
-3 6
-5 -3 1 5
0 2 3 7 9
3 2 3 4
-7 4 -1 6
-5 3 3 5
1 3 9 10 -1
-4 3 17 5
10 -10 -6 -3 -2
-3 -2 16 0

4 2 9 1
-1
-2 -1 6 2
0 5
-2 1 20 3
-8 -7 -10
-2 0 2 4
0 -8 3 -3
3 -2 13 5
0 4 7 8 -4
-2 -3 2 0

1 -2 19 3
0 9 -4
5 2 17 4
-8 2 -4 5
0 3 1 5
0 1 2 5 -2
-1 0 0 4
1 4 7 -9
-3 -1 20 0

-3 -2 2 4
-8 2 3 -2
1 -3 9 4
8 10 5 -2
-2 1 11 3
-5 -4 -2
-1 -3 8 0

-4 -3 15 2
-1 7
1 -2 5 5
0 4 8 10 -3
4 0 13 1
4
-4 3 15 1
10
1 3 12 5
6 8 -10 -6 -2
0 3 5 4
-8 9 4 1
4 2 11 1
8
2 2 5 1
0
1 -2 1 3
-8 -7 -9
1 -1 8 1
10
-3 0 3 3
-7 -4 -1
1 2 13 2
-8 -2
5 2 6 0

-4 -1 17 2
9 -3
4 -2 10 0

-3 1 5 3
7 -3 -9
2 -1 3 5
4 6 10 -9 -4
-4 -2 5 4
4 -4 -1 6
-3 2 13 5
1 2 10 -8 -3
5 2 8 2
-6 -10
-4 0 5 3
0 9 -3
-4 1 6 4
0 -4 -2 -10
2 0 20 4
2 -4 -10 7
2 -3 19 4
-4 5 -2 7
0 -3 19 3
8 10 3
-5 -2 11 3
10 -2 -9
5 -3 12 3
-5 -4 7
1 1 1 5
2 9 -10 -8 -7
4 -1 12 5
1 2 4 5 -9
3 0 7 1
5
3 0 6 5
0 2 7 8 -7
4 1 11 5
10 -8 -5 -3 -1
3 -1 3 2
0 -2
-1 0 15 0

3 -1 0 4
7 -3 -1 6
-3 3 1 3
2 6 -2
-1 -2 5 3
2 -4 -2
-3 -2 13 0

2 2 10 1
10
-2 -2 7 3
-7 -3 -10
2 -1 18 1
-9
-3 1 4 5
0 4 5 8 -4
-3 3 10 1
-7
-3 -2 10 0

-2 3 16 2
0 7
1 2 7 5
2 4 -10 -7 -3
2 3 17 1
9
-1 0 14 3
0 9 -3
2 -1 2 4
0 3 5 7
-4 2 19 0

0 0 9 4
-7 3 4 1
3 -1 18 0

5 3 1 4
2 -5 5 -2
3 3 2 5
9 10 -4 -3 -2
5 0 11 0

2 0 1 2
8 -7
0 -2 0 3
-7 -5 -9
3 2 4 5
1 3 4 -7 -6
-1 1 7 0

-3 -2 6 5
1 6 8 -7 -5`

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Embedded solver from 789B.go.
func expected(b1, q, l int64, bad map[int64]bool) string {
	if abs(b1) > l {
		return "0"
	}
	if b1 == 0 {
		if bad[0] {
			return "0"
		}
		return "inf"
	}
	if q == 0 {
		if !bad[b1] {
			if bad[0] {
				return "1"
			}
			return "inf"
		}
		if bad[0] {
			return "0"
		}
		return "inf"
	}
	if q == 1 {
		if bad[b1] {
			return "0"
		}
		return "inf"
	}
	if q == -1 {
		if bad[b1] && bad[-b1] {
			return "0"
		}
		return "inf"
	}
	count := 0
	cur := b1
	for abs(cur) <= l {
		if !bad[cur] {
			count++
		}
		cur *= q
	}
	return strconv.Itoa(count)
}

type testCase struct {
	b1  int64
	q   int64
	l   int64
	bad []int64
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesB), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	lineIdx := 0
	// first line t
	for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
		lineIdx++
	}
	if lineIdx >= len(lines) {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[lineIdx]))
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	lineIdx++
	cases := make([]testCase, 0, t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
			lineIdx++
		}
		if lineIdx >= len(lines) {
			return nil, fmt.Errorf("case %d: missing header", caseNum)
		}
		header := strings.Fields(lines[lineIdx])
		lineIdx++
		if len(header) != 4 {
			return nil, fmt.Errorf("case %d: bad header", caseNum)
		}
		b1, err1 := strconv.ParseInt(header[0], 10, 64)
		q, err2 := strconv.ParseInt(header[1], 10, 64)
		l, err3 := strconv.ParseInt(header[2], 10, 64)
		m, err4 := strconv.Atoi(header[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return nil, fmt.Errorf("case %d: parse error", caseNum)
		}
		var badVals []int64
		if m > 0 {
			for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
				lineIdx++
			}
			if lineIdx >= len(lines) {
				return nil, fmt.Errorf("case %d: missing bad numbers", caseNum)
			}
			badFields := strings.Fields(lines[lineIdx])
			lineIdx++
			if len(badFields) != m {
				return nil, fmt.Errorf("case %d: expected %d bad numbers got %d", caseNum, m, len(badFields))
			}
			badVals = make([]int64, m)
			for i, s := range badFields {
				v, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("case %d: bad bad-number", caseNum)
				}
				badVals[i] = v
			}
		} else {
			for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
				lineIdx++
			}
		}
		cases = append(cases, testCase{b1: b1, q: q, l: l, bad: badVals})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d %d\n", tc.b1, tc.q, tc.l, len(tc.bad))
		for i, v := range tc.bad {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		if len(tc.bad) > 0 {
			input.WriteByte('\n')
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		cmd := exec.CommandContext(ctx, bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		cancel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, stderr.String())
			os.Exit(1)
		}
		badMap := make(map[int64]bool, len(tc.bad))
		for _, v := range tc.bad {
			badMap[v] = true
		}
		expect := expected(tc.b1, tc.q, tc.l, badMap)
		outStr := strings.TrimSpace(out.String())
		if outStr != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
