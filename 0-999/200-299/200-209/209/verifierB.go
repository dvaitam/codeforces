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
	a, b, c int64
}

// Embedded testcases (previously stored in testcasesB.txt) so verifier is self contained.
const rawTestcasesB = `
4 8 3
3 10 9
6 9 4
7 7 10
10 5 1
5 9 1
7 9 10
5 3 3
0 4 1
3 5 2
5 6 0
1 2 3
0 9 10
8 9 10
1 0 1
10 3 9
9 1 6
1 5 1
0 9 0
3 2 1
7 3 0
10 0 8
6 9 1
4 1 3
1 10 4
5 6 2
0 8 7
0 9 1
6 3 4
5 7 9
2 10 3
0 10 2
2 5 8
4 1 9
7 10 2
0 7 10
6 9 8
4 10 5
6 10 4
2 8 0
7 1 5
0 8 4
2 3 7
5 9 4
10 5 9
10 9 2
4 6 6
10 1 0
9 3 5
2 3 3
10 7 6
10 9 6
0 6 9
6 10 0
2 7 1
4 2 7
8 7 8
9 0 0
7 5 4
7 0 6
3 8 10
1 2 0
6 10 6
5 0 3
0 0 10
8 9 1
3 1 9
10 3 4
4 2 1
7 6 10
1 0 4
7 1 4
2 10 8
10 10 5
1 2 4
1 0 0
3 10 4
8 5 5
9 0 9
10 7 10
7 10 6
5 8 2
3 6 9
4 0 2
2 4 5
5 5 1
5 9 0
0 4 2
2 9 4
5 6 8
2 4 1
7 3 0
4 2 8
1 4 6
5 4 6
1 1 8
7 7 5
5 1 7
1 7 6
0 4 5
`

func loadTestcases() ([]testCase, error) {
	fields := strings.Fields(rawTestcasesB)
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("unexpected token count %d (want multiple of 3)", len(fields))
	}
	cases := make([]testCase, 0, len(fields)/3)
	for i := 0; i < len(fields); i += 3 {
		a, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse a at token %d (%q): %w", i+1, fields[i], err)
		}
		b, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse b at token %d (%q): %w", i+2, fields[i+1], err)
		}
		c, err := strconv.ParseInt(fields[i+2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse c at token %d (%q): %w", i+3, fields[i+2], err)
		}
		cases = append(cases, testCase{a: a, b: b, c: c})
	}
	return cases, nil
}

// expected mirrors the logic from 209B.go so the verifier has no external oracle.
func expected(a, b, c int64) int64 {
	counts := []int64{a, b, c}
	ans := int64(-1)
	for i := 0; i < 3; i++ {
		u := counts[(i+1)%3]
		v := counts[(i+2)%3]
		if (u+v)&1 == 1 {
			continue
		}
		f := u
		if v > u {
			f = v
		}
		if ans == -1 || f < ans {
			ans = f
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		exp := expected(tc.a, tc.b, tc.c)
		input := fmt.Sprintf("%d %d %d\n", tc.a, tc.b, tc.c)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err2 := strconv.ParseInt(gotStr, 10, 64)
		if err2 != nil || got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx+1, exp, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
