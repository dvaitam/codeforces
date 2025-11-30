package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesA = `7 7 1 5 9 8 7 13 14
5 8 12 7 5 10 5
2 10 9 18
10 3 3 1 6 1 6 3 4 5 1 3
7 6 10 11 4 9 8 8 9
5 1 1 1 2 1 2
6 4 6 2 4 4 4 3
9 8 3 3 11 16 4 10 10 4 11
9 4 5 8 2 7 6 4 5 3 4
3 1 2 2 1
2 3 2 1
2 9 13 17
5 9 8 7 14 9 15
8 6 2 6 10 2 8 10 11 6
4 4 1 5 2 4
6 3 3 4 1 1 2 6
4 1 1 1 1 1
10 10 4 13 3 12 4 2 20 1 7 6
2 8 7 2
1 9 14
10 2 3 1 2 1 3 3 4 2 1 4
1 10 4
7 4 5 6 8 3 4 1 3
3 6 9 5 2
10 8 6 1 16 14 10 12 13 9 5 1
8 2 3 1 3 2 2 4 3 3
6 10 20 5 10 13 14 3
1 10 7
6 3 2 2 6 4 4 6
10 7 1 7 14 12 10 7 13 11 12 1
3 8 3 9 6
8 9 16 18 1 2 16 11 10 15
1 7 4
9 2 2 1 4 4 3 1 2 1 1
9 10 4 7 4 20 7 10 9 6 4
8 7 11 2 1 5 8 13 13 2
5 3 6 5 6 6 3
2 3 3 1
1 1 1
5 9 11 12 2 16 15
7 6 9 3 4 7 10 5 1
3 3 3 3 3
6 2 3 1 1 3 2 2
10 5 6 7 9 3 5 2 8 4 1 5
3 9 3 10 13
6 5 7 2 2 9 8 8
6 6 2 8 2 12 8 7
1 5 6
3 3 6 5 4
2 2 1 2
4 1 2 1 1 2
9 9 10 15 16 7 14 3 12 8 9
10 3 4 2 3 1 1 6 1 5 4 6
4 2 4 4 3 2
1 4 3
2 4 8 7
6 9 5 4 16 5 13 14
9 8 11 16 16 7 8 1 11 11 11
1 9 5
5 10 5 13 19 10 16
2 2 1 1
4 3 1 3 1 4
6 3 2 6 4 3 5 4
9 9 2 3 17 3 14 7 10 18 14
8 7 10 10 4 14 14 13 1 11
1 3 3
9 10 9 11 3 16 9 10 14 13 13
1 3 6
3 4 5 6 1
1 8 14
3 8 3 5 12
7 1 2 2 2 1 1 2 1
1 1 1
6 2 3 2 4 4 1 1
10 8 11 4 10 5 13 10 4 7 2 13
8 6 4 8 6 11 2 1 1 8
5 1 1 1 1 2 2
2 3 4 5
7 2 1 4 1 1 4 2 1
8 7 11 7 1 8 14 6 12 5
2 6 2 2
6 1 2 2 1 1 1 2
2 10 5 7
1 4 2
1 5 6
1 10 8
3 3 4 1 4
6 5 3 1 4 6 6 8
5 5 9 6 3 10 2
2 9 10 6
7 3 2 2 3 5 2 2 2
5 6 7 11 1 3 10
1 7 2
2 3 4 3
9 7 12 3 10 7 5 11 6 2 4
8 6 11 9 1 7 7 1 7 12
6 8 7 12 10 16 3 6
2 5 2 9
10 3 6 4 4 2 4 4 2 2 4 3
9 3 3 4 6 6 1 4 2 3 1`

// Embedded solver from 677A.go.
func solve(n, h int, heights []int) string {
	width := 0
	for _, a := range heights {
		if a > h {
			width += 2
		} else {
			width++
		}
	}
	return strconv.Itoa(width)
}

type testCase struct {
	n       int
	h       int
	heights []int
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesA), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err1 := strconv.Atoi(fields[0])
		h, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: bad n or h", idx+1)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d heights got %d", idx+1, n, len(fields)-2)
		}
		heights := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad height %d", idx+1, i+1)
			}
			heights[i] = val
		}
		cases = append(cases, testCase{n: n, h: h, heights: heights})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.n, tc.h, tc.heights)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.h)
		for i, v := range tc.heights {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
