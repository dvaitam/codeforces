package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n   int
	x   int64
	arr []int64
}

// solve is the embedded logic from 1704B.go.
func solve(tc testCase) int {
	n := tc.n
	x := tc.x
	a := tc.arr
	l, r := a[0]-x, a[0]+x
	changes := 0
	for i := 1; i < n; i++ {
		nl, nr := a[i]-x, a[i]+x
		if nl > r || nr < l {
			changes++
			l, r = nl, nr
		} else {
			if nl > l {
				l = nl
			}
			if nr < r {
				r = nr
			}
		}
	}
	return changes
}

// Embedded copy of testcasesB.txt (each line: n x a1 ... an).
const testcaseData = `
4 1 9 11 5 9
4 2 10 4 14 8
5 9 7 11 11 17 13
5 8 4 5 15 17 18
6 10 17 18 1 10 6 7
3 7 17 11 4
4 6 5 19 3 2
3 9 11 14 10
3 6 9 11 17
5 1 17 4 5 11 11
3 10 3 15 9
4 8 12 13 3 19
1 3 2
5 8 19 9 8 19 11
3 6 13 10 15
5 6 18 17 6 1 5
3 4 19 5 4
2 7 20 2
1 9 9
6 2 7 9 3 19 17 3
1 4 6
5 7 1 19 12 16 10
2 4 20 16
2 7 15 12
5 4 16 3 9 14 7
1 9 13
5 8 3 13 20 17 19
5 7 2 12 15 1 7
3 1 18 4 10
5 6 18 19 18 10 17
4 9 17 14 20 19
3 8 10 5 17
4 10 5 18 6 9
6 1 14 19 2 12 14 13
3 1 3 3 1
4 5 15 9 12 16
3 7 15 4 16
3 3 14 5 1
2 5 12 5
5 5 14 9 17 10 14
6 5 14 11 16 7 16 13
6 7 3 3 5 7 5 8
6 1 4 9 5 16 4 13
6 3 1 3 14 20 2 18
2 9 14 12
1 2 18
6 7 4 9 9 6 16 2
2 2 13 4
6 8 10 17 16 13 4 20
4 2 5 13 20 7
2 9 9 14
6 9 10 16 18 7 20 11
4 2 1 12 9 2
5 8 10 4 8 17 9
3 4 14 5 5
3 4 14 18 20
1 9 20
5 3 14 9 9 16 10
3 8 7 16 12
5 8 8 11 6 20 6
6 10 15 18 5 2 17 11
5 3 7 11 20 16 16
3 2 5 5 9
2 2 18 2
5 3 4 8 19 7 17
5 5 14 11 1 1 10
5 4 3 8 9 11 9
5 9 13 1 4 11 12
2 2 9 5
6 10 2 12 3 3 4 10
3 4 9 17 2
3 1 3 5 13
3 4 4 11 9
1 9 11
1 6 5
5 5 13 3 19 20 17
4 10 14 18 13 10
2 5 18 5
1 10 17
1 3 8
2 7 9 18
1 5 18
3 9 9 16 5
4 2 12 3 18 12
5 9 17 19 1 20 10
4 3 5 3 19 5
6 4 16 11 12 10 6 5
4 8 13 4 20 5
3 5 20 1 18
1 3 13
6 9 4 15 1 14 20 14
3 6 14 13 20
4 1 4 16 2 1
1 2 19
2 9 17 12
5 5 19 12 16 8 20
2 2 18 12
2 2 2 11
4 6 9 2 20 14
4 7 12 10 11 15
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		x, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad x: %v", i+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d values got %d", i+1, 2+n, len(fields))
		}
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			val, err := strconv.ParseInt(fields[2+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = val
		}
		tests = append(tests, testCase{n: n, x: x, arr: arr})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected int) error {
	var input strings.Builder
	input.WriteString("1\n")
	fmt.Fprintf(&input, "%d %d\n", tc.n, tc.x)
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.FormatInt(v, 10))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	val, err := strconv.Atoi(got)
	if err != nil {
		return fmt.Errorf("non-integer output %q", got)
	}
	if val != expected {
		return fmt.Errorf("expected %d got %d", expected, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
