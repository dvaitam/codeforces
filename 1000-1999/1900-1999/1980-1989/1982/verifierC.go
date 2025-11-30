package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesC = `100
2 4 6
5 6
4 1 14
3 10 1 7
1 2 4
7
3 9 12
3 10 7
3 6 7
4 8 6
6 9 9
7 7 1 7 6 8
2 6 10
8 2
2 2 6
2 9
5 3 15
8 7 3 7 7
2 4 11
6 9
2 6 13
2 8
7 4 8
1 8 10 8 1 4 5
1 5 13
10
2 7 14
2 8
7 4 12
7 5 1 2 5 1 1
3 7 15
10 7 8
1 5 10
5
7 4 13
2 1 2 5 5 9 6
1 9 15
4
7 3 4
7 5 5 9 3 10 9
6 4 12
2 7 9 7 5 5
4 6 15
3 3 2 2
4 7 14
3 9 5 6
6 8 14
4 8 8 9 6 8
6 1 8
5 3 8 1 10 4
1 6 13
7
1 9 9
2
6 7 7
6 1 2 10 1 5
7 5 8
3 10 5 4 2 7 8
6 6 12
3 6 7 7 3 8
6 3 11
6 3 4 3 8 6
7 7 13
8 7 4 4 8 4 10
6 1 15
7 1 4 2 3 6
1 3 6
10
3 10 10
9 5 6
4 8 8
9 9 7 10
4 8 12
8 4 6 5
1 1 1
3
3 1 5
1 3 2
7 7 10
10 7 9 4 8 4 6
5 2 11
2 6 6 9 8
3 5 5
9 1 4
3 2 5
9 6 4
7 4 8
5 5 9 7 5 8 6
7 4 4
5 9 2 1 8 8 8
1 7 14
8
4 2 3
2 4 2 3
4 4 11
10 2 7 9
7 7 7
3 4 8 4 3 5 6
3 7 8
9 5 10
5 4 15
5 8 9 10 8
5 5 9
4 1 2 10 2
2 7 10
4 5
6 1 12
9 9 7 1 2 7
6 5 6
10 6 4 9 5 4
6 4 5
9 5 6 4 6 8
3 10 11
3 1 9
5 6 11
10 1 3 7 3
2 9 9
3 4
7 8 11
4 3 4 7 6 10 10
2 8 9
10 1
5 10 12
8 8 5 1 4
5 3 13
8 8 9 6 2
3 3 12
7 4 6
7 5 11
1 4 1 6 4 6 8
6 4 8
6 3 5 1 6 10
5 1 12
3 6 1 8 1
1 4 4
1
2 6 7
1 6
6 7 9
4 8 7 3 6 5
2 6 12
7 1
4 5 13
9 8 1 10
1 7 13
3
1 9 10
10
7 9 15
3 2 6 4 3 4 1
7 3 14
9 3 2 7 10 2 10
6 8 10
10 10 1 5 6 7
1 1 8
2
3 5 15
3 8 4
5 6 8
7 6 5 8 7
1 5 13
5
5 8 8
9 10 9 5 1
4 7 8
7 6 8 1
1 5 5
5
6 10 15
5 4 9 9 6 7
7 5 8
2 10 6 4 10 9 6
2 3 8
1 10
7 1 10
3 6 6 5 5 6 8
7 7 13
3 1 3 10 1 8 3
3 1 12
8 5 10
2 2 10
7 5
7 3 11
3 2 3 10 2 9 9
5 7 13
5 5 5 1 7
7 5 9
9 9 9 6 6 4 7
7 3 3
9 3 10 7 6 8 1
5 7 10
1 6 9 3 4
6 6 13
1 4 10 4 5 3
7 7 8
10 8 4 8 9 2 4
2 8 9
7 7
3 5 11
6 10 6`

func solveCase(n int, l, r int64, arr []int64) int {
	sum := int64(0)
	wins := 0
	for _, v := range arr {
		sum += v
		if sum > r {
			sum = 0
		} else if sum >= l {
			wins++
			sum = 0
		}
	}
	return wins
}

type testCase struct {
	n   int
	l   int64
	r   int64
	arr []int64
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesC)
	pos := 0
	readInt := func() (int64, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v, err
	}
	t64, err := readInt()
	if err != nil {
		return nil, err
	}
	t := int(t64)
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n64, err := readInt()
		if err != nil {
			return nil, err
		}
		l, err := readInt()
		if err != nil {
			return nil, err
		}
		r, err := readInt()
		if err != nil {
			return nil, err
		}
		n := int(n64)
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			val, err := readInt()
			if err != nil {
				return nil, err
			}
			arr[j] = val
		}
		tests[i] = testCase{n: n, l: l, r: r, arr: arr}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.l, tc.r)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}

	allInput := buildAllInput(tests)
	allOutput, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}
	outLines := strings.Split(strings.TrimSpace(allOutput), "\n")
	if len(outLines) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outLines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.n, tc.l, tc.r, tc.arr)
		if strings.TrimSpace(outLines[i]) != strconv.Itoa(want) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput n=%d l=%d r=%d arr=%v\nexpected: %d\ngot: %s\n", i+1, tc.n, tc.l, tc.r, tc.arr, want, outLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
