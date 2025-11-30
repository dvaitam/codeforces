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
5 3 4 2 1 5
10 3 6 5 4 1 2 9 7 10 8
9 7 4 5 9 2 6 3 8 1
3 2 3 1
9 8 3 9 2 5 1 6 7 4
6 1 6 4 5 2 3
3 3 1 2
5 4 2 5 1 3
9 1 2 6 8 5 7 3 4 9
3 3 1 2
9 6 4 2 1 9 5 8 3 7
5 1 3 2 4 5
3 3 1 2
7 2 1 3 4 7 5 6
5 4 3 1 2 5
6 1 2 6 3 5 4
10 7 8 3 2 10 6 4 1 9 5
9 9 4 6 2 5 3 7 8 1
10 4 9 2 3 8 5 1 10 7 6
7 3 1 4 6 5 2 7
6 5 3 4 6 2 1
8 4 2 8 5 6 1 3 7
6 4 5 1 2 3 6
6 1 5 6 2 4 3
8 1 2 4 7 3 6 8 5
4 1 3 2 4
7 5 4 3 2 7 6 1
7 4 3 2 1 7 6 5
4 1 2 3 4
6 4 3 1 2 6 5
10 6 8 2 10 5 7 1 3 4 9
5 3 1 2 4 5
7 7 6 4 2 5 3 1
6 5 3 2 1 6 4
3 1 2 3
4 2 3 4 1
3 1 2 3
9 4 1 2 8 9 3 7 5 6
3 1 2 3
7 4 7 2 6 1 3 5
8 7 5 1 2 6 4 8 3
4 1 2 3 4
9 6 1 9 3 7 2 8 5 4
6 4 2 3 6 5 1
3 1 3 2
9 7 3 5 9 8 4 6 1 2
4 2 1 3 4
4 1 3 2 4
8 8 6 1 5 7 2 4 3
7 6 4 5 7 3 1 2
7 1 2 3 6 4 7 5
8 2 1 6 4 8 5 3 7
6 3 6 5 1 4 2
3 2 1 3
6 2 4 6 5 1 3
6 1 3 5 6 4 2
6 2 3 5 1 6 4
7 5 4 3 6 7 1 2
3 1 2 3
7 2 1 3 6 4 7 5
8 2 4 8 3 7 5 1 6
5 2 1 4 5 3
7 7 4 2 1 6 5 3
6 6 5 1 2 4 3
9 6 3 1 9 2 7 5 8 4
4 1 4 3 2
7 1 4 6 2 7 3 5
10 8 4 1 2 5 6 10 7 3 9
3 1 3 2
10 8 10 1 6 2 3 9 4 5 7
10 4 2 9 7 5 10 3 8 1 6
10 6 3 2 10 7 5 1 8 4 9
6 1 3 6 4 5 2
3 2 3 1
9 3 9 7 2 5 4 6 1 8
5 5 4 1 3 2
9 4 9 3 5 7 6 1 2 8
8 4 8 3 5 6 7 2 1
9 3 6 1 7 5 4 2 8 9
8 4 3 7 2 6 5 1 8
6 5 2 6 1 3 4
3 3 1 2
4 2 3 1 4
7 7 1 6 5 3 4 2
10 7 6 5 2 1 8 3 10 9 4
8 7 8 1 4 3 6 5 2
7 4 3 5 7 6 1 2
8 8 7 3 1 2 6 4 5
5 5 2 1 3 4
4 4 3 2 1
10 6 9 8 4 7 5 3 2 10 1
3 2 1 3
4 4 3 1 2
9 8 6 4 1 9 7 2 3 5
4 2 3 1 4
6 6 4 2 1 3 5
7 4 3 1 5 7 6 2
3 2 3 1
3 3 2 1
7 5 6 4 3 1 7 2
10 1 2 10 4 3 8 6 5 7 9
3 2 3 1
6 5 1 2 6 3 4
4 4 2 1 3`

type testCase struct {
	input    string
	expected string
}

func solve(n int, perm []int) string {
	pos1, pos2, posn := -1, -1, -1
	for i, v := range perm {
		idx := i + 1
		if v == 1 {
			pos1 = idx
		} else if v == 2 {
			pos2 = idx
		} else if v == n {
			posn = idx
		}
	}
	minPos := pos1
	if pos2 < minPos {
		minPos = pos2
	}
	maxPos := pos1
	if pos2 > maxPos {
		maxPos = pos2
	}
	var iOut, jOut int
	if minPos < posn && posn < maxPos {
		iOut, jOut = pos1, pos2
	} else if posn < minPos {
		iOut, jOut = posn, minPos
	} else {
		iOut, jOut = maxPos, posn
	}
	return fmt.Sprintf("%d %d", iOut, jOut)
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
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: incomplete permutation", caseIdx+1)
			}
			val, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", caseIdx+1, err)
			}
			perm[i] = val
			pos++
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solve(n, perm),
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
